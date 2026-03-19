FROM cgr.dev/chainguard/wolfi-base:latest AS builder

RUN apk add --no-cache go-1.26=1.26.1-r1

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd.version=${VERSION} -X github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd.commit=${COMMIT} -X github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd.date=${DATE}" \
    -o /apk-datasource ./cmd/apk-datasource

FROM cgr.dev/chainguard/wolfi-base:latest

COPY --from=builder /apk-datasource /usr/local/bin/apk-datasource

USER nonroot:nonroot

ENTRYPOINT ["apk-datasource"]
