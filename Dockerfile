FROM golang:1.14 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o finnhub_exporter .

FROM scratch
WORKDIR /app
COPY --from=builder /app/finnhub_exporter .
ENTRYPOINT ["./finnhub_exporter"]
