class CreateCompanies < ActiveRecord::Migration[6.0]
  def change
    create_table :companies, comment: '公司' do |t|
      t.string :uuid, comment: 'uuid'
      t.string :name, comment: '注册名'
      t.string :alias_name, comment: '别名'
      t.text :intro, comment: '介绍'
      t.string :legal_person, comment: '法人'
      t.string :tel, comment: '电话'
      t.string :website, comment: '官网'
      t.string :email, comment: '邮箱'
      t.string :address, comment: '地址'
      t.datetime :established_at, comment: '成立日期'
      t.string :unified_social_credit_code, comment: '统一社会信用代码'
      t.string :logo, comment: 'logo'
      t.datetime :deleted_at, comment: '删除时间'

      t.timestamps
    end
  end
end
