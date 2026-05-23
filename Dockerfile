from node:22-alpine3.22 as app-builder

workdir /app

copy ./app .

run npm i
run npm run build

from golang:1.26.3-alpine as svc-builder

workdir /app

copy . .

copy --from=app-builder /app/dist /app/cmd/web/dist

run apk add --no-cache make build-base libwebp-dev

arg CGO_ENABLED=1

run go build -v -ldflags "-s -w" -o service ./cmd

from alpine:3.23.4

workdir /app

copy --from=svc-builder /app/service /app

run apk add --no-cache ca-certificates libwebp

cmd  ["/app/service"]
