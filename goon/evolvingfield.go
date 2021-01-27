package goon

import "goon/goon/field"

type EvolvingField struct {
	timestep int
	steps    []field.Field
	config   Config
	idx_now  int
	idx_pre  int
}

// NewEvolvingField creates required data structure with provided initial config
func NewEvolvingField(cfg Config, ict string) EvolvingField {
	ef := EvolvingField{config: cfg}
	ef.timestep = 100
	if ict == "KINK" {
		f1 := field.NewField(cfg.grid.nx, cfg.grid.ny, cfg.field.ncomponents)
		f2 := field.NewField(cfg.grid.nx, cfg.grid.ny, cfg.field.ncomponents)
		f1.Kink(cfg.grid.nx, cfg.grid.ny, cfg.field.ncomponents)
		f2.Kink(cfg.grid.nx, cfg.grid.ny, cfg.field.ncomponents)
		ef.steps = append(ef.steps, f1)
		ef.steps = append(ef.steps, f2)
		return ef
	}
	return EvolvingField{}

}

func (ef *EvolvingField) incrementSteps(timestep int) (int, int) {
	ef.timestep = timestep
	if ef.idx_now == 0 {
		ef.idx_now = 1
		ef.idx_pre = 0
	} else {
		ef.idx_now = 0
		ef.idx_pre = 1
	}
	return ef.idx_now, ef.idx_pre
}

func (ef *EvolvingField) ToFile() {
	nc := ef.config.field.ncomponents
	files := NewFieldFiles(nc, ef.config.outloc, ef.timestep)

	t := ef.idx_now
	nx, ny := ef.config.grid.nx, ef.config.grid.ny
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			for c := 0; c < nc; c++ {
				fm := ef.steps[t].Value(i, j, c)
				files.Write(c, FloatToString(fm)+" ")
			}
		}
		files.NewLine()
	}
	files.Close()

}

func dampingFactors(alpha, ht float64) (float64, float64) {
	return 1 / (1 + alpha*ht/2), (1 - alpha*ht/2)
}

func (ef *EvolvingField) Evolve() {
	var i, j, ip1, im1, jp1, jm1 int
	alpha := ef.config.damping.mag
	ht, ht2 := ef.config.ht, ef.config.ht2
	hx2, _ := ef.config.hx2, ef.config.hy2

	nsteps := ef.config.ntimesteps
	for timestep := 0; timestep < nsteps; timestep++ {
		if timestep > ef.config.damping.cut {
			alpha = 0.0
		}
		fac1, fac2 := dampingFactors(alpha, ht)
		t, tp := ef.incrementSteps(timestep)

		for i = 1; i < ef.config.grid.nx-1; i++ {
			for j = 1; j < ef.config.grid.ny-1; j++ {
				ip1 = i + 1
				im1 = i - 1
				jp1 = j + 1
				jm1 = j - 1
				ef.steps[t].Update(i, ip1, im1, j, jp1, jm1, hx2, ht2, fac1, fac2, &ef.steps[tp])
			}
		}

		i = 0
		ip1 = 1
		im1 = ef.config.grid.nx - 1
		for j = 1; j < ef.config.grid.ny-1; j++ {
			jp1 = j + 1
			jm1 = j - 1
			ef.steps[t].Update(i, ip1, im1, j, jp1, jm1, hx2, ht2, fac1, fac2, &ef.steps[tp])
		}

		i = ef.config.grid.nx - 1
		ip1 = 0
		im1 = ef.config.grid.nx - 2
		for j = 1; j < ef.config.grid.ny-1; j++ {
			jp1 = j + 1
			jm1 = j - 1
			ef.steps[t].Update(i, ip1, im1, j, jp1, jm1, hx2, ht2, fac1, fac2, &ef.steps[tp])
		}

		j = ef.config.grid.ny - 1
		jp1 = 0
		jm1 = ef.config.grid.ny - 2
		for i = 1; i < ef.config.grid.nx-1; i++ {
			ip1 = i + 1
			im1 = i - 1
			ef.steps[t].Update(i, ip1, im1, j, jp1, jm1, hx2, ht2, fac1, fac2, &ef.steps[tp])
		}

		j = 0
		jp1 = 1
		jm1 = ef.config.grid.ny - 1
		for i = 1; i < ef.config.grid.nx-1; i++ {
			ip1 = i + 1
			im1 = i - 1
			ef.steps[t].Update(i, ip1, im1, j, jp1, jm1, hx2, ht2, fac1, fac2, &ef.steps[tp])
		}

		if timestep%ef.config.outfreq == 0 {
			ef.ToFile()
		}

	}
}
