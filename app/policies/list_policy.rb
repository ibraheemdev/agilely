class ListPolicy < BoardPolicy

  def initialize(user, record)
    @user = user
    @record = record.class.name === "Board" ? record : record.board
  end
  
  def create?
    update?
  end

  def update?
    super
  end
  
  def destroy?
    update?
  end
end