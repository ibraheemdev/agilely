class BoardsController < ApplicationController

  def create
    board = Board.create(board_params)
    current_user.participations.create(participant: board)
    redirect_to show_board_path(slug: board.slug, data: { turbolinks: false })
  end

  def show
    @board = Board.includes(:participations, lists: [:cards]).find_by(slug: params[:slug])
  end

  def update
    Board.find_by(slug: params[:slug]).update(board_params)
  end

  def destroy
    Board.find_by(slug: params[:slug]).destroy
  end

  private

  def board_params
    params.require(:board).permit(:title)
  end
end
