FROM  golang:1.12.1  as builder
WORKDIR /go/src/github.com/Lasiar/au-back/
COPY ./ /go/src/github.com/Lasiar/au-back/
RUN go get ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM amd64/alpine
RUN mkdir /app
COPY --from=builder /go/src/github.com/Lasiar/au-back/main /app
COPY config.json /app/config.json
#COPY Tree.xml /app/Tree.xml
WORKDIR /app
EXPOSE 80
CMD ["/app/main"]
