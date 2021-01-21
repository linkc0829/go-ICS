#!/bin/sh
set -ex
mongo mongo:27017/users --eval "var user = '$MONGO_INITDB_ROOT_USERNAME', pwd = '$MONGO_INITDB_ROOT_PASSWORD'" scripts/createAdmin.js
echo "server admin created"