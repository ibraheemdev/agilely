FactoryBot.define do
  factory :user do
    name  { "John" }
    email { "test@test.com" }
    password { "password" }
    confirmed_at { Time.now }
  end
end