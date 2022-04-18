FROM golang:1.18 as api-plant-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM api-plant-base as api-plant-builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/api-plant ./cmd/api-plant/.

FROM alpine:latest
COPY --from=api-plant-builder /app/build/api-plant /usr/bin
EXPOSE 8080
ENTRYPOINT ["api-plant"]
