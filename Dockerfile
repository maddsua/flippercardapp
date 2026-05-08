from node:22 as web-builder

workdir /app/web

copy ./web .

run npm i
run npm run build

from golang:1.24.4-alpine as svc-builder

workdir /app

copy . .

copy --from=web-builder /app/web/dist /app/cmd/web/dist

run go build -v -ldflags "-s -w" -o service ./cmd

from alpine:3.23.4

workdir /app

copy --from=svc-builder /app/service /app

run apk add ca-certificates

cmd  ["/app/service"]
