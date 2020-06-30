class CardsController < ApplicationController

  def create
    list = authorize List.find(params[:list_id])
    card = list.cards.create(card_params)
    json_response({type: "card", resource: card})
  end

  def update
    card = authorize Card.find(params[:id])
    card.update(card_params)
    head :no_content
  end

  def destroy
    card = authorize Card.find(params[:id])
    card.destroy
    head :no_content
  end

  private

  def card_params
    params.require(:card).permit(:title, :position, :list_id)
  end
end
