Rails.application.routes.draw do
  root to: 'pages#home'
  get '/pricing', to: 'pages#pricing'
  devise_for :users, path: '', path_names: { sign_in: 'login', sign_out: 'logout', sign_up: 'signup' }
end
