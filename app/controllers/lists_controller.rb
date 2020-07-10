class ListsController < ApplicationController

  def create
    board = authorize Board.find_by(slug: params[:board_slug]), :update?, policy_class: BoardPolicy
    list  = board.lists.create(list_params)
    json_response(list)
  end

  def update
    list = authorize List.find(params[:id])
    list.update(list_params)
    json_response(list)
  end

  def destroy
    list = authorize List.find(params[:id])
    list.destroy
    head :no_content
  end

  private

  def list_params
    params.require(:list).permit(:title, :position)
  end

  def authorize_board(params)
    authorize Board.find_by(params), :update?, policy_class: BoardPolicy
  end
end
