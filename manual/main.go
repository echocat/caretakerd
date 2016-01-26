package main

import (
    "github.com/russross/blackfriday"
    "io/ioutil"
)

func main() {
    bytes, err := ioutil.ReadFile("manual/docs/configuration/examples.md")
    if err != nil {
        panic(err)
    }
    content := blackfriday.MarkdownCommon(bytes)
    err = ioutil.WriteFile("target/test.html", content, 0)
    if err != nil {
        panic(err)
    }

}
