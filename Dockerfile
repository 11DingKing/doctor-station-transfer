FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/doctor ./cmd/server
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/doctor /app/doctor
EXPOSE 8080
ENTRYPOINT ["/app/doctor"]
