require "mongo_beautiful_logger"

Mongoid.logger = MongoBeautifulLogger.new if Rails.env.development?