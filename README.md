# go-git-webhook

[![Build Status](https://travis-ci.org/lifei6671/go-git-webhook.svg?branch=master)](https://travis-ci.org/lifei6671/go-git-webhook)
[![Build status](https://ci.appveyor.com/api/projects/status/tpm2k23umrqri2dd/branch/master?svg=true)](https://ci.appveyor.com/project/lifei6671/go-git-webhook/branch/master)

一个基于 Golang 开发的用于迅速搭建并使用 WebHook 进行自动化部署和运维系统，支持：Github / GitLab / GitOsc。

界面和开发思路参考于 [git-webhook](https://github.com/NetEaseGame/git-webhook) 。

在原作者的基础上解耦了WebHook和Server之间关系，实现了多对多关系。

因与服务器通信使用的是SSH方式，请注意保管服务器账号和密码。

# 如何使用？

**1、拉取源码**

```
git clone github.com/lifei6671/go-git-webhook.git

```

**2、编译源码**

```
#更新依赖
go get -d ./...

#编译项目
go build -v -tags "pam" -ldflags "-w"
```

**3、运行**

```
chmod 0777 go-git-webhook

#恢复数据库，请提前创建一个空的数据库
./go-git-webhook orm syncdb webhook

#创建管理员账户
./go-git-webhook install -account=admin -password=123456 -email=admin@163.com

```

# 后台运行

**使用nohup后台运行**

```bash
nohup ./go-git-webhook &
```


**使用supervisor运行**

```bash
[program:go-git-webhook]
command=/var/www/go-git-webhook/go-git-webhook > /dev/null > 2>&1
autostart=true
autorestart=true
startsecs=10

```


# 使用技术

go-git-webhook 基于beego框架1.7.2版本开发。编译于golang 1.8版本。使用glide作为包管理工具。

# 界面预览

![WebHook](https://github.com/lifei6671/go-git-webhook/blob/master/static/uploads/1.png?raw=true)

![New WebHook](https://github.com/lifei6671/go-git-webhook/blob/master/static/uploads/2.png?raw=true)

![WebHook And Server List](https://github.com/lifei6671/go-git-webhook/blob/master/static/uploads/4.png?raw=true)

![New Server](https://github.com/lifei6671/go-git-webhook/blob/master/static/uploads/7.png?raw=true)

![Scheduler List](https://github.com/lifei6671/go-git-webhook/blob/master/static/uploads/6.png?raw=true)

# 问题反馈

如发现 BUG 请在 [issues](https://github.com/lifei6671/go-git-webhook/issues) 中反馈。
