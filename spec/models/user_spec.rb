require 'rails_helper'

RSpec.describe User, type: :model do
  let(:user) { create(:user) }
  let!(:board) { create(:board) }
  let!(:participation) { board.participations.create(user_id: user.id, role: "admin") }

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
    it { expect(user).to validate_length_of(:password).is_at_least(6) }
    it { expect(user.admin).to be false }
    it { expect(user).to_not allow_value("").for(:password) }
    it { expect(user.dup).to_not allow_value(user.email.upcase).for(:email) }

    it "is expected to follow valid email regex" do
      addresses = %w[user@foo,com user_at_foo.org example.user@foo. foo@bar_baz.com foo@bar+baz.com foo@bar..com]
      addresses.each do |invalid_address|
        expect(user).to_not allow_value(invalid_address).for(:email)
      end
    end
  end

  describe "associations" do
    it { expect(user).to have_many(:participations) }
    it { expect(user).to have_many(:boards).through(:participations).source(:participant).dependent(:destroy) }
  end

  describe ".has_participation_in?" do
    it { expect(user.has_participation_in? board).to be true }
  end

  describe ".participation_in" do
    it { expect(user.participation_in board).to eq(participation) }
  end

  describe ".role_in" do
    it { expect(user.role_in board).to eq(participation.role) }
  end

  describe "#can_edit?" do
    it { expect(user.can_edit? board).to be true }
  end

  describe ".titles" do
    it { expect(user.boards_titles.first.as_json).to eq({"id"=> nil, "title"=> board.title, "slug"=> board.slug}) }
  end
end
