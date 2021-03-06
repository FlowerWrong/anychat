# frozen_string_literal: true

class CreateChatMessages < ActiveRecord::Migration[6.0]
  def change
    create_table :chat_messages, comment: '私聊消息' do |t|
      t.string :uuid, comment: 'uuid'
      t.integer :from, comment: '发送人'
      t.integer :to, comment: '接收人'
      t.text :content, comment: '内容'
      t.string :ack, comment: 'req ack'
      t.datetime :read_at, comment: '已读时间'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
