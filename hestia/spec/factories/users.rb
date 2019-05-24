FactoryBot.define do
  factory :user do
    username { "MyString" }
    password { "MyString" }
    mobile { "MyString" }
    email { "MyString" }
    avatar { "MyString" }
    first_login_at { "2019-04-29 14:01:40" }
    last_active_at { "2019-04-29 14:01:40" }
    deleted_at { "2019-04-29 14:01:40" }
  end
end
