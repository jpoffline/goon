package main

import (
	"fmt"
	"goon/goon"
)

func run() {

}

func main() {
	cfg := goon.NewConfig()
	cfg.Print()
	ef := goon.NewEvolvingField(cfg, "KINK")
	ef.Evolve()
	fmt.Println("DONE")
}
