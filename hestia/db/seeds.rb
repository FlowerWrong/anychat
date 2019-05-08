# setting
Faker::Config.locale = 'zh-CN'

# company seed
1.upto(100).each do |i|
  company = Company.new(
    id: i,
    uuid: SecureRandom.uuid,
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

  app = App.new(
    uuid: SecureRandom.uuid,
    name: company.name,
    token: SecureRandom.urlsafe_base64,
    company_id: company.id,
    intro: company.name,
    domains: [Faker::Internet.domain_name]
  )
  app.save!

  1.upto(3).each do |j|
    role = j == 1 ? 'admin' : 'member'
    user = User.new(
      uuid: SecureRandom.uuid,
      username: Faker::Name.unique.name,
      password: '123456',
      password_confirmation: '123456',
      mobile: Faker::PhoneNumber.cell_phone.gsub(/\D/, ''),
      email: Faker::Internet.email,
      avatar: Faker::Avatar.image(j, '50x50', 'jpg'),
      company_id: company.id,
      app_id: app.id,
      role: role,
      first_login_at: Time.now,
      last_active_at: Time.now
    )
    user.save!
  end
end
