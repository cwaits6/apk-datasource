FROM cgr.dev/chainguard/wolfi-base:latest AS builder

RUN apk add --no-cache go-1.26

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X github.com/cwaits6/apk-datasource/cmd/apk-index/cmd.version=${VERSION} -X github.com/cwaits6/apk-datasource/cmd/apk-index/cmd.commit=${COMMIT} -X github.com/cwaits6/apk-datasource/cmd/apk-index/cmd.date=${DATE}" \
    -o /apk-index ./cmd/apk-index

FROM cgr.dev/chainguard/wolfi-base:latest

COPY --from=builder /apk-index /usr/local/bin/apk-index

USER nonroot:nonroot

ENTRYPOINT ["apk-index"]
