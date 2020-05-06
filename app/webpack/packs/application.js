require("@rails/ujs").start();
require("turbolinks").start();
require("@rails/activestorage").start();
require("channels");
require("stylesheets/application");

// use with <%= image_pack_tag 'rails.png' %> or `imagePath`
const images = require.context("../images", true);
const imagePath = (name) => images(name, true);
