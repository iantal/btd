FROM golang:alpine as builder

ENV GO111MODULE="" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=github.com/iantal

WORKDIR /build

COPY . .

RUN apk add git

ARG GT

RUN echo ${GT}
RUN git config --global url."https://golang:${GT}@github.com".insteadOf "https://github.com"

RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .
RUN cp /build/config.yml .

FROM golang:alpine as deploy
COPY --from=builder /dist .
RUN apk update && apk add unzip && apk add bash && apk add git

ENV BASE_PATH="/opt/data"
VOLUME [ "/opt/data" ]
EXPOSE 8004
CMD ["./main"]