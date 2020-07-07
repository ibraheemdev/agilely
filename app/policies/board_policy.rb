class BoardPolicy < ApplicationPolicy
  
  def create?
    true
  end

  def show?
    user&.admin? ||
    record["public"] || 
    user&.has_participation?("Board", record["_id"]) ||
    false
  end

  def update?
    user&.can_edit?("Board", record["_id"]) || false
  end
  
  def destroy?
    user&.can_edit?("Board", record["_id"]) || false
  end
end