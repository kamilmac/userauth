FROM golang
COPY . /go/src/github.com/kamilmac/userauth
WORKDIR /go/src/github.com/kamilmac/userauth
RUN go get github.com/dgrijalva/jwt-go
RUN go build main.go
CMD ["./main"]
EXPOSE 5000