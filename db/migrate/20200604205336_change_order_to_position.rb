class ChangeOrderToPosition < ActiveRecord::Migration[6.0]
  def change
    rename_column :cards, :order, :position
  end
end
