FROM golang:1.22  as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -C ./src/content-service -o /app/publish/content-service

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/publish/* ./

CMD ["/app/content-service"]