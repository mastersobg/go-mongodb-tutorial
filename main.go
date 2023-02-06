package main

import (
	"context"

	"go-mongodb-tutorial/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
