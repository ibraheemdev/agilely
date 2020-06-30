class ListsController < ApplicationController

  def create
    board = authorize Board.find_by(slug: params[:board_slug]), :update?, policy_class: BoardPolicy
    list = board.lists.create(list_params)
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
    json_response(list)
  end

  private

  def list_params
    params.require(:list).permit(:title, :position)
  end
end
