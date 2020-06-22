FactoryBot.define do
  factory :user do
    name  { "John" }
    email { "test@test.com" }
    password { "password" }
    password_confirmation { "password" }
  end
end