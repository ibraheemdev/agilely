module SignupHelper
  def sign_up_with(credentials)
    visit new_user_registration_path
    fill_in 'Email', with: credentials[:email]
    fill_in 'Name', with: credentials[:name]
    fill_in 'Password', with: credentials[:password]
    click_button 'Sign up'
    @email = ActionMailer::Base.deliveries.last
  end
  
  def confirm_account
    visit @email.body.raw_source.match(/(?:"https?\:\/\/.*?)(\/.*?)(?:")/)[1]
  end
end