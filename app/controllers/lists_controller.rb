class ListsController < ApplicationController

  def create
    board = authorize_board params[:board_slug]
    list  = board.lists.create(list_params)
    json_response(list)
  end

  def update
    board = authorize_board params[:board_slug]
    list  = board.lists.find(params[:id])
    list.update(list_params)
    json_response(list)
  end

  def destroy
    board = authorize_board params[:board_slug]
    board.lists.find(params[:id]).destroy
    head :no_content
  end

  private

  def list_params
    params.require(:list).permit(:title, :position)
  end

  def authorize_board(slug)
    authorize Board.find_by(slug: slug), :update?, policy_class: BoardPolicy
  end
end
