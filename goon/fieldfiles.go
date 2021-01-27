package goon

import (
	"os"
	"strconv"
)

type FieldFiles struct {
	fs []*os.File
	nf int
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

func (ff *FieldFiles) Write(id int, v string) {
	ff.fs[id].WriteString(v)
}

func outputFileName(loc, id string, timestep int) string {
	loc = loc + "/"
	ts := strconv.Itoa(timestep)
	fn := loc + id + "_" + Pad(ts, 5) + ".dat"
	return fn
}
