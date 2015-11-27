package templates

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

	"/cloudformation/templates/build/ecs-service.json": {
		local:   "cloudformation/templates/build/ecs-service.json",
		size:    5171,
		modtime: 1448595006,
		compressed: `
H4sIAAAJbogA/6xXT2/bNhS/+1MIOhZt4mTYsOnm2EnrIcuM2EsPRQ40TUdEJFIgqRbpkO++R1qSSYqU
vCBui9rk4+/9/z3y30mSpLOv6w0pqwIpcsNFidQDEZJylmZJejm9mH6a/gF/049adoUEKokCAdjVp2Ft
g+TzDSpp8dKt6dWXimiEtRKUPZnTZn1BJBa0Uo2CTU6SvTmcMEBO+D5RsKQAM6EsqSVJzcnXj0dlC7Kn
jDYIb1FId4QpuqdEtPqu5+vEhU4U72sHsXlRS/D/XTU3mEGVc84UooyIOwjPG7XiFkNrQEohnOtvWr0k
4jvFMZ0rLlRI511dbiEGQzorOJroMPr6eUWYfXKP6kIr+X3qhvr26p20Fxztki0qEMP/z4L5ek1wLah6
+Sx4XYVMgd7Jsuv5ZZY5olm23EXtm0HQD7LJkxZO9lwcSuH2Cv5HKsmRTBDGRMo2TXaVLJlU2hfpJu2h
wsMWPqzmg3bpuIFMZ07rUXLw3lG2rreMKBlSOOdliRakoCVVZHdLpYpqbFC0k1gQ4J8uDJQ1+iaNzvTv
WlW1pRHqHuHnRmlnxAMqamMFwfITOFISkWX6e7DOIaq3UBxXTW2EgNoFWLphWfYnp9ryb90qrHfu6Y+9
A3u5UlV2fu6IJBaoBf2ZqJlSHniz7xv6sS+yuFsbhvB2Hp3fr+7BNBsz7J7sdTTbZnTBJiE17bfXE4gz
EORW4VE+iBeeOXE8S97F6+rrnkheC91VXYWB128lgGPJrwTQjVCUSNcwI+b3Q58XDs2Q7AVQqB4ZEBgr
Z7rpoaMzJ3Gtz5oQumUr86lj6ZI9CSAar+7cQkiXFbihOOaFBla48uomvRG89Oi6l9JgDXk4G/4eKHO6
E0udrnR6Zv6cT+0Dr+PFGqMFN+0FkorioyzM4SwLNmq0DPpEanvb7o6ncSSBx+i5ZTqJBDFy3J+I4aBa
Nmr6J+xwXRwwz47ZexRAOyJHsNx7ziCiXf9fNpvVuOtfCCpUPs8Jfvazu0HiifQNi00Ys9fj/MCY0Ib5
lK4/vv8nxsD4Fhg158ND5jFWVE1IXjY5ME7OC01cl47APyzvi1xMHZklg6nwHelM/OLubGhJeK3j+muw
YcBVRrAm24UAl6FbV7yg+MVPzzVD24LsDNUJGCgRJb9Nj1piPLJurh3xybHWk+MgNEoW/RFqlqMz03Uf
hg0VZDfnNdPmX9htavXfSKvG3iO+Na7cMF1HnhtRyPF2tR0aMdRn++gdJ9zm97zw0R26NKk1QsGkRN+z
NpAndHLZeba5pbec/ZVlRmK07mZS1qVBOzTMgmP4zfxc6Su5Is3Gt+GL5fV+D61ojCkK/qNHWmALhYFa
mT7v81d6bKw+Dybm+n+GSvSTM/RDnmFe9pnt0VvpMV06w01OAlQrlcyOYRm5dU9CO3YVrJDKdSzsh0Jq
Yn3IxUA/HjLS1HhqP3X8CTaUOiMRT19fb3NkLI1GaCCOZp8c7lH6nb7t7lGHq/GWfAhixk8J8qTvHKJ7
Jeur6cjjKY5376N9pSo/CQ1fnuACCM1qlXNBf5LgrTxw7jEY4/b5olPxITDGTy7Q6OV4ov+9Tv4LAAD/
/4LZvwIzFAAA
`,
	},

	"/cloudformation/templates/build/ecs-stack.json": {
		local:   "cloudformation/templates/build/ecs-stack.json",
		size:    12360,
		modtime: 1448595006,
		compressed: `
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
`,
	},

	"/cloudformation/templates/build/vpc.json": {
		local:   "cloudformation/templates/build/vpc.json",
		size:    3348,
		modtime: 1448510745,
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

	"/cloudformation": {
		isDir: true,
		local: "/cloudformation",
	},

	"/cloudformation/templates": {
		isDir: true,
		local: "/cloudformation/templates",
	},

	"/cloudformation/templates/build": {
		isDir: true,
		local: "/cloudformation/templates/build",
	},
}
