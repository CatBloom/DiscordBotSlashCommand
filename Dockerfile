FROM golang:1.20 AS prod

WORKDIR /go/src/work
COPY go.mod go.sum ./
RUN go mod tidy
COPY . ./

RUN go fmt && go build -o app .

CMD ["./app"]

FROM golang:1.20 AS dev

WORKDIR /go/src/work
COPY go.mod go.sum ./
RUN go mod tidy
COPY . ./

RUN go install github.com/cosmtrek/air@v1.40.4 
RUN go fmt && go build -o app .

CMD ["air", "-c", ".air.toml"]