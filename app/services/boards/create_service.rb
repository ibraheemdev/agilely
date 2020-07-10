module Boards
  class CreateService < ApplicationService
    attr_reader :board_params, :user
    hash_initializer :board_params, :user

    def execute
      board = Board.create(board_params)
      user.participations.create(participant: board, role: "admin")
      board
    end
  end
end
