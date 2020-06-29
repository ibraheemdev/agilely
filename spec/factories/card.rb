FactoryBot.define do
  factory :card do
    title  { "a card" }
    
    after(:create) do |card|
      card.update(title: "the #{card.id.ordinalize} card")
    end
  end
end