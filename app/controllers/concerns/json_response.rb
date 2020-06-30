module JsonResponse
  extend ActiveSupport::Concern

  def json_response(resource)
    if resource.valid?
      render json: { resource.class.name.downcase => resource }.as_json
    else
      render json: { errors: resource.errors.full_messages }
    end
  end
end
