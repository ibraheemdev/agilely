require 'rails_helper'
require 'support/pundit_matcher'

RSpec.describe ListPolicy do
  subject { ListPolicy.new(user, list) }
  let(:board) { create(:board) }
  let(:list) { create(:list, board_id: board.id) }
  let(:user) { create(:user) }
  let!(:participation) { user.participations.create!(participant: board, role: "admin") }

  actions = [:create, :update, :destroy]
  
  actions.each do |action|
    describe "#{action.to_s}?" do
      context "when user has edit access for board" do
        it { is_expected.to permit(action) }
      end
  
      context "when user has read only access for board" do
        before { participation.update(role: "viewer") }
        it { is_expected.not_to permit(action) }
      end
    end
  end
end