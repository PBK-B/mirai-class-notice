# 课表通知 QQ 机器人使用教程

1, 创建文件夹

`mkdir class_notice && cd class_notice`

2, 下载安装包

`wget https://github.com/PBK-B/mirai-class-notice/releases/download/v0.1-Beta/linux_v0.1_beta.tar.gz`

3, 解压缩

`tar -zxvf linux_v0.1_beta.tar.gz`

4, 编辑配置文件

`vim conf/app.conf`

5, 修改配置文件

```
appname = class_notice
httpaddr = 0.0.0.0
httpport = 8089
runmode = prod
dbhost = 127.0.0.1:3306
dbdriver = mysql
dbusername = root
dbpassword = mysqlpassword
dbdatabase = test_class_notice
```

6, 运行 class_notice 程序

`nohup ./class_notice &`

7, 访问 http://ip:8089/admin/

8, 登陆后台，默认用户名 `admin`，默认密码 `admin`

9, 设置上课时间

![docs/images/img-01.png](/docs/images/img-01.png)

10, 添加课程

![docs/images/img-02.png](/docs/images/img-02.png)

11, 设置开学时间，上课周数以及登陆 QQ 账号和选择需要通知的 QQ 群

![docs/images/img-03.png](/docs/images/img-03.png)