FROM golang:1.13 AS build

WORKDIR /go/src/app
COPY . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app cmd/action-labels/main.go

FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
CMD ["/app"]
