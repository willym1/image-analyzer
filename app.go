package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    
    "github.com/willym1/image-analyzer/test"
)

func main() {
    // initialize colors dictionary
    colors_file, err := ioutil.ReadFile("./colors.json")
    if err != nil {
        log.Fatal("colors.json not found")
    }
    var colors interface{}
    json.Unmarshal(colors_file, &colors)

    test.Run(true)
}
