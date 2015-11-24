package gasp

import (
	"io"
	"os"
	"path/filepath"
)

type Env struct {
	NS Map
}

func (me *Env) EvalFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	r := NewReader(f)
	for {
		obj, err := r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		Eval(obj, *me)
	}
}

func (env *Env) RunProject(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	fis, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		err := env.EvalFile(filepath.Join(dir, fi.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
