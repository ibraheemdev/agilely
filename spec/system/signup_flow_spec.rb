require 'rails_helper'
require 'support/signup_helper'
include SignupHelper

RSpec.describe "User signs up", type: :system do
  context 'with valid credentials' do
    let(:credentials) { { email: 'valid@example.com', name: 'john', password: 'password', password_confirmation: 'password' } }
    
    it "displays flash notice" do
      sign_up_with credentials
      expect(page).to have_content('A message with a confirmation link has been sent to your email address.')
    end

    it "creates the user" do
      expect { sign_up_with credentials }.to change{ User.count }.by(1)
    end

    it "sends a confirmation email" do
      expect { sign_up_with credentials }.to change{ ActionMailer::Base.deliveries.count }.by(1)
      expect(@email).to have_content("You can confirm your account email through the link below:")
    end

    it "confirms the user" do
      sign_up_with credentials
      confirm_account
      expect(User.find_by(email: 'valid@example.com').confirmed?).to be true
    end

    it "redirects to login" do
      sign_up_with credentials
      confirm_account
      expect(page).to have_current_path(new_user_session_path)
    end
  end

  context 'with invalid credentials' do
    let(:credentials) { { email: 'invalid$email', name: 'john', password: 'ps', password_confirmation: 'p' } }
    it "displays error message" do
      sign_up_with credentials
      expect(page).to have_content("We ran into a couple errors")
    end
    it "doesn't create the user" do
      expect { sign_up_with credentials }.to change{ User.count }.by(0)
    end
    it "doesn't send a confirmation email" do
      expect { sign_up_with credentials }.to change{ ActionMailer::Base.deliveries.count }.by(0)
    end
  end
end