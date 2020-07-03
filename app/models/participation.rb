class Participation
  include ApplicationDocument
  
  field :role, type: Symbol
  validates :role, presence: true, inclusion: { in: [:viewer, :editor, :admin] }

  belongs_to :participant, polymorphic: true
  validates :participant_id, uniqueness: { scope: [:participant_type, :user_id] }

  embedded_in :user
end
