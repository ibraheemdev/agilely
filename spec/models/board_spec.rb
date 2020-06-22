require 'rails_helper'

RSpec.describe Board, type: :model do
  let(:board) { create(:board) }

  describe "attributes" do
    it { expect(board).to respond_to(:slug) }
    it { expect(board).to respond_to(:title) }
    it { expect(board).to respond_to(:created_at) }
    it { expect(board).to respond_to(:updated_at) }
  end

  describe 'validations' do
    it { expect(board).to validate_presence_of(:slug) }
    it { expect(board).to validate_length_of(:slug).is_equal_to(8) }
    it { expect(board).to validate_presence_of(:title) }
    it { expect(board).to validate_length_of(:title).is_at_most(512) }
    it { expect(board.slug).to match(/^[a-zA-Z0-9]*$/) }
  end

  describe "associations" do
    it { expect(board).to have_many(:participations).dependent(:destroy) }
    it { expect(board).to have_many(:lists).dependent(:destroy).order(position: :asc) }
  end
end
