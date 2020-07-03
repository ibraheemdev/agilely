class User
  include ApplicationDocument
  
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

  field :admin, type: Boolean, default: false

  embeds_many :participations

  index({ confirmation_token: 1 }, { unique: true, name: "index_users_on_confirmation_token" })
  index({ reset_password_token: 1 }, { unique: true, name: "index_users_on_reset_password_token" })
  index({ email: 1 }, { unique: true, name: "index_users_on_email" })

  delegate :titles, to: :boards, prefix: true

  def boards
    Board.in(id: participations.where(participant_type: "Board").pluck(:participant_id))
  end

  def has_participation_in?(record)
    participations.where(participant: record).exists?
  end

  def participation_in(record)
    participations.find_by(participant: record)
  end

  def role_in(record)
    participations.find_by(participant: record)&.role || :guest
  end

  def can_edit?(record)
    participation = participation_in(record)
    participation&.role === :admin || participation&.role === :editor
  end
end

