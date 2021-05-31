const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
  mode: "development",
  entry: path.resolve(__dirname, "src", "index.jsx"),
  output: {
    path: path.resolve(__dirname, "./dist"),
    filename: "assets/bundle.js",
  },
  devServer: {
    historyApiFallback: true,
    contentBase: path.resolve(__dirname, "./dist"),
    open: true,
    compress: true,
    port: 8080,
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, "src", "index.html"),
    }),
  ],
  module: {
    rules: [
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ["babel-loader"],
      },
    ],
  },
};
