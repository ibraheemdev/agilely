module Boards
  class CreateService < ApplicationService
    def initialize(params)
      @board_params = params[:board_params]
      @user = params[:user]
    end

    def execute
      board = Board.create(@board_params)
      @user.participations.create(participant: board, role: "admin")
      board
    end
  end
end
