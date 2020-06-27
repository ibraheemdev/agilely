class Board < ApplicationRecord
  has_many :participations, as: :participant, dependent: :destroy
  has_many :lists, -> { order(position: :asc) }, dependent: :destroy

  validates :slug, presence: true, length: { is: 8 }, uniqueness: true
  validates :title, presence: true, length: { maximum: 512 }
  validates :public, presence: true
  
  before_validation :set_slug, on: :create

  class << self
    def full(slug)
      includes(:participations, lists: [:cards]).find_by!(slug: slug)
    end

    def titles
      select(:title, :slug)
    end
  end
  
  def role_of(user)
    user.try(:participations).try(:role_in, self) || 'guest'
  end

  def to_param
    slug
  end

  def full_json
    as_json(include: { lists: { include: { cards: {} } }, participations: {} })
  end

  private

  def set_slug
    loop do
      self.slug = SecureRandom.alphanumeric(8)
      break unless Board.where(slug: slug).exists?
    end
  end
end
