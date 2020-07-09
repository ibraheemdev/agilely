class BoardsController < ApplicationController
  skip_before_action :authenticate_user!, only: [:show]

  def create
    board = Boards::CreateService.execute(board_params: board_params, user: current_user)
    redirect_to show_board_path(slug: board.slug, data: { turbolinks: false })
  end

  def show
    @board = authorize Board.find_by(slug: params[:slug]), policy_class: BoardPolicy
    @boards_titles = current_user&.boards_titles
    @role = current_user&.role_in(@board) || "guest"
    @board = @board.full
  end

  def update
    board = authorize Board.find_by(slug: params[:slug])
    board.update(board_params)
    json_response(board)
  end

  def destroy
    board = authorize Board.find_by(slug: params[:slug])
    board.destroy
    json_response(board)
  end

  private

  def board_params
    params.require(:board).permit(:title, :public)
  end
end
