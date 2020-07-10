require 'rails_helper'

RSpec.describe Board, type: :model do
  subject { create(:board) }
  let!(:user) { create(:user) }
  let!(:participation) { user.participations.create(participant: subject, role: "admin") }
  let!(:list) { create(:list, board_id: subject.id) }
  let!(:card) { create_list(:card, 1, list_id: list.id, board_id: subject.id) }

  describe "attributes" do
    it { is_expected.to be_mongoid_document }
    it { is_expected.to have_timestamps }

    it { is_expected.to have_field(:title).of_type(String) }
    it { is_expected.to validate_presence_of(:title) }
    it { is_expected.to validate_length_of(:title).with_maximum(512) }

    it { is_expected.to have_field(:slug).of_type(String) }
    it { is_expected.to validate_presence_of(:slug) }
    it { is_expected.to validate_length_of(:slug).is(8) }

    it { is_expected.to have_field(:public).of_type(Mongoid::Boolean) }
    it { is_expected.to validate_inclusion_of(:public).to_allow([true, false]) }
    
    it { is_expected.to have_many(:lists).ordered_by(:position.asc) }
    it { is_expected.to have_many(:cards).ordered_by(:position.asc) }
  end

  describe "#users" do
    let!(:user2) { create(:user, email: "idk@idk.com") }
    it { expect(subject.users).to include(user) }
    it { expect(subject.users).not_to include(user2) }
  end

  describe "#full" do
    it { expect(subject.full).to eq(FullBoardQuery.execute(subject)) }
  end

  describe "#to_param" do
    it { expect(subject.to_param).to eq(subject.slug) }
  end

  describe ".titles" do
    it { expect(Board.titles).to eq([[subject.title, subject.slug]]) }
  end
end
