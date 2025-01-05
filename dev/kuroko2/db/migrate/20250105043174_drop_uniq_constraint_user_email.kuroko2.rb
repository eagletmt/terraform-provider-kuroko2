# This migration comes from kuroko2 (originally 27)
class DropUniqConstraintUserEmail < ActiveRecord::Migration[5.0]
  def up
    remove_index :users, name: "email"
    add_index :users, :email, name: "email", using: :btree
  end
end
