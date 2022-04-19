# 一个做到一半摆烂的项目

## 简介
这个项目里已经做好了用户注册登录、找回密码、Github登录、绑定等和用户相关的服务。可以以此为扩展做业务逻辑。

## 使用准备
- 腾讯短信包（免费有100条额度），并且申请了模板和签名
- 阿里云OSS对象存储
- Redis
- Mysql8
- Github App

## 项目里有什么
- gin
- gorm and mysql
- cors
- tencent cms
- oss
- mail
- viper
- log(distinct debug)
- docker
- snow Flake
- jwt auth
- random
- websocket

## How to use
- [ ] Globally replace the package name with your own repository
- [ ] Edit config/vars GlobalConfig. **It is recommended to make changes on the existing basis. Try not to change the existing structure, if you change, you need to change part of the code synchronously**
- [ ] Exec `go build cmd/main.go` and run, `config.json` will generate under `config/`
- [ ] Complete the config
- [ ] If you deploy with docker engine, edit `docker-copmose.yml`, Especially port mappings and service names

