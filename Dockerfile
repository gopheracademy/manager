FROM node:14.5.0-alpine3.12 AS front_builder

ADD ./www /www
WORKDIR /www
RUN npm install && npm run build

# Backend Build Step
FROM golang:1.15.4-alpine3.12 AS builder

# Prerequisites
RUN apk update && apk add --no-cache upx

# Dependencies
WORKDIR $GOPATH/src/github.com/gopheracademy/manager
COPY . .
RUN go mod download
RUN go mod verify

# Copy frontend build
COPY --from=front_builder /www/build $GOPATH/src/github.com/gopheracademy/manager/www/build/

# Build
RUN CGO_ENABLED=0 go build \
			-o /tmp/manager \
      github.com/gopheracademy/manager

# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/manager /go/bin/manager
VOLUME [ "/data" ]
COPY --from=front_builder /www/build /data/www/build/
WORKDIR /data
EXPOSE 8000
ENTRYPOINT ["/go/bin/manager"]
