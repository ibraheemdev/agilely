module ApplicationHelper
  def react_component(id, props)
    content_tag(:div, { id: id, data: { props: props } }) do
    end
  end
end