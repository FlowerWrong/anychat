# frozen_string_literal: true

class CreateRoomMessages < ActiveRecord::Migration[6.0]
  def change
    create_table :room_messages, comment: '群消息' do |t|
      t.string :uuid, comment: 'uuid'
      t.integer :from, comment: '发送人'
      t.integer :room_id, comment: '群'
      t.text :content, comment: '内容'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
