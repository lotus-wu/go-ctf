package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("sourse")))
	fmt.Println("files should put on 'sourse' dir")
	fmt.Println("listen on 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}
}
