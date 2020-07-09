class FullBoardQuery

  def self.execute(*args, &block)
    new(*args, &block).execute
  end
  
  def initialize(board)
    @board = board
  end

  def execute
    @board.as_json.merge("lists" => lists, "participants" => participants)
  end

  private

  def participants
    @board.users.map do |u|
      u.participation_in(@board)
        .as_json
        .merge("name" => u.name, "email" => u.email)
    end
  end

  def lists
    List.collection.aggregate([
      find_lists, 
      lookup_cards,
      project_lists,
      sort_lists
    ]).as_json
  end

  def find_lists
    { '$match' => { board_id: @board.id } }
  end

  def lookup_cards
    {
      '$lookup' => {
        from: "cards",
        pipeline: [
          { '$match' => { board_id: @board.id } },
          { '$sort': { position: 1 } }
       ],
       as: "cards"
      }
    }
  end

  def project_lists
    {
      '$project' => { 
        _id: true, 
        title: true, 
        position: true,
        created_at: true,
        updated_at: true,
        cards: list_cards("$_id")
       }
    }
  end

  def list_cards(list_id)
    { '$filter' => {
      input: "$cards",
      as: "card",
      cond: {
        '$eq' => ["$$card.list_id", list_id]
      }
    }}
  end

  def sort_lists
    { '$sort': { position: 1 } }
  end
end