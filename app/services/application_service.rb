class ApplicationService
  def self.execute(*args, &block)
    new(*args, &block).execute
  end
  
  def self.hash_initializer(*attr_names)
    define_method(:initialize) do |*args|
      data = args.first || {}
      attr_names.each do |attr_name|
        instance_variable_set "@#{attr_name}", data[attr_name]
      end
    end
  end
end