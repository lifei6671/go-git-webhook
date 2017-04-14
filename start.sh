#!/bin/sh
set -e

cd /usr/src/github.com/lifei6671/go-git-webhook/

goFile="go-git-webhook"

if [ ! -f go-git-webhook] ; then
    go build -ldflags "-w"
fi

chmod +x $goFile

if [ ! -f "conf/app.conf" ] ; then
    cp conf/app.conf.example conf/app.conf
fi

if [ ! -z $db_host ] ; then
    sed -i 's/^db_host.*/db_host='$db_host'/g' conf/app.conf
fi

if [ ! -z $db_port ] ; then
    sed -i 's/^db_port.*/db_port='$db_port'/g' conf/app.conf
fi

if [ ! -z $db_database ] ; then
    sed -i 's/^db_database.*/db_database='$db_database'/g' conf/app.conf
fi

if [ ! -z $db_username ] ; then
    sed -i 's/^db_username.*/db_username='$db_username'/g' conf/app.conf
fi

if [ ! -z $db_password ] ; then
    sed -i 's/^db_password.*/db_password='$db_password'/g' conf/app.conf
fi

if [ ! -z $httpport ] ; then
    sed -i 's/^httpport.*/httpport='$httpport'/g' conf/app.conf
fi

#初始化数据
if [ ! -z $install ] ; then
    ./$goFile orm syncdb
    ./$goFile install -account=$account -password=$password -email=$email
fi

exec ./$goFile