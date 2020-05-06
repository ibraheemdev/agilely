# CHANGELOG

### initial setup - May 6, 2020

- ran `rvm use 2.7.1`
- ran `rails new agilely --database=postgresql --webpack=react -T --skip-sprockets`
- delete `/app/assets` and `<%= stylesheet_link_tag ...%>` and set `config.generators.assets = false`
- change webpack entry point to `app/webpack`
- move stylesheets and images to `/app/webpack`
- setup tailwindcss and purge in `/app/webpack/stylesheets/application.css`
- setup rspec and configure database_cleaner and factory_bot
