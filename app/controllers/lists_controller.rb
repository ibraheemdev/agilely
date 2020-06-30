class ListsController < ApplicationController

  def create
    board = authorize Board.find_by(slug: params[:board_slug])
    list = board.lists.create(list_params)
    json_response(type: "list", resource: list)
  end

  def update
    list = authorize List.find(params[:id])
    list.update(list_params)
    head :no_content
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
end
