class MongoidLogger
  
  WHITE     = "\e[37m"
  CYAN      = "\e[36m"
  MAGENTA   = "\e[35m"
  BLUE      = "\e[34m"
  YELLOW    = "\e[33m"
  GREEN     = "\e[32m"
  RED       = "\e[31m"
  BLACK     = "\e[30m"
  BOLD      = "\e[1m"
  CLEAR     = "\e[0m"

  FIND      = { regex: /STARTED.*find"=>/, color: BLUE }
  UPDATE    = { regex: /STARTED.*update"=>/, color: YELLOW }
  INSERT    = { regex: /STARTED.*insert"=>/, color: GREEN }
  DELETE    = { regex: /STARTED.*delete"=>/, color: RED }
  AGGREGATE = { regex: /STARTED.*aggregate"=>/, color: MAGENTA }
  SUCCEEDED = { regex: /SUCCEEDED/, color: GREEN }
  FAILED    = { regex: /FAILED/, color: RED }

  ACTIONS   = [ FIND, UPDATE, INSERT, DELETE, AGGREGATE, SUCCEEDED, FAILED ]

  def initialize(logger)
    @logger = logger
  end
  
  def colorize_message(msg)
    for i in 0 ... ACTIONS.size
      msg = color(msg, ACTIONS[i][:color]) if msg.match?(ACTIONS[i][:regex])
    end
    msg = "#{msg}\n"
  end

  def remove_unnecessary(msg)
    if msg.include?("opology") || msg.include?("Server description")
      ""
    else
      msg
      .sub("agilely_development.", "")
      .sub(" localhost:27017 #1", "")
      .sub(" localhost:27017", "")
      .sub(/ \| \[\d+\]/, "")
    end
  end

  %w[debug info warn error fatal unknown].each_with_index do |method, severity|
    define_method(method) do |message = nil, &block|
      message = block.call if message.nil? and block
      message = send(:remove_unnecessary, message.to_s)
      message = send(:colorize_message, message.to_s)
      @logger.add(severity, nil, message, &block)
    end
  end

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
end

if Rails.env.development?
  logger = Logger.new($stdout)
  logger.formatter = proc do |severity, datetime, progname, msg|
    "#{msg}"
  end
  Mongoid.logger = MongoidLogger.new(logger)
end
