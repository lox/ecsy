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

	"/templates/src/ecs-service.yml": {
		local:   "templates/src/ecs-service.yml",
		size:    5154,
		modtime: 1480301724,
		compressed: `
H4sIAAAJbogA/+xXTW/bOBO+61dMjAIGijpf7ftiVzfHcbYGvFkjctJD0QNN0xERmdSSVII0yH/foURZ
pCUnTrB76xeQcGaeh8N5OBwNBoNo+C2Zs3WeEcMupFoTc8OU5lLE0D89PjkeHP+O//rROdNU8dxUlvEo
gYSpe05ZDMP6x09AYE70HZyzFRfc+gIRS/wP4+lZP4pmRJE1M0gQRwA3s9FkaX8AmD/mCPQtiePx6DSO
0RLHk2VpCojnKQO+ZMLwFWcK5MqCgJGgCgFcRBXqrFhknCbFQjATMgw3FJV1N0teYoAu3RAaDC5aMp0z
atmX8MBNWiVR8yp+j8dYQZ+8l3nFlTaQV1hbG9iX/PS95JpRiRV7GzuqYZQVGuvqkyZGcXHbTWP1Q6sQ
Wz1iDKFpxVFJCVcd8kRoQwRlCaOF4ubxDyWLfEdygUtnjkObYekEt9YLOYkBigollDKty51xR6ntFqyg
L8iaZ4/7JrcqvUGg1K1CbVbG3go8xUKz/xKzuXf74oaXybjShFi2Qo5kJIUhXDB1iRvZl4PWQX6ppV/t
AHomlfGhL4v1gqkdtxR9QVbqDGhkzoQLWZEiMzH8dlwKanr2XvxMkiUsSGaV8SLHV0Yyk45SRu+uVbbv
KV1fTS1oyg08pEzAUqIrpCUWtVgaVkquqxpNz0LiI8ubJNMRU7aaFO9u2AF2MQ/FlgTwJiAO0AbIlR9W
UnWT93pRhNVblmLRFeu1Zl+Nyaccr7io+wLAwfjvgmQavsPBFVu1dvwJseBH5CPoFsSlNDZ+HySL9Vdh
8sKU20pQenfleZRYNyQr8Gh6jOoBJrdGjtj+7ETZcw1oinU/c2UPAg8mK/judgXQ295w79PGdoBtF3qp
tcZHRx+evs7nswD38PwysVfqOf7w5CT63Bm/CX812gX/cFnUr3WQgD25xtbVmXzPxhZFV0zLQmGLjKtL
9dbuXHrMFF4hLBrTdXGhNHoKBZRu2LCtDJHQdj68D8JY0WIONcBNTidLjHqC/oXAN2Btr3CZRT+uzrH/
4emSmQep7ko9PA/KkD481xDBRifiVuG7sNkgDGCS48aNpDKLwdB8YwG4wAtathd3tFUtPIe5fNE84ks1
yWM4Piz/Hh2/kfXLl88dZOFqm8O2rC1FtSuYEW04bXywm8SxH1JGbNpAvN0BdhS8mkeC4y3PpjXFRR21
aYdtK7Fl3zFPOL96swGwn+SL1auRPafgUfNcm1Lak3cG791o+OdE3TLjpNuz3njJA9hnbAjBi7O5+zXk
4zxFDacyW8ZwurFdi7RlPWkUNxE4nN0T3OLnZnHO10wWuJ3/uSXciWDUVvxc4Y5QFTOJVXtsEhgLssgY
YhtVsDbQ/zf662i0/4YA9S8F7q3AbrKwf7ybI/EsrVGl+xX/dTG2X+7gUU3so1o93e64t7Rdf5VtpOcW
nBnfWa7YciQLgYQntQQ9AQQyDIf/rdLbNf+VCab5l2XiEzpoO1l1TVQVUKtbeOv+Mg4/FcOVzNrjjl2M
ai0F307NuNMshsUoAbcLMhn+Gccb1HY1hloX6zK0ksK5pPi7MM0Z40BiWLhkD368WqGUkCXL5INnsSQc
O2FOsjhYBghE4/8ZAI64h2RNfkpBHvQhlevAZ0ibz0c/SmOLbBJw5hkxqf34cL/ZrPxhDsOqTKuqerO1
3yR2nMXO83jtTHZn4fKvXhL7ObfYvCTVyLlgH98SotitVaaqW6K2k1jrSdoL7Gob6hs36etQ9PS1naPH
sDCpVPwn6xpsW1H1ZI8fRh970T8BAAD//yAivo4iFAAA
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    9603,
		modtime: 1480302185,
		compressed: `
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
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    3175,
		modtime: 1480302232,
		compressed: `
H4sIAAAJbogA/7xWX2/aPhR951Pcop+Ul9JfklWTljcGtEPbWgSokzrtwQS3tRpsK3baoYrvvhvH5B8k
pa36BynC9/j43ON7b+j1ep3+r9mcrmREND0T8YroKxorJngAju96bs/9gh+nM6QqjJnUWWQ0mMHVZAAX
VD+K+D4AwSnIZBGxEFSy4FQD4UvAIMiYPSC3XVbHQMJYKGVi/Wt14nQ6P4mUjN+qoAMwM7CB4DfsNv0O
6TkBPMFgPJwG4Lkn5v9/7zNsTDjb4NYgPkL80wrEq0E+7UL8GuR0C+lcJlom2ihEPeOllUaihAZwNKU3
6bJZG/2VItZZHOCCrFIA0kP3vyf0OghmmoT36fqmZ6i6nTxtd2IszDZXHbcAazK6Vj+/wtCuxNlVUtnt
FIq8SXZ9zZK87Q03atpyvFLUdntJlf+sKv9ZVf7bVPmFqilVIolDuq2OjGa+lkhiNo8GfhBs62MSC0lj
zTJ4+jdgy/hrJEJso6Mzxpdjjg0BvyutcAxdJOjiIy3PLvyxe8dcacJDOqccH+sAlvSGJJG24REni4gO
uZol0mQIOk5oPfhNKM0xN1UJz8ltLhGgB98p0qcW5GtVW7tVo0xZn6NDj2S935Ex1zTGDC2o7g48bUoU
fY3EdyvKdaO9O0hbHZLypbrE2mg4KLeyIgd7PMurvAszlmEeSa8Ul0fjSVmTU4jCiNNw2lCsCENNDzJM
OS6I3uNViaoANDH2I6whknaBEWirFyWclCMbx8Kz6spzqY6P/ROp7nkGaZCDNZztHctL/oMkPLyrlFf/
gbCILFjE9PoaXx+pZBrRUGPlu8dwdE41vh7AcfJaP7BPrPDdXtlzdfuG3IuyfOcsvBdk4b9PFt7bs/AP
ymIqEk1VW70ZxDwdWg15NLO2WfMa2mE2Zm1LGoYW8qYZQJVm3DRmyc/tLxzXotrnUSE+B5SNrM6WV+ps
PqP0Cj00nUJNTlYslQZPdkJbMRey+kqJkJmTG1JonXYHO2mb8gOkeVVrn72BYgp8gDj/heL+BQAA///6
Ev11ZwwAAA==
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

	"/templates/src": {
		isDir: true,
		local: "/templates/src",
	},
}
