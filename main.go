package main

import (
	"exceapp/cmd/config"
	"exceapp/pkg/google"
	"fmt"
)

func main() {
	config.ConnectDB()

	url := google.GetLoginUrl("khan")
	fmt.Println(url)
	fmt.Println("docker runnning....🎉")
}
