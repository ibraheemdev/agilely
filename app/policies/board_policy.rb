class BoardPolicy < ApplicationPolicy
  
  def create?
    true
  end

  def show?
    user&.admin? ||
    record.public? || 
    user&.has_participation_in?(record) ||
    false
  end

  def update?
    user&.can_edit?(record) || false
  end
  
  def destroy?
    user&.can_edit?(record) || false
  end
end