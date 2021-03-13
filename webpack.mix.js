let mix = require("laravel-mix");

// 后台登陆页面
mix
  .js("src/pages/login/index.js", "static/js/login.js")
  .react()
  .sass("src/pages/login/scss/index.scss", "static/css/login.css");

// 后台管理页面
mix
  .js("src/pages/admin/index.js", "static/js/admin.js")
  .react()
  .sass("src/pages/admin/scss/index.scss", "static/css/admin.css");
