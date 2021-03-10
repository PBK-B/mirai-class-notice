let mix = require('laravel-mix');

// 后台登陆页面
mix
  .js("src/pages/login/index.js", "static/js/login.js")
  .react()
  .sass("src/pages/login/scss/index.scss", "static/css/login.css");

  