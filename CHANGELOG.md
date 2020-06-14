# CHANGELOG

### boards, lists, cards - May 16, 2020

- boards, lists, cards CRUD actions
- card dnd
- move images and css to `app/assets`
- Next steps: 
 - card detail popup
 - pundit authorization
 - user invites

### devise setup - May 7, 2020

- styled 'pages#home' and 'pages#pricing'
- generated devise user and styled devise views
- override `ActionView::Base.field_error_proc`
- add `name` field to devise user model

### react setup - May 6, 2020

- created Procfile for use with foreman gem
- created `react_component()` and `mount()` helpers

### initial setup - May 6, 2020

- ran `rvm use 2.7.1`
- ran `rails new agilely --database=postgresql --webpack=react -T --skip-sprockets`
- delete `/app/assets` and `<%= stylesheet_link_tag ...%>` and set `config.generators.assets = false`
- change webpack entry point to `app/webpack`
- move stylesheets and images to `/app/webpack`
- setup tailwindcss and purge in `/app/webpack/stylesheets/application.css`
- setup rspec and configure database_cleaner and factory_bot
