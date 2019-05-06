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

ActiveRecord::Schema.define(version: 2019_04_29_083226) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"
  enable_extension "postgis"

  create_table "apps", comment: "企业应用", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.string "name", comment: "名称"
    t.integer "company_id", comment: "公司"
    t.text "intro", comment: "介绍"
    t.string "domains", comment: "域名列表", array: true
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["company_id"], name: "index_apps_on_company_id"
    t.index ["domains"], name: "index_apps_on_domains", using: :gin
  end

  create_table "chat_messages", comment: "聊天消息", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.integer "from", comment: "发送人"
    t.integer "to", comment: "接收人"
    t.text "content", comment: "内容"
    t.string "ack", comment: "req ack"
    t.datetime "read_at", comment: "已读时间"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "companies", comment: "公司", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
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

  create_table "spatial_ref_sys", primary_key: "srid", id: :integer, default: nil, force: :cascade do |t|
    t.string "auth_name", limit: 256
    t.integer "auth_srid"
    t.string "srtext", limit: 2048
    t.string "proj4text", limit: 2048
  end

  create_table "users", comment: "用户", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.string "username", comment: "用户名"
    t.string "password_digest", comment: "密码"
    t.string "mobile", comment: "手机"
    t.string "email", comment: "邮箱"
    t.string "avatar", comment: "头像"
    t.string "note", comment: "备注"
    t.string "ua", comment: "user agent"
    t.string "ip", comment: "IP地址"
    t.string "lan_ip", comment: "LAN IP地址"
    t.string "os", comment: "操作系统"
    t.string "browser", comment: "浏览器"
    t.float "latitude", comment: "纬度"
    t.float "longitude", comment: "经度"
    t.integer "company_id", comment: "公司"
    t.integer "app_id", comment: "应用"
    t.string "role", comment: "角色"
    t.datetime "first_login_at", comment: "第一次登录时间"
    t.datetime "last_active_at", comment: "最后一次活跃时间"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.index ["app_id"], name: "index_users_on_app_id"
    t.index ["company_id"], name: "index_users_on_company_id"
  end

end
