arg APP_VERSION="development"
arg APP_BUILD_TS="unknown"
arg APP_SOURCE_VCS
arg APP_SOURCE_REPO

# build app frontend

from node:22-alpine3.22 as app-builder

arg APP_VERSION
arg APP_BUILD_TS
arg APP_SOURCE_VCS
arg APP_SOURCE_REPO

workdir /app

copy ./app .

run npm i

env VITE_APP_VERSION=${APP_VERSION}
env VITE_APP_BUILD_TS=${APP_BUILD_TS}
env VITE_APP_PLATFORM="docker-containerized"

env VITE_APP_SOURCE_VCS=${APP_SOURCE_VCS}
env VITE_APP_SOURCE_REPO=${APP_SOURCE_REPO}

run npm run build

# build backend service

from golang:1.26.3-alpine as svc-builder

arg APP_VERSION
arg APP_BUILD_TS

workdir /app

copy . .

copy --from=app-builder /app/dist /app/cmd/service/web/dist

run apk add --no-cache make build-base libwebp-dev

env CGO_ENABLED=1

run go build -v -ldflags "-s -w -X main.Version=${APP_VERSION} -X main.BuildTS=${APP_BUILD_TS}" -o service ./cmd/service

from alpine:3.23.4

workdir /app

copy --from=svc-builder /app/service /app

run apk add --no-cache ca-certificates libwebp

cmd  ["/app/service"]
