#!/usr/bin/env bash

oto -template server.go.plush \
	-out server.gen.go \
	-pkg main \
	./def
gofmt -w server.gen.go server.gen.go
echo "generated server.gen.go"

oto -template data.go.plush \
	-out data.gen.go \
	-pkg main \
	./def
gofmt -w data.gen.go data.gen.go
echo "generated data.gen.go"

oto -template client.js.plush \
	-out www/src/client.gen.js \
	-pkg main \
	./def
echo "generated client.gen.js"
