#!/usr/bin/env bash

oto -template templates/server.go.plush \
	-out server.gen.go \
	-pkg main \
	./def
gofmt -w server.gen.go server.gen.go
echo "generated server.gen.go"

oto -template templates/store.go.plush \
	-out store.gen.go \
	-pkg main \
	./def
gofmt -w store.gen.go store.gen.go
echo "generated store.gen.go"

oto -template templates/client.js.plush \
	-out www/src/client.gen.js \
	-pkg main \
	./def
echo "generated client.gen.js"
