class User < ApplicationRecord
  has_many :participations
  has_many :boards, through: :participations, source: :participant, source_type: "Board", dependent: :destroy

  devise :database_authenticatable, :registerable,
         :recoverable, :rememberable, :validatable,
         :confirmable

  validates :name, presence: true, length: { maximum: 50 }

  def board_titles
    boards.titles
  end

  def has_participation_in?(record)
    participations.has_participation_in?(record)
  end

  def role_in(record)
    participations.role_in(record)
  end
end
