class LoggerDecorator
  WHITE = "\e[37m"
  CYAN = "\e[36m"
  MAGENTA = "\e[35m"
  BLUE = "\e[34m"
  YELLOW = "\e[33m"
  GREEN = "\e[32m"
  RED = "\e[31m"
  BLACK = "\e[30m"
  BOLD = "\e[1m"
  CLEAR = "\e[0m"
  STARTED_FIND = /STARTED.*find"=>/
  STARTED_UPDATE = /STARTED.*update"=>/
  STARTED_INSERT = /STARTED.*insert"=>/
  STARTED_DELETE = /STARTED.*delete"=>/

  def initialize(logger)
    @logger = logger
  end

  %w[debug info warn error fatal unknown].each_with_index do |method, severity|
    define_method(method) do |message = nil, &block|
      message = block.call if message.nil? and block
      message = send(:colorize_message, message.to_s)
      message = send(:remove_prefix, message.to_s)
      @logger.add(severity, nil, message, &block)
    end
  end

  # Proxy everything else to the logger instance
  def respond_to?(method)
    super || @logger.respond_to?(method)
  end

  private

  def method_missing(method, *args, &block)
    @logger.send(method, *args, &block)
  end

  def color(text, color, bold=false)
    color = self.class.const_get(color.to_s.upcase) if color.is_a?(Symbol)
    bold  = bold ? BOLD : ""
    "#{bold}#{color}#{text}#{CLEAR}"
  end

  def colorize_message(message)
    message = "#{message}\n"
    message = case
    when message.match?(STARTED_FIND)
      color(message, BLUE)
    when message.match?(STARTED_UPDATE)
      color(message, YELLOW)
    when message.match?(STARTED_INSERT)
      color(message, MAGENTA)
    when message.match?(STARTED_DELETE)
      color(message, RED)
    when message.match?(/SUCCEEDED/)
      color(message, GREEN)
    when message.match?(/FAILED/)
      color(message, RED)
    when message.include?("opology") || message.include?("Server description")
      ""
    else
      message
    end
  end

  def remove_prefix(message)
    message
      .sub("agilely_development.", "")
      .sub(" localhost:27017 #1", "")
      .sub(" localhost:27017", "")
      .sub(/ \| \[\d+\]/, "")
  end
end

if Rails.env.development?
  logger = Logger.new($stdout)
  logger.formatter = proc do |severity, datetime, progname, msg|
    "#{msg}"
  end
  Mongoid.logger = LoggerDecorator.new(logger)
end