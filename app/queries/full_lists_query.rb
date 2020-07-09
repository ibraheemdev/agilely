class FullListsQuery

  def self.execute(*args, &block)
    new(*args, &block).execute
  end
  
  def initialize(board_id)
    @board_id = board_id
  end

  def execute
    aggregation.as_json
  end

  private

  def aggregation
    List.collection.aggregate([
      find_lists, 
      lookup_cards,
      project_lists,
      sort_lists
    ])
  end

  def find_lists
    { '$match' => { board_id: @board_id } }
  end

  def lookup_cards
    {
      '$lookup' => {
        from: "cards",
        pipeline: [
          { '$match' => { board_id: @board_id } },
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
        cards: filter_cards
       }
    }
  end

  def filter_cards
    { '$filter' => {
      input: "$cards",
      as: "card",
      cond: {
        '$eq' => ["$$card.list_id", "$_id"]
      }
    }}
  end

  def sort_lists
    { '$sort': { position: 1 } }
  end
end