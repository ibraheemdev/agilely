const path = require('path')

module.exports = {
  resolve: {
    alias: {
      '@redux': path.resolve(__dirname, '..', '..', 'app/views/boards/show/redux'),
    },
  },
};
