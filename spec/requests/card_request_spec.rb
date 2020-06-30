require "rails_helper"

RSpec.describe "Card", type: :request do
  let!(:user) { create(:user) }
  let!(:board) { create(:board) }
  let!(:list) { create(:list, board_id: board.id) }
  let!(:card) { create(:card, list_id: list.id) }

  describe "#create" do
    context "user is signed in" do
      before { sign_in user }
      context "user has edit access" do
        let!(:participation) { board.participations.create(user_id: user.id, role: "admin") }
        before { post list_cards_path(list_id: list.id), params: { card: { title: "a new card" } } }
    
        it "returns 200" do
          expect(response).to have_http_status(200)
        end
    
        it "returns the card" do
          expect(json_response["card"]).not_to be_blank
        end
      end
    
      context "user has read only access" do
        let!(:participation) { board.participations.create(user_id: user.id, role: "viewer") }
    
        it "renders the 404 page" do
          post list_cards_path(list_id: list.id), params: { card: { title: "a new card" } }
          expect(response).to have_http_status(404)
          expect(response.body).to include("The requested page was not found")
        end
      end
    end
    
    context "user is not signed in" do
      it "redirects to login" do
        post list_cards_path(list_id: list.id), params: { card: { title: "a new card" } }
        expect(response).to redirect_to(new_user_session_path)
      end
    end
  end

  describe "#update" do
    context "user is signed in" do
      before { sign_in user }
      context "user has edit access" do
        let!(:participation) { board.participations.create(user_id: user.id, role: "admin") }
        before { put card_path(id: card.id), params: { card: { title: "a new card title" } } }

        it "returns 200" do
          expect(response).to have_http_status(200)
        end

        it "returns the card" do
          expect(json_response["card"]["title"]).to eq("a new card title")
        end
      end

      context "user has read only access" do
        let!(:participation) { board.participations.create(user_id: user.id, role: "viewer") }

        it "renders the 404 page" do
          put card_path(id: card.id), params: { card: { title: "a new card title" } }
          expect(response).to have_http_status(404)
          expect(response.body).to include("The requested page was not found")
        end
      end

      context "user moves card to unauthorized list" do
        let!(:board2) { create(:board) }
        let!(:list2) { create(:list, board_id: board2.id) }
        before { put card_path(id: card.id), params: { card: { title: "a new card title", list_id: list2.id } } }

        it "renders the 404 page" do
          expect(response.body).to include("The requested page was not found")
        end
      end
    end

    context "user is not signed in" do
      it "redirects to login" do
        put card_path(id: card.id)
        expect(response).to redirect_to(new_user_session_path)
      end
    end
  end

  describe "#destroy" do
    context "user is signed in" do
      before { sign_in user }
      context "user has edit access" do
        let!(:participation) { board.participations.create(user_id: user.id, role: "admin") }
        before { delete card_path(id: card.id) }
    
        it "returns 200" do
          expect(response).to have_http_status(200)
        end
    
        it "returns the card" do
          expect(json_response["card"]).to eq(card.as_json)
        end
      end
    
      context "user has read only access" do
        let!(:participation) { board.participations.create(user_id: user.id, role: "viewer") }
    
        it "renders the 404 page" do
          delete card_path(id: card.id)
          expect(response).to have_http_status(404)
          expect(response.body).to include("The requested page was not found")
        end
      end
    end
    
    context "user is not signed in" do
      it "redirects to login" do
        delete card_path(id: card.id)
        expect(response).to redirect_to(new_user_session_path)
      end
    end
  end
end
