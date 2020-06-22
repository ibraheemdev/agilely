module Boards
  class CreateService < ApplicationService
    def initialize(params)
      @board_params = params[:board_params]
      @user = params[:user]
    end

    def execute
      board = Board.create(@board_params)
      @user.participations.create(participant: board)
      board
    end
  end
end
