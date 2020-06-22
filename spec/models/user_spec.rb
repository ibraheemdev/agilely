require 'rails_helper'

RSpec.describe User, type: :model do
  let(:user) { create(:user) }

  describe "attributes" do
    it { expect(user).to respond_to(:email) }
    it { expect(user).to respond_to(:encrypted_password) }
    it { expect(user).to respond_to(:reset_password_token) }
    it { expect(user).to respond_to(:reset_password_sent_at) }
    it { expect(user).to respond_to(:remember_created_at) }
    it { expect(user).to respond_to(:confirmation_token) }
    it { expect(user).to respond_to(:confirmed_at) }
    it { expect(user).to respond_to(:confirmation_sent_at) }
    it { expect(user).to respond_to(:unconfirmed_email) }
    it { expect(user).to respond_to(:name) } 
    it { expect(user).to respond_to(:created_at) }
    it { expect(user).to respond_to(:updated_at) }
    it { expect(user).to respond_to(:admin) }
  end

  describe 'validations' do
    it { expect(user).to validate_presence_of(:name) }
    it { expect(user).to validate_length_of(:name).is_at_most(50) }
    it { expect(user).to validate_presence_of(:email) }
    it { expect(user).to validate_presence_of(:password) }
    it { expect(user).to validate_presence_of(:password_confirmation).on(:create) }
    it { expect(user.admin).to be false }

    it "is expected to validate that :password confirmation should match :password" do
      user.password_confirmation = "notthepassword"
      expect(user.valid?).to be false
    end

    it "is expected to validate that :email is unique" do
      user_with_same_email = user.dup
      user_with_same_email.email = user.email.upcase
      expect(user_with_same_email.valid?).to be false
    end

    it "is expected to follow valid email regex" do
      addresses = %w[user@foo,com user_at_foo.org example.user@foo. foo@bar_baz.com foo@bar+baz.com foo@bar..com]
      addresses.each do |invalid_address|
        user.email = invalid_address
        expect(user.valid?).to be false
      end
    end
  end

  describe "associations" do
    it { expect(user).to have_many(:participations) }
    it { expect(user).to have_many(:boards).through(:participations).source(:participant).dependent(:destroy) }
  end
end
