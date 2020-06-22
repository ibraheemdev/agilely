require 'rails_helper'

RSpec.describe Midstring do
  include Midstring

  describe "#midstring" do
    context "basic cases" do
      it { expect(midstring("a", "z")).to eq("n")}
      it { expect(midstring("a", "c")).to eq("b")}
      it { expect(midstring("x", "z")).to eq("y")}
      it { expect(midstring("abc", "abchi")).to eq("abcd")}
      it { expect(midstring("abcde", "abchi")).to eq("abcf")}
    end

    context "consecutive characters" do
      it { expect(midstring("abh", "abit")).to eq("abhn")}
      it { expect(midstring("abhs", "abit")).to eq("abhw")}
    end

    context "a and b cases" do
      it { expect(midstring("a", "ab")).to eq("aan")}
      it { expect(midstring("abc", "abcb")).to eq("abcan")}
      it { expect(midstring("abc", "abcaah")).to eq("abcaad")}
    end
  end
end
