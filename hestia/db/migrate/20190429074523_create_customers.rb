class CreateCustomers < ActiveRecord::Migration[6.0]
  def change
    create_table :customers, comment: '客户' do |t|
      t.string :name, comment: '姓名'
      t.string :mobile, comment: '手机'
      t.string :email, comment: '邮箱'
      t.string :avatar, comment: '头像'
      t.string :note, comment: '备注'
      t.string :ua, comment: 'user agent'
      t.string :ip, comment: 'IP地址'
      t.string :os, comment: '操作系统'
      t.string :browser, comment: '浏览器'
      t.point :location, comment: '地理位置'
      t.datetime :first_login_at, comment: '第一次登录时间'
      t.datetime :last_active_at, comment: '最后一次活跃时间'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
