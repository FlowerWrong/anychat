FactoryBot.define do
  factory :user_room_message do
    uuid { "MyString" }
    user_id { 1 }
    room_message_id { 1 }
    read_at { "2019-05-27 13:45:49" }
    deleted_at { "2019-05-27 13:45:49" }
  end
end
