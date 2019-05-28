# frozen_string_literal: true

class CreateLoginLogs < ActiveRecord::Migration[6.0]
  def change
    create_table :login_logs do |t|
      t.integer :user_id, comment: '用户'
      t.string :ua, comment: 'user agent'
      t.string :ip, comment: 'IP地址'
      t.string :lan_ip, comment: 'LAN IP地址'
      t.string :os, comment: '操作系统'
      t.string :browser, comment: '浏览器'
      t.decimal :latitude, precision: 20, scale: 17, comment: '纬度'
      t.decimal :longitude, precision: 20, scale: 17, comment: '经度'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
