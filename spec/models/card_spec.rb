require 'rails_helper'

RSpec.describe Card, type: :model do
  let(:board) { create(:board) }
  let(:list) { create(:list, board_id: board.id) }
  let!(:card) { create(:card, list_id: list.id) }

  describe "attributes" do
    it { expect(card).to respond_to(:list_id) }
    it { expect(card).to respond_to(:description) }
    it { expect(card).to respond_to(:title) }
    it { expect(card).to respond_to(:position) }
    it { expect(card).to respond_to(:created_at) }
    it { expect(card).to respond_to(:updated_at) }
  end

  describe 'validations' do
    it { expect(card).to validate_presence_of(:title) }
    it { expect(card).to validate_presence_of(:position) }
  end

  describe "associations" do
    it { expect(card).to belong_to(:list) }
  end

  describe "#set_position" do
    let!(:card2) { create(:card, list_id: list.id) }
    
    it "sets first card position to 'c'" do
      expect(card.position).to eq('c')
    end

    it "sets next list positions to midstring" do
      expect(card2.position).to eq('o')
    end
  end

  describe "#board" do
    it { expect(card.board).to eq(board) }
  end
end