require 'rails_helper'

RSpec.describe Participation, type: :model do
  let(:user) { create(:user) }
  let(:board) { create(:board) }
  let!(:participation) { Participation.create!(user_id: user.id, participant: board) }

  describe 'validations' do
    it { expect(participation).to validate_uniqueness_of(:participant_id).scoped_to([:participant_type, :user_id]) }
    
    it "is expected that user cannot participate in the same board twice" do
      participation2 = Participation.create(user_id: user.id, participant: board) 
      expect(participation2.valid?).to be false
    end

    it "is expected that user can participate in multiple boards" do
      board2 = create(:board)
      participation2 = Participation.create!(user_id: user.id, participant: board2)
      expect(participation2.valid?).to be true
    end
  end

  describe "associations" do
    it { expect(participation).to belong_to(:user) }
    it { expect(participation).to belong_to(:participant) }
  end
end
