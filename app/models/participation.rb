class Participation
  include ApplicationDocument
  
  field :role, type: Integer
  validates :role, presence: true
  # enum role: [:viewer, :editor, :admin]
  # enumerize :role, in: { viewer: 1, editor: 2, admin: 3 }, predicates: true, scope: :shallow

  belongs_to :participant, polymorphic: true, index: true
  validates :participant_id, uniqueness: { scope: [:participant_type, :user_id] }

  embedded_in :user
  
  def self.has_participation_in?(record)
    exists?(participant: record)
  end

  def self.participation_in(record)
    find_by(participant: record)
  end

  def self.role_in(record)
    find_by(participant: record)&.role || 'guest'
  end

  def can_edit?
    admin? || editor?
  end
end
