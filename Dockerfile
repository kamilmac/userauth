FROM golang
COPY . /app/userauth
WORKDIR /app/userauth
RUN go get github.com/dgrijalva/jwt-go
RUN go build main.go
CMD ["./main"]
# EXPOSE 8080