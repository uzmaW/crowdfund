export default {
  devServer: {
    proxy: {
      '^/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
        pathRewrite: { '^/api': '' },
      },
    },
  },
  configureWebpack: {
    devServer: {
      static: './dist', // Serve from the 'dist' directory
    },
  },
};