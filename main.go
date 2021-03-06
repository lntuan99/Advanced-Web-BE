package main

import (
    "advanced-web.hcmus/api/routers"
    "advanced-web.hcmus/config"
    "advanced-web.hcmus/model"
    "advanced-web.hcmus/services/smtp"
    "fmt"
    "os"
    "os/signal"
)

func main() {
    fmt.Println("=====================================================")
    fmt.Println("Start API Server......")
    fmt.Println("=====================================================")

    // Fetch configs
    config.FetchEnvironmentVariables()

    // Initialize configs
    model.Initialize()
    defer model.DBInstance.Close()

    // Initialize smtp service
    smtp.Initialize()

    go func() {
        r := routers.Initialize()
        port := os.Getenv("PORT")
        if port == "" {
            port = config.Config.ApiPort
        }
        r.Run(fmt.Sprintf(":%s", port))
    }()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    <-c

    fmt.Println("=====================================================")
    fmt.Println("API Server has stopped!")
}
