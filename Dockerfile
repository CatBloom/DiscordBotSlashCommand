FROM golang:1.2 
WORKDIR /go/src/work
COPY go.mod ./
RUN go mod download
COPY . ./
RUN go build .