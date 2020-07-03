class Card
  include ApplicationDocument
  include Midstring
  
  field :title, type: String
  validates :title, presence: true

  field :description, type: String

  field :position, type: String
  validates :position, presence: true

  belongs_to :list, index: true

  before_validation :set_position, on: :create

  delegate :board, to: :list

  private
  
  def set_position
    self.position = list.cards.blank? ? 'c' : midstring(list.cards.last.position, '')
  end
end
