### userauth
------
## Routes:
### api/v1/login 
    Input: {"username":"admin","password":"admin"}
    Output: {"status":"SUCCESS/ERROR","message":"MSG","data":{"token":"NEWTOKEN"}}

### api/v1/auth
    Input: {"token":"TOKEN"}
    Output: {"status":"SUCCESS/ERROR","message":"MSG","data":{"token":"NEWTOKEN"}}


## Flags || Env_vars || defaults
    -port, -adminUser, -adminPass, -signingKey
    USERAUTH_PORT|_ADMINUSER|_ADMINPASS|_SIGNINGKEY
    Defaults: 5000, admin, admin, secretkey 

Run "go test" while running main package to test api calls