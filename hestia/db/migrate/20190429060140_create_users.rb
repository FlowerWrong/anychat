# frozen_string_literal: true

class CreateUsers < ActiveRecord::Migration[6.0]
  def change
    create_table :users, comment: '用户' do |t|
      t.string :uuid, comment: 'uuid'
      t.string :username, comment: '用户名'
      t.string :password_digest, comment: '密码'
      t.string :mobile, comment: '手机'
      t.string :email, comment: '邮箱'
      t.string :avatar, comment: '头像'
      t.string :note, comment: '备注'
      t.datetime :first_login_at, comment: '第一次登录时间'
      t.datetime :last_active_at, comment: '最后一次活跃时间'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
