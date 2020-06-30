class CardsController < ApplicationController

  def create
    list = authorize List.find(params[:list_id]), :update?, policy_class: ListPolicy
    card = list.cards.create(card_params)
    json_response(card)
  end

  def update
    card = authorize Card.find(params[:id])
    card_params[:list_id] && ( authorize List.find(card_params[:list_id]), :update?, policy_class: ListPolicy )
    card.update(card_params)
    json_response(card)
  end

  def destroy
    card = authorize Card.find(params[:id])
    card.destroy
    json_response(card)
  end

  private

  def card_params
    params.require(:card).permit(:title, :position, :list_id)
  end
end
