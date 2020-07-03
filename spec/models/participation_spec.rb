require 'rails_helper'

RSpec.describe Participation, type: :model do
  let(:user) { create(:user) }
  let(:board) { create(:board) }
  
  subject { user.participations.create!(participant: board, role: "admin") }

  describe "attributes" do
    it { is_expected.to be_mongoid_document }
    it { is_expected.to have_timestamps }

    it { is_expected.to have_field(:role).of_type(Symbol) }
    it { is_expected.to validate_presence_of(:role) }
    it { is_expected.to validate_inclusion_of(:role).to_allow(:viewer, :editor, :admin) }

    it { is_expected.to be_embedded_in(:user) }

    it { is_expected.to belong_to(:participant) }
    it { is_expected.to validate_uniqueness_of(:participant_id).scoped_to([:participant_type, :user_id]) }
    
    it "is expected that user can't be on the same board twice" do
      user.participations.create(participant: board, role: "admin")
      expect(user.participations.create(participant: board, role: "admin") ).to be_invalid
    end
  end
end
