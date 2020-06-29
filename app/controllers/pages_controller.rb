class PagesController < ApplicationController
  skip_before_action :authenticate_user!, only: [:home, :pricing]

  def home
  end

  def pricing
  end

  def dashboard
    @boards = current_user.boards
  end
end
