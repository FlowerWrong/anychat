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
    avatar: Faker::Avatar.image(i, '50x50', 'jpg'),
    first_login_at: Time.now,
    last_active_at: Time.now
  )
  user.save!
end
