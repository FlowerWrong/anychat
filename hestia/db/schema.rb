# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `rails
# db:schema:load`. When creating a new database, `rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2019_04_29_060140) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "companies", comment: "公司", force: :cascade do |t|
    t.string "name", comment: "注册名"
    t.string "alias_name", comment: "别名"
    t.text "intro", comment: "介绍"
    t.string "legal_person", comment: "法人"
    t.string "tel", comment: "电话"
    t.string "website", comment: "官网"
    t.string "email", comment: "邮箱"
    t.string "address", comment: "地址"
    t.datetime "established_at", comment: "成立日期"
    t.string "unified_social_credit_code", comment: "统一社会信用代码"
    t.string "logo", comment: "logo"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "users", comment: "用户", force: :cascade do |t|
    t.string "username", comment: "用户名"
    t.string "password_digest", comment: "密码"
    t.string "mobile", comment: "手机"
    t.string "email", comment: "邮箱"
    t.string "avatar", comment: "头像"
    t.integer "company_id", comment: "公司"
    t.datetime "first_login_at", comment: "第一次登录时间"
    t.datetime "last_active_at", comment: "最后一次活跃时间"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

end
