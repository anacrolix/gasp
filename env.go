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
	fmt.Fprintf(w, "}\n")
	if env.Outer != nil {
		env.Outer.WriteNS(w)
	}
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
(def or (fn [& next] (if (empty? next) false (if (first next) (first next) (apply or (rest next))))))
(def == (fn [a b] (not (or (< a b) (> b a)))))
(def reduce (fn (f cum vals)
	(if (empty? vals)
		cum
		(reduce f (f cum (first vals)) (rest vals)))))
(def flip (fn (f) (fn (a b) (f b a))))
(def reverse (fn (coll) (reduce (flip cons) () coll)))
(def list (fn (& elems) (reverse (reduce (fn (coll e) (cons e coll)) () elems))))
(def len (fn [x] (if (empty? x) 0 (+ 1 (len (rest x))))))
(def nth (fn [l i] (if (== i 0) (first l) (nth (rest l) (- i 1)))))
(def apply (fn [f args] (eval (cons f args))))
(def defn (macro (fn (name params & body) (list 'def name (concat (list 'fn params) body)))))
(defn conj (a & b) (concat a b))
(defn and [& a]
	(if (empty? a) true
		(if (first a)
			(apply and (rest a))
			false)))
`)
	for _, o := range objs {
		Eval(o, ret)
	}
	return
}
