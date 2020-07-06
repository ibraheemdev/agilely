class List
  include ApplicationDocument
  include Midstring
  
  field :title, type: String
  validates :title, presence: true, length: { maximum: 512 }

  field :position, type: String
  validates :position, presence: true

  embedded_in :board
  has_many :cards, order: :position.asc, dependent: :delete_all

  before_validation :set_position, on: :create

  private
  
  def set_position
    self.position = self.board.lists.length === 1  ? 'c' : midstring(self.board.lists[-2].position, '')
  end
end
