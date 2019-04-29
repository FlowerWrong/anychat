FactoryBot.define do
  factory :chat_message do
    from { 1 }
    to { 1 }
    content { "MyText" }
    read_at { "2019-04-29 16:32:26" }
    deleted_at { "2019-04-29 16:32:26" }
  end
end
