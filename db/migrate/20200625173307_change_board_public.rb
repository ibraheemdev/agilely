class ChangeBoardPublic < ActiveRecord::Migration[6.0]
  def change
    change_column_default(:boards, :public, nil)
  end
end
