package main

import (
    "context"
    "fmt"
    "github.com/liamg/github-profile-magic-action/config"
    "github.com/liamg/github-profile-magic-action/profile"
    "os"
    "time"
)

const outputDir = "build"

func main() {

    conf, err := config.Derive()
    if err != nil {
        fail("Configuration error: %s", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    if err := profile.New(conf).Generate(ctx, outputDir); err != nil {
        fail("Failed to generate profile: %s", err)
    }

    // TODO: commit + push content to profile readme
}

func fail(format string, args ...interface{}) {
    _, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
    os.Exit(1)
}
