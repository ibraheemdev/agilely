class FullBoardQuery
  def initialize(slug)
    @slug = slug
  end

  def execute
    board = aggregation ||
    ( raise Mongoid::Errors::DocumentNotFound.new(Board, @slug) )
    handle_null_lists(board)
  end

  private

  def aggregation
    Board.collection.aggregate([
      find_board, 
      lookup_cards, 
      lookup_users, 
      project_board
    ]).first
  end

  def find_board
    { '$match' => { slug: @slug } }
  end

  def lookup_cards
    {
      '$lookup' => {
        from: "cards",
        localField: "lists._id",
        foreignField: "list_id",
        as: "cards"
      }
    }
  end

  def lookup_users
    {
      '$lookup' => {
        from: "users",
        as: "users",
        let: {
          "board_id": "$_id"
        },
        pipeline: [
          unwind_participations,
          find_users_with_participations,
          project_users
        ],
      }
    }
  end

  def unwind_participations
    { '$unwind' => '$participations' }
  end

  def find_users_with_participations
    {
      '$match' => {
        '$expr' => {
          '$and' => [{
              '$eq' => ["Board", "$participations.participant_type"]
            },
            {
              '$eq' => ["$$board_id", "$participations.participant_id"]
            },
          ]
        }
      }
    }
  end

  def project_users
    {
      '$project' => {
        email: true,
        name: true,
        participation: '$participations'
      }
    }
  end

  def project_board
    {
      '$project' => {
        _id: true,
        title: true,
        slug: true,
        public: true,
        users: '$users',
        lists: merge_cards_with_lists,
      }
    }
  end

  def merge_cards_with_lists
    {
      '$map' => {
        input: "$lists",
        as: "list",
        in: {
          '$mergeObjects' => [
            "$$list",
            {
              cards: {
                '$filter' => {
                  input: "$cards",
                  cond: {
                    '$eq' => ["$$this.list_id", "$$list._id"]
                  }
                }
              }
            }
          ]
        }
      }
    }
  end

  def handle_null_lists(board)
    board["lists"] ||= []
    return board
  end
end