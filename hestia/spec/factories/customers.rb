FactoryBot.define do
  factory :customer do
    name { "MyString" }
    mobile { "MyString" }
    email { "MyString" }
    avatar { "MyString" }
    note { "MyString" }
    ua { "MyString" }
    ip { "MyString" }
    os { "MyString" }
    browser { "MyString" }
    location { "" }
    first_login_at { "2019-04-29 15:45:23" }
    last_active_at { "2019-04-29 15:45:23" }
    deleted_at { "2019-04-29 15:45:23" }
  end
end
