# frozen_string_literal: true

FactoryBot.define do
  factory :room do
    uuid { 'MyString' }
    name { 'MyString' }
    intro { 'MyText' }
    creator_id { 1 }
    logo { 'MyString' }
    deleted_at { '2019-05-27 13:45:33' }
  end
end
