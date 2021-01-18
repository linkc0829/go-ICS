#!/bin/sh
app="gql-server"
src="$srcPath/$app/$pkgFile"

printf "\nStart running: $app\n"
time /$GOPATH/bin/realize start run
printf "\nStopped running: $app\n\n"