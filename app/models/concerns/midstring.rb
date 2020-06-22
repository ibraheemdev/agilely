module Midstring
  extend ActiveSupport::Concern
    def midstring(prev, nxt)
      pos = p = n = 0
      str = nil
      while p == n
        p = pos < prev.length ? prev[pos].ord : 96
        n = pos < nxt.length ? nxt[pos].ord : 123
        pos += 1
      end
      str = prev[0...pos - 1]
      if p == 96
        while n == 97
          n = pos < nxt.length ? nxt[pos].ord : 123
          pos += 1
          str += 'a'
        end
        if n == 98
          str += 'a'
          n = 123
        end
      elsif p + 1 == n
        str += p.chr
        n = 123
        while (p = pos < prev.length ? prev[pos].ord : 96) == 122
          pos += 1
          str += 'z'
        end
      end
      str + ((p.to_f + n) / 2).ceil.chr
    end
end
