package main

import (
    "fmt"
    "log"
    "net/http"
)

type Config struct {}
const port = "80"

func main () {
    app := Config{}

    log.Println("Starting mail service on port:", port)

    srv := &http.Server{
        Addr: fmt.Sprintf(":%s", port),
        Handler: app.routes(),
    }

    err := srv.ListenAndServe()
    if err!= nil {
        log.Panic(err)
    }

}