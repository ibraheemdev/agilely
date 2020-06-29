class BoardPolicy < ApplicationPolicy
  def show?
    user&.admin? ||
    record.public? || 
    user&.has_participation_in?(record)
  end
end