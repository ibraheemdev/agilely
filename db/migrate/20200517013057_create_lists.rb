class CreateLists < ActiveRecord::Migration[6.0]
  def change
    create_table :lists do |t|
      t.references :board, null: false, foreign_key: true
      t.string :title

      t.timestamps
    end
  end
end
