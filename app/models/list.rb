class List < ApplicationRecord
  include Midstring
  
  belongs_to :board
  has_many :cards, -> { order(position: :asc) }, dependent: :delete_all
  validates :title, presence: true, length: { maximum: 512 }
  validates :position, presence: true

  before_validation :set_position, on: :create

  private
  
  def set_position
    self.position = board.lists.blank? ? 'c' : midstring(board.lists.last.position, '')
  end
end
