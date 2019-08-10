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

ActiveRecord::Schema.define(version: 2019_05_28_050507) do

  create_table "chat_messages", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", comment: "私聊消息", force: :cascade do |t|
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

  create_table "login_logs", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.integer "user_id", comment: "用户"
    t.string "ua", comment: "user agent"
    t.string "ip", comment: "IP地址"
    t.string "lan_ip", comment: "LAN IP地址"
    t.string "os", comment: "操作系统"
    t.string "browser", comment: "浏览器"
    t.decimal "latitude", precision: 20, scale: 17, comment: "纬度"
    t.decimal "longitude", precision: 20, scale: 17, comment: "经度"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "room_messages", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", comment: "群消息", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.integer "from", comment: "发送人"
    t.integer "room_id", comment: "群"
    t.text "content", comment: "内容"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "room_users", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", comment: "群用户", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.integer "user_id"
    t.integer "room_id"
    t.string "nickname", comment: "群昵称"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "rooms", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", comment: "群", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.string "name", comment: "群名称"
    t.text "intro", comment: "群介绍"
    t.integer "creator_id", comment: "创建者"
    t.string "logo", comment: "群logo"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "user_room_messages", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", comment: "群人员消息", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.integer "user_id", comment: "群人员"
    t.integer "room_message_id", comment: "群消息"
    t.datetime "read_at", comment: "已读时间"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "users", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", comment: "用户", force: :cascade do |t|
    t.string "uuid", comment: "uuid"
    t.string "username", comment: "用户名"
    t.string "password_digest", comment: "密码"
    t.string "mobile", comment: "手机"
    t.string "email", comment: "邮箱"
    t.string "avatar", comment: "头像"
    t.string "note", comment: "备注"
    t.datetime "first_login_at", comment: "第一次登录时间"
    t.datetime "last_active_at", comment: "最后一次活跃时间"
    t.datetime "deleted_at", comment: "删除时间"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

end
