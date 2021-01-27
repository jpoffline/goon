package field

type FieldValue float64

type Field struct {
	ncomponents int
	field       [][][]float64
	dpot        []float64
}

func NewField(nx, ny, nc int) Field {

	f := make([][][]float64, ny)
	for c := 0; c < ny; c++ {
		f[c] = make([][]float64, nx)
	}
	for x := 0; x < nx; x++ {
		for y := 0; y < ny; y++ {
			f[x][y] = make([]float64, nc)
		}
	}
	return Field{field: f, ncomponents: nc, dpot: make([]float64, nc)}
}

func (f *Field) Value(i, j, c int) float64 {
	return f.field[i][j][c]
}

func (f *Field) DpotValue(c int) float64 {
	return f.dpot[c]
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

func (f *Field) LaplacianForComponent(c, i, ip1, im1, j, jp1, jm1 int, h2 float64) float64 {
	ff := 2 * f.Value(i, j, c)
	ddx := f.Value(ip1, j, c) + f.Value(im1, j, c) - ff
	ddy := f.Value(i, jp1, c) + f.Value(i, jm1, c) - ff
	return (ddx + ddy) / h2
}

func (f *Field) CalcDpot(i, j int) {
	mf := 0.0
	for c := 0; c < f.ncomponents; c++ {
		mf += f.Value(i, j, c) * f.Value(i, j, c)
	}

	for c := 0; c < f.ncomponents; c++ {
		f.dpot[c] = f.Value(i, j, c) * (mf - 1.0)
	}

}

func (f *Field) Update(i, ip1, im1, j, jp1, jm1, nc int, hx2, ht2, fac1, fac2 float64, fp1 *Field) {
	for c := 0; c < nc; c++ {
		f.CalcDpot(i, j)
		lap := f.LaplacianForComponent(c, i, ip1, im1, j, jp1, jm1, hx2)
		eom := lap - f.DpotValue(c)
		nv := 2.0*f.Value(i, j, c) + eom*ht2 - fac2*fp1.Value(i, j, c)
		fp1.Set(i, j, c, nv*fac1)
	}
}
