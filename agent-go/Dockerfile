FROM golang:1.22 AS build
WORKDIR /src
ENV CGO_ENABLED=1 GOPROXY=https://proxy.golang.org,https://goproxy.io,direct GOSUMDB=off
RUN apt-get update && apt-get install -y --no-install-recommends libzmq3-dev libsodium-dev ca-certificates && rm -rf /var/lib/apt/lists/*
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o /out/agent ./cmd/agent

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends libzmq5 libsodium23 ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/agent /agent
ENV GOTRACEBACK=single
EXPOSE 9090
ENTRYPOINT ["/agent"]
