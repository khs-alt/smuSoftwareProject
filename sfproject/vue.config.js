const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: "/",
  devServer: {
    // port: 8080,
    // disableHostCheck: true
    proxy: {
      "/api": {
        target: "http://localhost:8000",
        //     // changeOrigin: true,
        ws: false,
      },
    },
  },
});