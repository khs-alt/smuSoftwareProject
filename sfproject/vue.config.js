const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    // port: 8080,
    // disableHostCheck: true
    proxy: {
      "/": {
        target,
        changeOrigin: true,
      },
    },
  },
});
