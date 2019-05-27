FactoryBot.define do
  factory :room_user do
    uuid { "MyString" }
    user_id { 1 }
    room_id { 1 }
    nickname { "MyString" }
    deleted_at { "2019-05-27 13:45:38" }
  end
end
