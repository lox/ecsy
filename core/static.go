package core

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDir struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDir) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDir{fs: _escLocal, name: name}
	}
	return _escDir{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/templates/build/ecs-service.json": {
		local:   "templates/build/ecs-service.json",
		size:    15948,
		modtime: 1448349201,
		compressed: `
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
AAAA//8tQQvPTD4AAA==
`,
	},

	"/templates/build/ecs-stack.json": {
		local:   "templates/build/ecs-stack.json",
		size:    11786,
		modtime: 1448349201,
		compressed: `
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
`,
	},

	"/templates/build/vpc.json": {
		local:   "templates/build/vpc.json",
		size:    3348,
		modtime: 1448349201,
		compressed: `
H4sIAAAJbogA/9SW32/TMBDH3/tXRH4OEBuERN5Cx6aCBFMbFWmIBze9DWupHcUOqEL537GTtHXcpCla
ux/aqtT2fe27z93F/TvyPBR9n8WwylKq4FLkK6rmkEsmOAo9RAIcvAo+6H/kG9tvhcoKJfWSkeqJeZZs
B2ZI0wKsCT01hVuzkzFsJsvqWfr1DrNiwcHasnuXSx6GnwUzTv3Yzup5H/nW0F7yLHnblfrEALWWS/8/
tPgBWuJordHPkfutQTVqztA7SVHkCfRnIF5nBp3Jahh+GpMwnF+Pt5DQdS4yyBUD2eY7Zsv8YyqSe6PF
wevq7w1+b+FFEy4V5QnEwPVjbSyXcEuLVNlWMb2TTpbaQNAXqLRf6QpQm1xH5tsIq6hmiib3lbqXZNlH
sim6K13sf+j6MLgJV5DrjG2MD0PsPCBS2tdfK+BqMEf7isGkOQ5Olg65DbVNANslq2BN/fQKrabdMe1q
3+BweLXRcEDRb8pSumApU+sbwd1KqN4CM0ghUU6FeV4w0INGegUqupEd5eUU2BTuzPvPMSqP6Fsb7H5L
Yd1S5B06NXz85PDxC4BPzgOfPDl88gLgvz01/KkoVIthB/vKJqaLFIb5n8Kdi+YyHHZq2J8LkIpxqnQ2
WjQ3N3Ngs3zQ+39HqVffwD7uJjgmM7Xp7uRISpGwKtphNLW419m9X3dnDBY/j2DxowRLnkew5FzBjsyn
HP0LAAD//1i9GBwUDQAA
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/templates": {
		isDir: true,
		local: "/templates",
	},

	"/templates/build": {
		isDir: true,
		local: "/templates/build",
	},
}
