# This migration comes from kuroko2 (originally 29)
class AddExecutionIdToProcessSignals < ActiveRecord::Migration[5.0]
  def change
    add_reference :process_signals, :execution, foreign_key: false
  end
end
