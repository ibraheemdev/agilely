class Card
  include ApplicationDocument
  include Midstring
  
  field :title, type: String
  validates :title, presence: true

  field :description, type: String

  field :position, type: String
  validates :position, presence: true

  belongs_to :board, index: true
  belongs_to :list

  before_validation :set_position, on: :create

  private
  
  def set_position
    self.position = self.list.cards.length === 1 ? 'c' : midstring(self.list.cards[-2].position, '')
  end
end
