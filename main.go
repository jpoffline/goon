package main

func run() {

}

func main() {
	cfg := NewConfig()
	cfg.Print()
	ef := NewEvolvingField(cfg, "KINK")
	ef.Evolve()
}
