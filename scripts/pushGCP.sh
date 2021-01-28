#!/bin/bash
#use this script after gccloud auth login
printf "starting...\n"

docker build -f ./docker/prod.dockerfile -t icsharing .

docker tag icsharing asia.gcr.io/red-atlas-303101/icsharing:https

docker push asia.gcr.io/red-atlas-303101/icsharing:https

printf "finish pushing image to GCP container registry"