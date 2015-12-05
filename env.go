package gasp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Env struct {
	NS    Map
	Outer *Env
}

func (me *Env) Define(name, val Object) {
	me.NS = me.NS.Assoc(name, val)
}

func (me *Env) Lookup(name Symbol) (ret Object) {
	ret = me.NS.Get(name)
	if ret != nil {
		return
	}
	if me.Outer == nil {
		return
	}
	return me.Outer.Lookup(name)
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
		Eval(obj, me)
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

func (env *Env) String() string {
	var buf bytes.Buffer
	env.WriteNS(&buf)
	return buf.String()
}

func (env *Env) WriteNS(w io.Writer) {
	fmt.Fprintf(w, "#{ ")
	for _, key := range env.NS.Keys() {
		fmt.Fprintf(w, "%s ", key.String())
	}
	fmt.Fprintf(w, "}")
}

func NewStandardEnv() (ret *Env) {
	ret = &Env{
		NS: NewMap(),
	}
	for _, b := range builtins {
		ret.NS = ret.NS.Assoc(b.Symbol, b.Object)
	}
	objs := ReadString(`
(def not (fn (x) (if x false true)))
(def > (fn (a b) (< b a)))
(def <= (fn (a b) (> b a)))
`)
	for _, o := range objs {
		Eval(o, ret)
	}
	return
}
