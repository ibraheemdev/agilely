require 'rails_helper'

RSpec.describe Participation, type: :model do
  let(:user) { create(:user) }
  let(:board) { create(:board) }
  let!(:participation) { board.participations.create!(user_id: user.id, role: "admin") }

  describe "attributes" do
    it { expect(participation).to respond_to(:role) }
    it { expect(participation).to respond_to(:participant_type) }
    it { expect(participation).to respond_to(:participant_id) }
    it { expect(participation).to respond_to(:user_id) }
    it { expect(participation).to respond_to(:created_at) }
    it { expect(participation).to respond_to(:updated_at) }
  end

  describe 'validations' do
    it { expect(participation).to validate_uniqueness_of(:participant_id).scoped_to([:participant_type, :user_id]) }
    it { expect(participation).to validate_presence_of(:role) }
    
    it "is expected that user cannot participate in the same board twice" do
      expect(board.participations.create(user_id: user.id, role: "admin") ).to be_invalid
    end

    it "is expected that user can participate in multiple boards" do
      board2 = create(:board)
      expect(board2.participations.create!(user_id: user.id, role: "admin")).to be_valid
    end
  end

  describe "associations" do
    it { expect(participation).to belong_to(:user) }
    it { expect(participation).to belong_to(:participant) }
  end

  describe ".has_participation_in?" do
    it { expect(described_class.has_participation_in? board).to be true }
  end

  describe ".participation_in" do
    it { expect(described_class.participation_in board).to eq(participation) }
  end

  describe ".role_in" do
    it { expect(described_class.role_in board).to eq(participation.role) }
  end

  describe "#can_edit?" do
    it "expects that a admin can edit" do
      expect(participation.can_edit?).to be true
    end

    it "expects that a editor can edit" do
      participation.update(role: "editor")
      expect(participation.can_edit?).to be true
    end

    it "expects that a viewer can't edit" do
      participation.update(role: "viewer")
      expect(participation.can_edit?).to be false
    end
  end
end
