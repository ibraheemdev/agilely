class ListsController < ApplicationController

  def create
    list = Board.find_by(slug: params[:board_slug]).lists.create(list_params)
    json_response({type: "list", resource: list})
  end

  def update
    List.find(params[:id]).update(list_params)
    head :no_content
  end

  def destroy
    List.find(params[:id]).destroy
    head :no_content
  end

  private

  def list_params
    params.require(:list).permit(:title, :position)
  end
end
