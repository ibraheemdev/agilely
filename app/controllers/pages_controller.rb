class PagesController < ApplicationController
  before_action :authenticate_user!, only: [:app]

  def home
  end

  def pricing
  end

  def dashboard
    @boards = current_user.boards
  end
end
