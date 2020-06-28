class ApplicationController < ActionController::Base
  include JsonResponse
  include Pundit
  
  before_action :configure_permitted_parameters, if: :devise_controller?
  rescue_from Pundit::NotAuthorizedError, with: :render404
  rescue_from ActiveRecord::RecordNotFound, with: :render404

  private

  def render404
    render 'pages/404'
  end

  protected

  def configure_permitted_parameters
    devise_parameter_sanitizer.permit(:sign_up, keys: [:name])
    devise_parameter_sanitizer.permit(:account_update, keys: [:name])
  end

  def after_sign_in_path_for(resource)
    stored_location_for(resource) || dashboard_path
  end
end
