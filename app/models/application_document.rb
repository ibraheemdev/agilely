module ApplicationDocument
  extend ActiveSupport::Concern
  
  include Mongoid::Document
  include Mongoid::Timestamps
end