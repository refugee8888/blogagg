package main

import (
	"fmt"

	"github.com/refugee8888/blogagg/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := cfg.SetUser("Paul"); err != nil {
		fmt.Println(err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg)
}
