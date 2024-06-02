FROM golang AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w -extldflags '-static'" -o sunoapi

FROM alpine:latest

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata gcc \
    && update-ca-certificates 2>/dev/null || true

COPY --from=builder /build/sunoapi /
COPY --from=builder /build/template /template/

EXPOSE 8000

ENTRYPOINT ["/sunoapi"]
