# frozen_string_literal: true

class CreateRooms < ActiveRecord::Migration[6.0]
  def change
    create_table :rooms, comment: '群' do |t|
      t.string :uuid, comment: 'uuid'
      t.string :name, comment: '群名称'
      t.text :intro, comment: '群介绍'
      t.integer :creator_id, comment: '创建者'
      t.string :logo, comment: '群logo'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
