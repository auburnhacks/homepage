FROM golang:1.10 as builder

# Using dep v0.5.0
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/auburnhacks/homepage/
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix nocgo -o homepage .

# Release Image
FROM heroku/heroku:16
FROM alpine:latest
WORKDIR /app

# Run update on alpine and create ssl certificates
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

# Copy executable from builder step
COPY --from=builder go/src/github.com/auburnhacks/homepage/homepage ./
COPY ./static/. ./static
RUN chmod +x homepage
CMD [ "./homepage" , "--static_dir", "./static" ]
