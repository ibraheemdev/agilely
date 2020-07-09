class Board
  include ApplicationDocument

  field :title, type: String
  validates :title, presence: true, length: { maximum: 512 }

  field :slug, type: String
  validates :slug, presence: true, length: { is: 8 }, uniqueness: true

  field :public, type: Boolean
  validates :public, inclusion: { in: [ true, false ] }

  has_many :lists, order: :position.asc
  has_many :cards, order: :position.asc

  before_validation :set_slug, on: :create

  def users
    User.where('participations.participant_type': "Board")
        .and('participations.participant_id': self.id)
        .all
  end

  def full_json
    as_json.merge(
      "lists" => List.full(self.id),
      "participants" => participants
    )
  end

  def participants
    users.map do |u|
      u.participation_in(self)
        .as_json
        .merge("name" => u.name, "email" => u.email)
    end
  end
  
  def to_param
    slug 
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