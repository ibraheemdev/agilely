class Card < ApplicationRecord
  include Midstring
  
  belongs_to :list
  
  validates :title, presence: true
  validates :position, presence: true

  before_validation :set_position, on: :create

  private
  
  def set_position
    self.position = list.cards.blank? ? 'c' : midstring(list.cards.last.position, '')
  end
end
