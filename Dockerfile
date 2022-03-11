FROM golang:1.17

WORKDIR /src
COPY . .

RUN go build -o /bin/action

ENTRYPOINT ["/bin/action"]
