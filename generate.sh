#!/usr/bin/env bash
mkdir -p backend/generated
oto -template templates/server.go.plush \
	-out backend/generated/conference.gen.go \
	-pkg generated \
	./def
gofmt -w backend/generated/conference.gen.go backend/generated/conference.gen.go
echo "generated conference.gen.go"

oto -template templates/client.js.plush \
	-out admin/src/components/client/conference.gen.js \
	-pkg services \
	./def
echo "generated conference.gen.js"

