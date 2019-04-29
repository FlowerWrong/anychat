# setting
Faker::Config.locale = 'zh-CN'

# company seed
1.upto(100).each do |i|
  company = Company.new(
    id: i,
    name: Faker::Company.name,
    alias_name: '',
    intro: Faker::Company.bs,
    legal_person: Faker::Name.name,
    tel: Faker::Company.ein,
    website: Faker::Internet.url,
    email: Faker::Internet.email,
    address: "#{Faker::Address.country} #{Faker::Address.city} #{Faker::Address.street_address}",
    established_at: Time.now,
    unified_social_credit_code: Faker::Company.duns_number,
    logo: Faker::Company.logo
  )
  company.save!

  1.upto(3).each do |j|
    user = User.new(
      username: Faker::Name.unique.name,
      password: '123456',
      password_confirmation: '123456',
      mobile: Faker::PhoneNumber.cell_phone.gsub(/\D/, ''),
      email: Faker::Internet.email,
      avatar: Faker::Avatar.image(j, '50x50', 'jpg'),
      company_id: company.id,
      first_login_at: Time.now,
      last_active_at: Time.now
    )
    user.save!
  end
end
