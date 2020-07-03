# This file is copied to spec/ when you run 'rails generate rspec:install'
require 'spec_helper'
ENV['RAILS_ENV'] ||= 'test'
require File.expand_path('../config/environment', __dir__)
abort("The Rails environment is running in production mode!") if Rails.env.production?
require 'rspec/rails'
require 'database_cleaner/mongoid'
require "capybara/rspec"
require 'support/json_response'
require 'mongoid-rspec'

RSpec.configure do |config|

  config.include Devise::Test::IntegrationHelpers, type: :request
  config.include Mongoid::Matchers, type: :model

  config.fixture_path = "#{::Rails.root}/spec/fixtures"
  config.use_transactional_fixtures = true

  config.before(:each, type: :system) do
    driven_by :rack_test
  end
  
  config.before(:each, type: :system, js: true) do
    driven_by :selenium_chrome_headless
  end

  config.infer_spec_type_from_file_location!
  
  config.filter_rails_from_backtrace!

  config.include FactoryBot::Syntax::Methods
  
  config.before(:suite) do
    DatabaseCleaner[:mongoid].strategy = :truncation
  end
  
  config.around(:each) do |example|
    DatabaseCleaner.cleaning do
      example.run
    end
  end
end
