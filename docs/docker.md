## Docker

### 1. build docker

```sh
$ docker build -t go-git-webhook .
```

### 2. 配置

```sh
$ cp conf/app.conf.example conf/app.conf
```

### 3. RUN

*注意:* `app.conf` 数据库连接修改为:

```ini
#数据库配置
db_host=mysql
db_port=3306
db_database=webhook_db
db_username=root
db_password=
```

```sh
# 起一个 mysql
$ docker run --name mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true -v `pwd`/create-db.sql:/docker-entrypoint-initdb.d/create-db.sql  -d mysql
# init db
$ docker run --rm --link mysql:mysql  -v `pwd`/conf/app.conf:/go-git-webhook/conf/app.conf go-git-webhook go-git-webhook orm syncdb webhook
# create admin
$ docker run --rm --link mysql:mysql  -v `pwd`/conf/app.conf:/go-git-webhook/conf/app.conf go-git-webhook go-git-webhook install -account=admin -password=123456 -email=admin@163.com
# 起服务
$ docker run -d -p 8080:8080 --name go-git-webhook --link mysql:mysql  -v `pwd`/conf/app.conf:/go-git-webhook/conf/app.conf go-git-webhook go-git-webhook
```

## Docker-compose

### 1. build

```sh
$ docker-compose build
```

### 2. 配置

```sh
$ cp conf/app.conf.example conf/app.conf
```

### 3. RUN

*注意:* `app.conf` 数据库连接修改为:

```ini
#数据库配置
db_host=mysql
db_port=3306
db_database=webhook_db
db_username=root
db_password=
```

```sh
# init db
$ docker-compose run --rm go-git-webhook go-git-webhook orm syncdb webhook
# create admin
$ docker-compose run --rm go-git-webhook go-git-webhook install -account=admin -password=123456 -email=admin@163.com
# start
$ docker-compose up -d
```