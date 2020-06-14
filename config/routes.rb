Rails.application.routes.draw do
  root to: 'pages#home'
  get '/pricing', to: 'pages#pricing'
  get '/dashboard', to: 'pages#dashboard'
  devise_for :users, path: '', path_names: { sign_in: 'login', sign_out: 'logout', sign_up: 'signup' }
  get '/b/:slug', to: 'boards#show', as: 'show_board'
  resources :boards, param: :slug, only: [:create, :destroy, :update], shallow: true do
    resources :lists, only: [:create, :destroy, :update], shallow: true
  end
  resources :lists, only: [], shallow: true do
    resources :cards, only: [:create, :destroy], shallow: true
  end
  patch '/cards/:id/reorder', to: 'cards#reorder'
  patch '/lists/:id/reorder', to: 'lists#reorder'
  constraints subdomain: 'api' do
    resources :boards
  end
end
