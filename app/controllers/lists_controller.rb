class ListsController < ApplicationController

  def create
    board_lists = Board.find_by(slug: params[:board_slug]).lists
    new_position = board_lists.count.zero? ? 'c' : midstring(board_lists.last.position, '')
    list = board_lists.create(list_params.merge(position: new_position))
    render json: { list: list }.as_json(include: { cards: {} })
  end

  def reorder
    List.find(params[:id]).update(
      position: midstring(params[:above], params[:below])
    )
    head :no_content
  end

  def update
    List.find(params[:id]).update(list_params)
  end

  def destroy
    List.find(params[:id]).destroy
  end

  private

  def list_params
    params.require(:list).permit(:title)
  end
end
