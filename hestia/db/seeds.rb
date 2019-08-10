# frozen_string_literal: true

# setting
Faker::Config.locale = 'zh-CN'

# company seed
1.upto(100).each do |i|
  user = User.new(
    uuid: SecureRandom.uuid,
    username: Faker::Name.unique.name,
    password: '123456',
    password_confirmation: '123456',
    mobile: Faker::PhoneNumber.cell_phone.gsub(/\D/, ''),
    email: Faker::Internet.email,
    avatar: Faker::Avatar.image,
    first_login_at: Time.now,
    last_active_at: Time.now
  )
  user.save!
end

user1 = User.new(
  uuid: SecureRandom.uuid,
  username: 'yang',
  password: '123456',
  password_confirmation: '123456',
  mobile: '13560474456',
  email: 'sysuyangkang@gmail.com',
  avatar: Faker::Avatar.image,
  first_login_at: Time.now,
  last_active_at: Time.now
)
user1.save!

user2 = User.new(
  uuid: SecureRandom.uuid,
  username: 'lin',
  password: '123456',
  password_confirmation: '123456',
  mobile: '13560474457',
  email: 'linzi@gmail.com',
  avatar: Faker::Avatar.image,
  first_login_at: Time.now,
  last_active_at: Time.now
)
user2.save!
