class User < ApplicationDocument
  devise :database_authenticatable, :registerable,
         :recoverable, :rememberable, :validatable,
         :confirmable

  field :email,              type: String, default: ""
  field :encrypted_password, type: String, default: ""

  field :reset_password_token,   type: String
  field :reset_password_sent_at, type: Time

  field :remember_created_at, type: Time

  field :confirmation_token,   type: String
  field :confirmed_at,         type: Time
  field :confirmation_sent_at, type: Time
  field :unconfirmed_email,    type: String

  field :name, type: String
  validates :name, presence: true, length: { maximum: 50 }

  field :admin, type: Boolean

  embeds_many :participations

  def boards
    Board.in(id: participations.where(participant_type: "Board").pluck(:participant_id))
  end

  index({ confirmation_token: 1 }, { unique: true, name: "index_users_on_confirmation_token" })
  index({ reset_password_token: 1 }, { unique: true, name: "index_users_on_reset_password_token" })
  index({ email: 1 }, { unique: true, name: "index_users_on_email" })

  delegate :role_in, :participation_in, :has_participation_in?, to: :participations
  delegate :titles, to: :boards, prefix: true

  def can_edit?(record)
    participation_in(record)&.can_edit? || false
  end
end

