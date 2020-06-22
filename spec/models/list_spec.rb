require 'rails_helper'

RSpec.describe List, type: :model do
  let(:board) { create(:board) }
  let!(:list) { create(:list, board_id: board.id) }

  describe "attributes" do
    it { expect(list).to respond_to(:board_id) }
    it { expect(list).to respond_to(:title) }
    it { expect(list).to respond_to(:position) }
    it { expect(list).to respond_to(:created_at) }
    it { expect(list).to respond_to(:updated_at) }
  end
  
  describe 'validations' do
    it { expect(list).to validate_presence_of(:title) }
    it { expect(list).to validate_length_of(:title).is_at_most(512) }
    it { expect(list).to validate_presence_of(:position) }
  end

  describe "associations" do
    it { expect(list).to belong_to(:board) }
    it { expect(list).to have_many(:cards).dependent(:delete_all).order(position: :asc) }
  end

  describe "#set_position" do
    let!(:list2) { create(:list, board_id: board.id) }
    
    it "sets first list position to 'c'" do
      expect(list.position).to eq('c')
    end

    it "sets next list positions to midstring" do
      expect(list2.position).to eq('o')
    end
  end
end