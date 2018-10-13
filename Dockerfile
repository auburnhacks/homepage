FROM golang:1.10-alpine as build

RUN mkdir -p /go/src \
    && mkdir -p /go/bin \
    && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

WORKDIR $GOPATH/src/homepage
RUN apk update && apk add git
RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.* ./
RUN dep ensure -vendor-only
#COPY . .
ADD ./metadata .
RUN CGO_ENABLED=0 go install -a std
RUN CGO_ENABLED=0 GOOS='linux' go build -a -ldflags '-extldflags "-static"' -installsuffix cgo -o homepage .

FROM alpine:latest
WORKDIR /app
# Run update on alpine and create ssl certificates
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*
# Copy executable from current directory
COPY --from=build /go/src/homepage/homepage .
# Copy static assets from current directory
COPY ./static/. ./static
EXPOSE 8321
ENTRYPOINT [ "./homepage" ]
