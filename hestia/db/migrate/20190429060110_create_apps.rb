class CreateApps < ActiveRecord::Migration[6.0]
  def change
    create_table :apps, comment: '企业应用' do |t|
      t.string :uuid, comment: 'uuid'
      t.string :name, comment: '名称'
      t.integer :company_id, comment: '公司'
      t.text :intro, comment: '介绍'
      t.string :domains, array: true, comment: '域名列表'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end

    add_index :apps, :company_id
    add_index :apps, :domains, using: 'gin'
  end
end
