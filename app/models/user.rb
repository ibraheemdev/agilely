class User < ApplicationRecord
  has_many :participations
  has_many :boards, through: :participations, source: :participant, source_type: "Board", dependent: :destroy

  devise :database_authenticatable, :registerable,
         :recoverable, :rememberable, :validatable,
         :confirmable

  validates :name, presence: true, length: { maximum: 50 }
  
  delegate :role_in, :participation_in, :has_participation_in?, to: :participations
  delegate :titles, to: :boards, prefix: true

  def can_edit?(record)
    participation_in(record)&.can_edit? || false
  end
end
