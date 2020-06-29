require 'rails_helper'

RSpec.describe Board, type: :model do
  let(:board) { create(:board) }
  let!(:user) { create(:user) }
  let!(:participation) { board.participations.create(user_id: user.id, role: "admin") }
  let!(:list) { create(:list, board_id: board.id) }
  let!(:card) { create_list(:card, 6, list_id: list.id) }

  describe "attributes" do
    it { expect(board).to respond_to(:slug) }
    it { expect(board).to respond_to(:public) }
    it { expect(board).to respond_to(:title) }
    it { expect(board).to respond_to(:created_at) }
    it { expect(board).to respond_to(:updated_at) }
  end

  describe 'validations' do
    it { expect(board).to validate_presence_of(:slug) }
    it { expect(board).to validate_length_of(:slug).is_equal_to(8) }
    it { expect(board.slug).to match(/^[a-zA-Z0-9]*$/) }
    it { expect(board).to validate_presence_of(:title) }
    it { expect(board).to validate_length_of(:title).is_at_most(512) }
    it { expect(board).to validate_inclusion_of(:public).in_array([true, false]) }
  end

  describe "associations" do
    it { expect(board).to have_many(:participations).dependent(:destroy) }
    it { expect(board).to have_many(:lists).dependent(:destroy).order(position: :asc) }
  end

  describe "#to_param" do
    it { expect(board.to_param).to eq(board.slug) }
  end

  describe "#full_json" do
    it { expect(board.full_json["lists"]).to eq(board.lists.as_json(include: { cards: {} })) }
    it { expect(board.full_json["lists"][0]["cards"]).to eq(board.lists.first.cards.as_json) }
    it { expect(board.full_json["participations"]).to eq(board.participations.as_json) }
  end

  describe ".titles" do
    it { expect(Board.where(id: board.id).titles.first.as_json).to eq({"id"=> nil, "title"=> board.title, "slug"=> board.slug}) }
  end

  describe "#set_slug" do
    let(:board) { create(:board, slug: "abcdefgh") }
    let(:board2) { create(:board, slug: "abcdefgh") }
    it { expect(board.slug).not_to eq(board2.slug) }
  end
end
