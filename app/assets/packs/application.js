require("@rails/ujs").start();
require("turbolinks").start();
require("@rails/activestorage").start();
require("javascript/channels");
require("stylesheets/application");
const images = require.context("images", true);