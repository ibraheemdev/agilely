module JsonResponse
  extend ActiveSupport::Concern

  def json_response(params)
    @type = params[:type]
    @resource = params[:resource]
    if @resource.valid?
      render json: { @type => @resource }.as_json
    else
      render json: { errors: @resource.errors.full_messages }
    end
  end
end
