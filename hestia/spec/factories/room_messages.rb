FactoryBot.define do
  factory :room_message do
    uuid { "MyString" }
    from { 1 }
    room_id { 1 }
    content { "MyText" }
    deleted_at { "2019-05-27 13:45:43" }
  end
end
