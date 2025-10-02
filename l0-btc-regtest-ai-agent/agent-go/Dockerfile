FROM golang:1.22 as build
WORKDIR /src
COPY . .
RUN go build -o /out/agent ./cmd/agent

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/agent /agent
COPY .env /env/.env
ENV GOTRACEBACK=single
EXPOSE 9090 7070
ENTRYPOINT ["/agent"]
