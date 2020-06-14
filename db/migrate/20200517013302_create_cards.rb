class CreateCards < ActiveRecord::Migration[6.0]
  def change
    create_table :cards do |t|
      t.string :title
      t.text :description
      t.references :list, null: false, foreign_key: true
      t.string :order

      t.timestamps
    end
  end
end
