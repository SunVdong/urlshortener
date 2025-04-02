package main

import "github.com/sunvdong/urlshortener/application"

func main() {
	a := application.Application{}

	if err := a.InitApp("./config/config.yaml"); err != nil {
		panic(err)
	}

	a.Run()
}
