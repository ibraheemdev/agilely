class CardsController < ApplicationController

  def create
    card = List.find(params[:list_id]).cards.create(card_params)
    render json: { card: card }.as_json
  end

  def update
    Card.find(params[:id]).update(card_params)
    head :no_content
  end

  def destroy
    Card.find(params[:id]).destroy
    head :no_content
  end

  private

  def card_params
    params.require(:card).permit(:title, :position, :list_id)
  end
end
