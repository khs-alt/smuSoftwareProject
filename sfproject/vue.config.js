const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: "/",
  devServer: {
    // port: 8080,
    // disableHostCheck: true
    proxy: {
      "/": {
        target: "http://localhost:8080",
        //     // changeOrigin: true,
        ws: false,
      },
    },
  },
});
