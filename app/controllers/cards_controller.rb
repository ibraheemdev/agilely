class CardsController < ApplicationController

  def create
    list_cards = List.find(params[:list_id]).cards
    new_position = list_cards.count.zero? ? 'c' : midstring(list_cards.last.position, '') 
    card = list_cards.create(card_params.merge(position: new_position))
    render json: { card: card }.as_json
  end

  def reorder
    new_position = { position: midstring(params[:above], params[:below]) }
    new_position.merge!(list_id: params[:new_list]) if params[:new_list]
    Card.find(params[:id]).update(new_position)
    head :no_content
  end

  def destroy
    Card.find(params[:id]).destroy
    head :no_content
  end

  private

  def card_params
    params.require(:card).permit(:title)
  end
end
