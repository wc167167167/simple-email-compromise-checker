FROM golang:alpine3.13

WORKDIR /src

COPY ./src .
RUN ls
RUN go build -o /app/ .

ENTRYPOINT [ "/app/testGo" ]
