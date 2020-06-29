FactoryBot.define do
  factory :list do
    title  { "a list" }

    after(:create) do |list|
      list.update(title: "the #{list.id.ordinalize} list")
    end
  end
end