FROM golang:1.13 AS build

COPY go.mod go.sum ./
RUN go get -d -v ./...
COPY . ./
RUN go build -o /go/bin/action-labels cmd/action-labels/main.go

FROM gcr.io/distroless/base
COPY --from=build /go/bin/action-labels /
CMD ["/action-labels"]
