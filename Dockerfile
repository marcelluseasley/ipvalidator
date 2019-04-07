FROM golang:1.11.7-stretch

LABEL maintainer="Marcellus Easley <marcellus.easley@gmail.com>"

WORKDIR $GOPATH/src/github.com/marcelluseasley/ipvalidator

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["ipvalidator"]