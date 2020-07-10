module Boards
  class CreateService < ApplicationService
    attr_reader :board_params, :user
    hash_initializer :board_params, :user

    def execute
      board = Board.create(board_params)
      board.participations.create(user_id: user.id, role: "admin")
      board
    end
  end
end
