class BoardPolicy < ApplicationPolicy
  def show?
    user&.admin? ||
    record.public? || 
    user && user.boards.exists?(record.id)
  end
end