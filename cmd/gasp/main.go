package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/anacrolix/gasp"
	"github.com/anacrolix/tagflag"
)

var args struct {
	ProjectDir string `type:"pos"`
}

type bootstrap struct {
	refs map[string]map[string]struct{}
}

func (me *bootstrap) addImport(pkg string) {
	if me.refs == nil {
		me.refs = make(map[string]map[string]struct{})
	}
	pkgRefs := me.refs[pkg]
	if pkgRefs == nil {
		pkgRefs = make(map[string]struct{})
		me.refs[pkg] = pkgRefs
	}
}

func (me *bootstrap) addRef(pkg, name string) {
	me.addImport(pkg)
	me.refs[pkg][name] = struct{}{}
}

func (me *bootstrap) walkObject(obj gasp.Object) error {
	// log.Printf("%T", obj)
	switch v := obj.(type) {
	case gasp.Symbol:
		// log.Printf("%q", v.Name())
		if !strings.HasPrefix(v.Name(), "go:") {
			break
		}
		parts := strings.Split(v.Name()[3:], ".")
		me.addRef(parts[0], parts[1])
	case gasp.List:
		for !v.Empty() {
			me.walkObject(v.First())
			v = v.Rest()
		}
	}
	return nil
}

func parseImport(obj gasp.Object) (full, short string, ok bool, err error) {
	l, ok := obj.(gasp.List)
	if !ok {
		return
	}
	if l.Empty() {
		return
	}
	i, ok := l.First().(gasp.Symbol)
	if !ok {
		return
	}
	if i.Name() != "import" {
		ok = false
		return
	}
	l = l.Rest()
	if l.Empty() {
		err = errors.New("bad import")
		return
	}
	s, ok := l.First().(gasp.String)
	if !ok {
		err = errors.New("bad import")
		return
	}
	full = s.Value
	index := strings.LastIndex(full, "/")
	if index == -1 {
		short = full
	} else {
		short = full[index+1:]
	}
	ok = true
	return
}

func (me *bootstrap) addRefs(r io.Reader) error {
	gr := gasp.NewReader(r)
	for {
		o, err := gr.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		shortMap := make(map[string]string)
		full, short, ok, err := parseImport(o)
		if err != nil {
			return err
		}
		if ok {
			shortMap[short] = full
		} else {
			me.walkObject(o)
		}
	}
}

func writeBootstrap(_w io.Writer) error {
	d, err := os.Open(args.ProjectDir)
	if err != nil {
		return err
	}
	defer d.Close()
	fis, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	var b bootstrap
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".gasp") {
			continue
		}
		f, err := os.Open(filepath.Join(args.ProjectDir, fi.Name()))
		if err != nil {
			return err
		}
		err = b.addRefs(f)
		f.Close()
		if err != nil {
			return err
		}
	}
	b.addImport("github.com/anacrolix/gasp")
	b.addImport("os")
	w := io.MultiWriter( /*os.Stderr,*/ _w)
	return template.Must(template.New("main").Parse(`
package main

import ({{ range $key, $value := . }}
    "{{ $key }}"
{{ end }})

func main() {
	env := gasp.Env{
		NS: gasp.NewMap(),
	}
	{{ range $pkg, $names := . }}
		{{ range $name, $value := $names }}
			env.NS = env.NS.Assoc(gasp.NewSymbol("go:{{$pkg}}.{{$name}}"), gasp.WrapGo({{$pkg}}.{{$name}}))
		{{ end }}
	{{ end }}
	err := env.RunProject(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running project: %s", err)
		os.Exit(1)
	}
}
`)).Execute(w, b.refs)
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	tagflag.Parse(&args)
	tempName := filepath.Join(os.TempDir(), fmt.Sprintf("gasp.%d.go", os.Getpid()))
	f, err := os.OpenFile(tempName, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	err = writeBootstrap(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %q", f.Name())
	goLook, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	err = syscall.Exec(goLook, []string{"go", "run", f.Name(), args.ProjectDir}, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error execing go run: %s\n", err)
		os.Exit(1)
	}
}
