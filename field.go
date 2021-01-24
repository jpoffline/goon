package main

import (
	"os"
	"strconv"
)

type FieldValue float64

type Field struct {
	ncomponents int
	field       [][][]float64
	dpot        []float64
}
type EvolvingField struct {
	timestep int
	steps    []Field
	config   Config
	idx_now  int
	idx_pre  int
}

func (f *Field) Value(i, j, c int) float64 {
	return f.field[i][j][c]
}

func (f *Field) Set(i, j, c int, v float64) {
	f.field[i][j][c] = v
}

func (f *Field) Modulus(nc int, i, j int) float64 {
	mod := 0.0
	for c := 0; c < nc; c++ {
		s := f.Value(i, j, c)
		mod += s * s
	}
	return mod
}

func (ef *EvolvingField) IncrementSteps(timestep int) (int, int) {
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

func newField(cfg Config) Field {
	nx := cfg.grid.nx
	ny := cfg.grid.ny
	nc := cfg.field.ncomponents
	f := make([][][]float64, ny)
	for c := 0; c < ny; c++ {
		f[c] = make([][]float64, nx)
	}
	for x := 0; x < nx; x++ {
		for y := 0; y < ny; y++ {
			f[x][y] = make([]float64, nc)
		}
	}
	return Field{field: f, ncomponents: cfg.field.ncomponents, dpot: make([]float64, cfg.field.ncomponents)}
}

func NewEvolvingField(cfg Config, ict string) EvolvingField {
	ef := EvolvingField{config: cfg}
	ef.timestep = 100
	if ict == "KINK" {
		f1 := newField(cfg)
		f2 := newField(cfg)
		f1.Kink(cfg)
		f2.Kink(cfg)
		ef.steps = append(ef.steps, f1)
		ef.steps = append(ef.steps, f2)
		return ef
	}
	return EvolvingField{}

}

type FieldFiles struct {
	fs []*os.File
	nf int
}

func (ff *FieldFiles) NewLine() {
	for c := 0; c < ff.nf; c++ {
		ff.fs[c].WriteString("\n")
	}
}
func (ff *FieldFiles) Close() {
	for c := 0; c < ff.nf; c++ {
		ff.fs[c].Close()
	}
}
func NewFieldFiles(nc int, loc string, ts int) FieldFiles {
	ff := FieldFiles{}
	ff.nf = nc
	for c := 0; c < nc; c++ {
		fl, _ := os.Create(outputFileName(loc, "field_"+strconv.Itoa(c), ts))
		ff.fs = append(ff.fs, fl)
	}
	return ff
}

func (ff *FieldFiles) Write(id int, v string) {
	ff.fs[id].WriteString(v)
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

func outputFileName(loc, id string, timestep int) string {
	loc = loc + "/"
	ts := strconv.Itoa(timestep)
	fn := loc + id + "_" + Pad(ts, 5) + ".dat"
	return fn
}
func (f *Field) LaplacianForComponent(c, i, ip1, im1, j, jp1, jm1 int, h2 float64) float64 {
	ff := 2 * f.Value(i, j, c)
	ddx := f.Value(ip1, j, c) + f.Value(im1, j, c) - ff
	ddy := f.Value(i, jp1, c) + f.Value(i, jm1, c) - ff
	return (ddx + ddy) / h2
}

func (f *Field) CalcDpot(i, j int) {
	mf := 0.0
	for c := 0; c < f.ncomponents; c++ {
		mf += f.Value(i, j, c)
	}

	for c := 0; c < f.ncomponents; c++ {
		f.dpot[c] = f.Value(i, j, c) * (mf - 1)
	}

}
func (ef *EvolvingField) Evolve() {
	alpha := ef.config.damping.mag
	ht, ht2 := ef.config.ht, ef.config.ht2
	hx2, _ := ef.config.hx2, ef.config.hy2
	nc := ef.config.field.ncomponents
	for timestep := 0; timestep < ef.config.ntimesteps; timestep++ {
		if timestep > ef.config.damping.cut {
			alpha = 0.0
		}
		fac1, fac2 := 1/(1+alpha*ht/2), (1 - alpha*ht/2)

		t, tp := ef.IncrementSteps(timestep)
		for i := 1; i < ef.config.grid.nx-1; i++ {
			for j := 1; j < ef.config.grid.ny-1; j++ {
				ip1 := i + 1
				im1 := i - 1
				jp1 := j + 1
				jm1 := j - 1
				for c := 0; c < nc; c++ {
					lap := ef.steps[t].LaplacianForComponent(c, i, ip1, im1, j, jp1, jm1, hx2)
					eom := lap - ef.steps[t].dpot[c]
					nv := 2.0*ef.steps[t].Value(i, j, c) + eom*ht2 - fac2*ef.steps[tp].Value(i, j, c)
					ef.steps[tp].Set(i, j, c, nv*fac1)
				}
			}
		}
		if timestep%ef.config.outfreq == 0 {
			ef.ToFile()
		}
	}
}
