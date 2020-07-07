class BoardPolicy < ApplicationPolicy
  
  def create?
    true
  end

  def show?
    user&.admin? ||
    record["public"] || 
    # the board json object contains the user
    record["users"].select { |b_user| b_user["_id"] === user._id }.any? ||
    false
  end

  def update?
    user&.can_edit?(record) || false
  end
  
  def destroy?
    user&.can_edit?(record) || false
  end
end