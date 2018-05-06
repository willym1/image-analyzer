package main

import (
    "github.com/willym1/image-analyzer/env"
    "github.com/willym1/image-analyzer/analyzer"
)

func main() {
    env.Init()
    
    analyzer.Serve()
}
