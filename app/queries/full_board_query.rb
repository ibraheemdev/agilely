class FullBoardQuery

  def self.execute(*args, &block)
    new(*args, &block).execute
  end
  
  def initialize(board)
    @board = board
  end

  def execute
    lists = @board.lists
    cards = @board.cards
    {
      board: @board.as_json.merge("list_ids": lists.map(&:id)),
      lists: add_cards_to_lists(lists, cards),
      participants: participants, 
      cards: cards
    }
  end

  private

  def add_cards_to_lists(lists, cards)
    lists.map do |list| 
      list.as_json.merge(
        card_ids: cards.select { |card| card.list_id === list.id }.pluck(:id)
      )
    end
  end

  def participants
    @board.users.map do |u|
      u.participation_in(@board)
        .as_json
        .merge("name" => u.name, "email" => u.email)
    end
  end
end
