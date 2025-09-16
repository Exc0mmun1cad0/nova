FROM golang:1.24-alpine3.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags '-static'" -o nova main.go


FROM gcr.io/distroless/static-debian12:latest-amd64 AS release-stage

WORKDIR /

COPY --from=build-stage /app/nova /nova

USER nonroot:nonroot

ENTRYPOINT [ "./nova" ]