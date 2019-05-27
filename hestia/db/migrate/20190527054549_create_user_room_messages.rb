# frozen_string_literal: true

class CreateUserRoomMessages < ActiveRecord::Migration[6.0]
  def change
    create_table :user_room_messages, comment: '群人员消息' do |t|
      t.string :uuid, comment: 'uuid'
      t.integer :user_id, comment: '群人员'
      t.integer :room_message_id, comment: '群消息'
      t.datetime :read_at, comment: '已读时间'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
