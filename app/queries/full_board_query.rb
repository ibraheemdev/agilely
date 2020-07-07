class FullBoardQuery
  def initialize(slug)
    @slug = slug
  end

  def execute
    Board.collection.aggregate(aggregation).first
  end

  def aggregation
    [{
      '$match' => {
        slug: @slug
      }
    },
    {
      '$lookup' => {
        from: "cards",
        localField: "lists._id",
        foreignField: "list_id",
        as: "cards"
      }
    },
    {
      '$lookup' => {
        from: "users",
        as: "users",
        let: {
          "board_id": "$_id"
        },
        pipeline: [{
            '$unwind' => '$participations'
          },
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
          },
          {
            '$project' => {
              email: true,
              name: true,
              participation: '$participations'
            }
          }
        ],
      }
    },
    {
      '$project' => {
        _id: true,
        title: true,
        slug: true,
        public: true,
        users: '$users',
        lists: {
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
        },
      }
    }]
  end
end