package main

import "fmt"

type DampingConfig struct {
	mag float64
	cut int
}
type FieldConfig struct {
	ncomponents int
}
type GridConfig struct {
	nx int
	ny int
}
type Config struct {
	ht         float64
	hx         float64
	hy         float64
	ht2        float64
	hx2        float64
	hy2        float64
	ntimesteps int
	damping    DampingConfig
	field      FieldConfig
	grid       GridConfig
	outloc     string
	outfreq    int
}

func NewConfig() Config {
	c := Config{}
	c.grid.nx = 105
	c.grid.ny = 105
	c.ht = 0.01
	c.hx = 0.1
	c.hy = c.hx
	c.hx2 = c.hx * c.hx
	c.hy2 = c.hx * c.hx
	c.ht2 = c.ht * c.ht
	c.ntimesteps = 100
	c.damping.cut = 100
	c.damping.mag = 5.0
	c.field.ncomponents = 1
	c.outloc = "out"
	c.outfreq = 20
	return c
}

func (c *Config) Print() {
	fmt.Println("CONFIG")
	fmt.Println("(Nx, Ny) = ", c.grid.nx, c.grid.ny)
	fmt.Println("(hx, ht) = ", c.hx, c.ht)
	fmt.Println("outloc = ", c.outloc)
}
