FROM golang:1.23.4-alpine AS builder

RUN apk update && apk add --no-cache git make

WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest

# install basic dependency
RUN apk --update upgrade && \
    apk add curl ca-certificates

# update certificates
RUN update-ca-certificates

# cleanup apk cache
RUN rm -rf /var/cache/apk/*

COPY --from=builder /app/bin/rest-api /bin/rest-api

EXPOSE 8080
CMD ["/bin/rest-api"]
