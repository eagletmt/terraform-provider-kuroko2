# This migration comes from kuroko2 (originally 28)
class ChangeExitStatusToUnsignedTinyint < ActiveRecord::Migration[5.0]
  def up
    change_column :executions, :exit_status, :integer, limit: 2
  end

  def down
    change_column :executions, :exit_status, :integer, limit: 1
  end
end
