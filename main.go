package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "log"
    "flag"
    "os"
    
    "github.com/kamilmac/userauth/users"
)

var (
    port            string
    signingKey      string
    adminUser       string
    adminPass       string
)

type App struct {
    userbase        *users.Users
}

type Response struct {
    Status          string      `json:"status"`
    Message         string      `json:"message"`
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

type AuthRes struct {
    Response
    Username        string      `json:"username"`
}

func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
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
        if ok, newToken := app.userbase.Login(req.Username, req.Password); !ok {
            res.Status = "error"
            res.Message = "Login failed"
        } else {
            res.Status = "success"
            res.Data.Token = newToken
        }
    }
    json, _ := json.Marshal(res)
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

func (app *App) HandleAuth(w http.ResponseWriter, r *http.Request) {
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
        if ok, username := app.userbase.Auth(req.Token); !ok {
            res.Status = "error"
            res.Message = "Login failed"
        } else {
            res.Status = "success"
            res.Username = username
        }
    }
    json, _ := json.Marshal(res)
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

func firstPresent(a, b string) string {
    if a == "" { return b }
    return a
}

func init() {
    flag.StringVar(
        &port, 
        "port", 
        firstPresent(
            os.Getenv("CFG_USERAUTH_PORT"), "5000",
        ), 
        "Specify port",
    )
    flag.StringVar(
        &adminUser, 
        "adminUser", 
        firstPresent(
            os.Getenv("CFG_USERAUTH_ADMINUSER"), "admin",
        ), 
        "Specify admin username",
    )
    flag.StringVar(
        &adminPass, 
        "adminPass", 
        firstPresent(
            os.Getenv("CFG_USERAUTH_ADMINPASS"), "admin",
        ), 
        "Specify admin password",
    )
    flag.StringVar(
        &signingKey, 
        "signingKey", 
        firstPresent(
            os.Getenv("CFG_USERAUTH_SIGNINGKEY"), "secretkey",
        ), 
        "Specify user token signing key",
    )
    flag.Parse()
}

func main() {
    app := App{}
    app.userbase = users.Init()
    app.userbase.SetSigningKey(signingKey)
    app.userbase.Register(adminUser, adminPass)
    // http.HandleFunc("/register", handleRegister)
    http.HandleFunc("/api/v1/login", app.HandleLogin)
    http.HandleFunc("/api/v1/auth", app.HandleAuth)
    log.Println("Running on port:", port)
    err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}