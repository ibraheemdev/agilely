class Board
  include ApplicationDocument
  
  field :title, type: String
  validates :title, presence: true, length: { maximum: 512 }

  field :slug, type: String
  validates :slug, presence: true, length: { is: 8 }, uniqueness: true

  field :public, type: Boolean
  validates :public, inclusion: { in: [ true, false ] }

  embeds_many :lists, order: :position.asc
  
  before_validation :set_slug, on: :create

  def users
    User.where('participations.participant_type': "Board")
        .and('participations.participant_id': self.id)
        .all
  end

  def to_param() slug end

  def self.full_json(slug)
    Board.collection.aggregate([
      { '$match' => { slug: slug } },
      { '$lookup' => 
        { from: "cards", localField: "lists._id",
          foreignField: "list_id", as: "cards" } 
      },
      { '$project' => {
        _id: true, title: true,
        lists: { '$map' => { input: "$lists", as: "list",
            in: { 
              '$mergeObjects' => [
                "$$list",
                { cards: {
                    '$filter' => {
                      input: "$cards",
                      cond: {
                        '$eq' => [ "$$this.list_id", "$$list._id" ]
                      }
                    }
                }}
              ]
            }
         }}
       }}
    ]).as_json
  end

  def self.titles
    pluck(:title, :slug)
  end

  private

  def set_slug
    loop do
      self.slug = SecureRandom.alphanumeric(8)
      break unless Board.where(slug: slug).exists?
    end
  end
end
