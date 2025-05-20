package main

import (
	"exceapp/pkg/google"
	"fmt"
)

func main() {
	url := google.GetLoginUrl("khan")
	fmt.Println(url)

}
