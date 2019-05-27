# frozen_string_literal: true

class CreateRoomUsers < ActiveRecord::Migration[6.0]
  def change
    create_table :room_users, comment: '群用户' do |t|
      t.string :uuid, comment: 'uuid'
      t.integer :user_id
      t.integer :room_id
      t.string :nickname, comment: '群昵称'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
