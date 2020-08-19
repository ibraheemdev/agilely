#!/usr/bin/env node

process.env.NODE_ENV = "development";

const fs = require("fs-extra");
const paths = require("react-scripts/config/paths");
const webpack = require("webpack");
const config = require("react-scripts/config/webpack.config")("development");

// update the webpack dev config in order to remove the use of webpack devserver
config.entry = config.entry.filter(
  (fileName) => !fileName.match(/webpackHotDevClient/)
);
config.plugins = config.plugins.filter(
  (plugin) => !(plugin instanceof webpack.HotModuleReplacementPlugin)
);

// update the paths in config
config.output.path = paths.appBuild;
config.output.publicPath = "";

fs.emptyDir(paths.appBuild)
  .then(() => {
    return new Promise((resolve, reject) => {
      const webpackCompiler = webpack(config);
      new webpack.ProgressPlugin().apply(webpackCompiler);

      webpackCompiler.watch({}, (err, stats) => {
        if (err) {
          return reject(err);
        }
        return resolve();
      });
    });
  })
  .then(() =>
    fs.copy(paths.appPublic, paths.appBuild, {
      dereference: true,
      filter: (file) => file !== paths.appHtml,
    })
  );
