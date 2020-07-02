class List < ApplicationDocument
  include Midstring
  
  field :title, type: String
  validates :title, presence: true, length: { maximum: 512 }

  field :position, type: String
  validates :position, presence: true

  belongs_to :board, index: true
  has_many :cards, order: :position.asc, dependent: :delete_all

  before_validation :set_position, on: :create

  private
  
  def set_position
    self.position = board.lists.blank? ? 'c' : midstring(board.lists.last.position, '')
  end
end
