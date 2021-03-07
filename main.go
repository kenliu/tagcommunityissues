package main

import (
	. "github.com/kenliu/TagCommunityIssues/TagCommunityIssues"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", TagCommunityIssues)
	_ = http.ListenAndServe(":8888", nil)
	log.Print("starting server")
}
