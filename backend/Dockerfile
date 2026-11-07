FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=build /out/server /app/server

ENV PORT=3000 \
    DATA_PATH=/data/kotoba.json \
    KANJI_DATA_PATH=/data/kanji.json \
    BUNPOU_DATA_PATH=/data/bunpou.json \
    LATIHAN_BUNPOU_DATA_PATH=/data/latihan/bunpou

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget -qO- http://localhost:3000/api/health || exit 1

ENTRYPOINT ["/app/server"]
