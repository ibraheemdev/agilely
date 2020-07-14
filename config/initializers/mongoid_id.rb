# override ObjectId so that model.as_json returns:
# { _id: "an_alphanumeric_string", ...} instead of
# { _id: $oid: {...}, ...}

module BSON
  class ObjectId
    def to_json(*args)
      to_s.to_json
    end

    def as_json(*args)
      to_s.as_json
    end
  end
end