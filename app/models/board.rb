class Board < ApplicationRecord
  has_many :participations, as: :participant, dependent: :destroy
  has_many :lists, -> { order(position: :asc) }, dependent: :destroy

  validates :title, presence: true, length: { maximum: 100 }
  before_create :set_slug

  def to_param
    slug
  end

  private

  def set_slug
    loop do
      self.slug = SecureRandom.alphanumeric(8)
      break unless Board.where(slug: slug).exists?
    end
  end
end
