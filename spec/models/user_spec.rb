require 'rails_helper'

RSpec.describe User, type: :model do
  subject { create(:user) }
  let!(:board) { create(:board) }
  let!(:participation) { subject.participations.create(participant: board, role: :admin) }

  describe "attributes" do
    it { is_expected.to be_mongoid_document }
    it { is_expected.to have_timestamps }

    it { is_expected.to have_field(:email).of_type(String).with_default_value_of("") }
    it { is_expected.to validate_uniqueness_of(:email) }
    it "is expected to follow valid email regex" do
      addresses = %w[user@foo,com user_at_foo.org example.user@foo. foo@bar_baz.com foo@bar+baz.com foo@bar..com]
      addresses.each do |invalid_address|
        is_expected.to validate_format_of(:email).not_to_allow(invalid_address)
      end
    end

    it { is_expected.to have_field(:encrypted_password).of_type(String).with_default_value_of("") }
    it { is_expected.to validate_length_of(:password).with_minimum(6) }
    it { is_expected.to validate_presence_of(:password) }

    it { is_expected.to have_field(:reset_password_token).of_type(String) }
    it { is_expected.to have_field(:reset_password_sent_at).of_type(Time) }

    it { is_expected.to have_field(:remember_created_at).of_type(Time) }

    it { is_expected.to have_field(:confirmation_token).of_type(String) }
    it { is_expected.to have_field(:confirmed_at).of_type(Time) }
    it { is_expected.to have_field(:confirmation_sent_at).of_type(Time) }
    it { is_expected.to have_field(:unconfirmed_email).of_type(String) }

    it { is_expected.to have_field(:name).of_type(String) }
    it { is_expected.to validate_presence_of(:name) }
    it { is_expected.to validate_length_of(:name).with_maximum(50) }

    it { is_expected.to have_field(:admin).of_type(Mongoid::Boolean).with_default_value_of(false) }

    it { is_expected.to embed_many(:participations) }

    it do 
      is_expected
      .to have_index_for(confirmation_token: 1)
      .with_options(unique: true, name: "index_users_on_confirmation_token" )
    end

    it do 
      is_expected
      .to have_index_for(reset_password_token: 1)
      .with_options(unique: true, name: "index_users_on_reset_password_token" )
    end

    it do 
      is_expected
      .to have_index_for(email: 1)
      .with_options(unique: true, name: "index_users_on_email" )
    end
  end

  describe "#boards" do
    let(:board2) { create(:board) }
    it { expect(subject.boards).to include(board) }
    it { expect(subject.boards).not_to include(board2) }
  end

  describe ".#board_titles" do
    it { expect(subject.boards_titles).to eq([[board.title, board.slug]]) }
  end

  describe "#has_participation_in?" do
    it { expect(subject.has_participation_in? board).to be true }
  end

  describe "#participation_in" do
    it { expect(subject.participation_in board).to eq(participation) }
  end

  describe "#role_in" do
    it { expect(subject.role_in board).to eq(participation.role) }
  end

  describe "#can_edit?" do
    it { expect(subject.can_edit? board).to be true }
  end

  describe "#can_edit?" do
    it { expect(subject.can_edit? board).to be true }
  end
end
