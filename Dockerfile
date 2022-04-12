FROM golang:1.17 as build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o app .

FROM debian:buster-slim
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt
COPY --from=build /src/app /app
ENTRYPOINT ["/app"]
