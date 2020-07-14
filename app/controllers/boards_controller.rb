class BoardsController < ApplicationController
  skip_before_action :authenticate_user!, only: [:show]

  def create
    board = Boards::CreateService.execute(board_params: board_params, user: current_user)
    redirect_to show_board_path(slug: board.slug, data: { turbolinks: false })
  end

  def show
    board = authorize_board(params[:slug])
    @react_props =
      board.full.merge(
        boards_titles: current_user&.boards_titles,
        current_user: current_user,
        role: current_user&.role_in(board) || "guest"
      )
    render 'boards/show/index'
  end

  def update
    board = authorize_board(params[:slug])
    board.update(board_params)
    json_response(board)
  end

  def destroy
    board = authorize_board(params[:slug])
    board.destroy
    json_response(board)
  end

  private

  def board_params
    params.require(:board).permit(:title, :public)
  end

  def authorize_board(slug)
    authorize Board.find_by(slug: slug)
  end
end
