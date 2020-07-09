class List
  include ApplicationDocument
  include Midstring
  
  field :title, type: String
  validates :title, presence: true, length: { maximum: 512 }

  field :position, type: String
  validates :position, presence: true

  belongs_to :board
  has_many :cards, 
    order: :position.asc, 
    dependent: :delete_all, 
    before_add: :set_board_on_card

  before_validation :set_position, on: :create

  def self.full(board_id)
    FullListsQuery.execute(board_id)
  end

  private

  def set_board_on_card(card)
    card.board = self.board
  end
  
  def set_position
    self.position = self.board.lists.length === 1  ? 'c' : midstring(self.board.lists[-2].position, '')
  end
end
