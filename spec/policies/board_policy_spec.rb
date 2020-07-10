require 'rails_helper'
require 'support/pundit_matcher'

RSpec.describe BoardPolicy do
  subject { BoardPolicy.new(user, board) }
  let(:board) { create(:board, public: false) }
  let(:user) { create(:user) }

  describe "#create?" do
    it { is_expected.to permit(:create) }
  end

  describe "#show?" do
    context "when admin user" do
      before { user.admin = true }
      it { is_expected.to permit(:show) }
    end

    context "when public board" do
      before { board.public = true }
      it { is_expected.to permit(:show) }
    end

    context "when user has participation" do
      before { user.participations.create!(participant: board, role: "admin") }
      it { is_expected.to permit(:show) }
    end

    context "when user does not have participation" do
      it { is_expected.not_to permit(:show) }
    end
  end

  describe "#update?" do
    context "when user has edit access" do
      before { user.participations.create!(participant: board, role: "admin") }
      it { is_expected.to permit(:update) }
    end

    context "when user has read only access" do
      before { user.participations.create!(participant: board, role: "viewer") }
      it { is_expected.not_to permit(:update) }
    end
  end

  describe "#destroy?" do
    context "when user has edit access" do
      before { user.participations.create!(participant: board, role: "admin") }
      it { is_expected.to permit(:destroy) }
    end

    context "when user has read only access" do
      before { user.participations.create!(participant: board, role: "viewer") }
      it { is_expected.not_to permit(:destroy) }
    end
  end
end