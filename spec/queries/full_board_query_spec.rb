require 'rails_helper'

RSpec.describe FullBoardQuery do
  let(:board) { create(:board) }
  let!(:user) { create(:user) }
  let!(:participation) { user.participations.create(participant: board, role: "admin") }
  let!(:list) { create(:list, board_id: board.id) }
  let!(:card) { create_list(:card, 1, list_id: list.id, board_id: board.id) }
  subject(:query) { FullBoardQuery.execute(board) }

  describe "#execute" do
    it "has all expected keys" do
      board_keys.each { |k| expect(query[k]).to be_present }
      list_keys.each { |k| expect(query["lists"][0][k]).to be_present }
      card_keys.each { |k| expect(query["lists"][0]["cards"][0][k]).to be_present }
      participant_keys.each { |k| expect(query["participants"][0][k]).to be_present }
    end

    it "only contains expected keys" do
      invalid_keys.each { |k| expect(query[k]).not_to be_present }
    end
  end

  def board_keys
    ["_id", "created_at", "public", "slug", "title", "updated_at", "lists", "participants"]
  end

  def list_keys
    ["_id", "created_at", "position", "title", "updated_at", "cards"]
  end

  def card_keys
    ["_id", "created_at", "position", "title", "updated_at", "list_id", "board_id"]
  end

  def participant_keys
    ["_id", "created_at", "participant_id", "participant_type", "updated_at", "name", "email"]
  end

  def invalid_keys
    ["idk", "invalid", "fail", "somethingelse"]
  end
end
