package main

import (
    "fmt"
    "time"
    "net/http"
    "encoding/json"
    "log"
    "errors"
    "flag"
    
    "github.com/dgrijalva/jwt-go"
)

type Response struct {
    Status          string      `json:"status"`
    Message         string      `json:"message"`
}

type User struct {
    ID              string
    Username        string
    Email           string
    Password        string
    Token           string
}

type LoginReq struct {
    Username        string      `json:"username"`      
    Password        string      `json:"password"`
}

type LoginData struct {
    Token           string      `json:"token"`
}

type LoginRes struct {
    Response
    Data            LoginData   `json:"data"`
}

type AuthReq struct {
    Token           string      `json:"token"`      
}

type AuthData struct {
    UserID          string      `json:"userid"`
}

type AuthRes struct {
    Response
    Data            AuthData    `json:"data"`
}

const (
    signingKey = "thisisnotrandomtext"
)

var (
    port int
    me = User{"149sdfinw9fas0315jfs9", "kamil", "kamil@gmail.com", "limak", ""}
    Users = map[string]User {
        me.ID: me,
    }
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
    var (
        res LoginRes
        req LoginReq
    ) 
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&req)
    if err != nil {
        res.Status = "error"
        res.Message = "Json req decoding error"
    } else {
        if user, err := getUserByName(req.Username); err != nil { 
            res.Status = "error"
            res.Message = fmt.Sprintf("%v", err)
        } else if user.Password != req.Password {
            res.Status = "error"
            res.Message = "Wrong password"
        } else {
            user.Token = createToken(user)
            res.Status = "success"
            res.Data.Token = user.Token
        }
    }
    json, _ := json.Marshal(res)
    w.Write(json)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
    var (
        res AuthRes
        req AuthReq
    ) 
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&req)
    if err != nil {
        res.Status = "error"
        res.Message = "Json req decoding error"
    } else {
        token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
            return []byte(signingKey), nil
        })
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            res.Status = "error"
            res.Message = "Unexpected signing method"
        } else if err != nil || !token.Valid {
            res.Status = "error"
            res.Message = "Token invalid or expired"
        } else {
            userID, _ := token.Claims["userid"]
            res.Status = "success"
            res.Data.UserID = fmt.Sprintf("%v", userID)
        }
    }
    json, _ := json.Marshal(res)
    w.Write(json)
}

func getUserByName(username string) (*User, error) {
    for _, v := range Users {
        if v.Username == username {
            return &v, nil
        }
    }
    return nil, errors.New("User doesn't exist")
}

func createToken(user *User) string {
    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims["userid"] = user.ID
    token.Claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
    tokenString, _ := token.SignedString([]byte(signingKey))
    return tokenString
}

func init() {
    flag.IntVar(&port, "port", 5000, "Specify port")
    flag.Parse()
    fmt.Println("Running on port:", port)
}

func main() {
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/auth", handleAuth)
    err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}