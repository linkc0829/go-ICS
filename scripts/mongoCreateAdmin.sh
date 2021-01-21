#!/bin/sh
#https://blog.ruanbekker.com/blog/2019/05/04/using-mongodb-inside-drone-ci-services-for-unit-testing/
set -ex
mongo --eval "var user = '$MONGO_INITDB_ROOT_USERNAME', pwd = '$MONGO_INITDB_ROOT_PASSWORD'" scripts/createAdmin.js
echo "server admin created"