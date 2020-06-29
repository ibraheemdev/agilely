class Participation < ApplicationRecord
  belongs_to :participant, polymorphic: true
  belongs_to :user

  validates :participant_id, uniqueness: { scope: [:participant_type, :user_id] }
  validates :role, presence: true

  enum role: [:viewer, :editor, :admin]
  
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
