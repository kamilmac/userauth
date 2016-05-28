package users

import (
    "time"
    "fmt"
    "github.com/dgrijalva/jwt-go"
)

var signingKey = ""

type user struct {
    username        string
    password        string
    token           string
    loggedIn        bool
}

type Users map[string]*user

func (u *user) grantToken() string {
    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims["username"] = u.username
    token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
    tokenString, _ := token.SignedString([]byte(signingKey))
    u.token = tokenString
    u.loggedIn = true
    return tokenString
}

func (u *Users) SetSigningKey(key string) {
   signingKey = key 
}

func (u *Users) Register(username, password string) {
    newUser := &user{
        username: username,
        password: password,
        token: "",
        loggedIn: false,
    }
    (*u)[username] = newUser
}

func (u *Users) Delete(username string) {
    delete(*u, username)
}

func (u *Users) Login(username, password string) (bool, string) {
    if _, present := (*u)[username]; !present {
        return false, ""
    }
    if password != (*u)[username].password {
        return false, ""
    }
    return true, (*u)[username].grantToken() 
}

func (u *Users) Auth(token string) (bool, string) {
    tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return []byte(signingKey), nil
    })
    if err != nil {
        return false, ""
    }
    tokenUser := tokenParsed.Claims["username"]
    username := fmt.Sprintf("%v", tokenUser)
    if _, ok := tokenParsed.Method.(*jwt.SigningMethodHMAC); !ok {
        (*u)[username].loggedIn = false
        return false, ""
    } 
    if !tokenParsed.Valid {
        (*u)[username].loggedIn = false
        return false, ""
    }
    tokenParsed.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
    return true, (*u)[username].username
}

func Init() *Users {
    return &Users{}
}