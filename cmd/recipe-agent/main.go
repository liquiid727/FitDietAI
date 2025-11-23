package main

import (
    "log"
    "cook/internal/recipe/server"
)

func main() {
    if err := server.Execute(); err != nil {
        log.Fatal(err)
    }
}