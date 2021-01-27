package field

//Kink sets the initial configuration of the field
// to be a crude kink.
func (f *Field) Kink(nx, ny, nc int) {
	nxf := float64(nx)
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			for c := 0; c < nc; c++ {
				if i < int(nxf*0.5) {
					f.Set(i, j, c, -1.0)
				} else {
					f.Set(i, j, c, 1.0)
				}

			}
		}
	}
}
