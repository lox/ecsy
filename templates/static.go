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

	"/templates/build/ecs-service.json": {
		local:   "templates/build/ecs-service.json",
		size:    5251,
		modtime: 1449012824,
		compressed: `
H4sIAAAJbogA/6xXTW/bOBO++1cIOhZt4uTFu9jVzbGT1ots1oi96aHIgabpiIhECiTVIl3kv++QlmSS
IiU3iNuiNjl85uuZGfLfSZKks6/rDSmrAilyw0WJ1AMRknKWZkl6Ob2Yfpr+AX/Tj1p2hQQqiQIB2NWn
YW2D5PMNKmnx0q3p1ZeKaIS1EpQ9mdNmfUEkFrRSjYJNTpK9OZwwQE74PlGwpAAzoSypJUnNydePR2UL
sqeMNghvUUh3hCm6p0S0+q7n68SFThTvawexeVFL8P9dNTeYQZVzzhSijIg7CM8bteIWQ2tASiGc629a
vSTiO8UxnSsuVEjnXV1uIQZDOis4mugw+vp5RZh9co/qQiv5feqG+vbqnbQXHO2SLSoQw79mwXy9JrgW
VL18FryuQqZA7WTZ9fwyyxzRLFvuovbNIOgH2eRJCyd7Lg5UuL2C/5FKciQThDGRsk2TzZIlk0r7It2k
PVR42MKH1XzQLh03kOnMaT1KDt47ytb1lhElQwrnvCzRghS0pIrsbqlUUY0NinYSCwL9pwsDZY2+SaMz
/btWVW1pBN4j/Nwo7Yx4QEVtrCBYfgJHSiKyTH8P8hyiegvkuGq4EQJqF2DphmXZn5xqy791q7Deuac/
9g7s5UpV2fm5I5JYoBb0Z6JmSnngzb5v6Me+yOJubTqEt/Po/H51D6bZmGH3ZK+j2RajCzYJqWm/vZ7Q
OANBbhUe5WN46yajp+KtbQb4eOEZFsez5F28jq/3RPJa6CrtGAtRfGtDOZbQSkD7EooS6RpmxPz66veZ
Q3ElewEtWY8gCIzFAd1EoENkDhFan3WD6ZYtJqWOpUv2JKBxeTx2iZUuK3BDccwLDaxw5fEwvRG89Np/
L6VBTno4G/4eKHO6E0udrnR6Zv6cT+0Dr+Pkj7UZN+0FkorioyzM9SwLFn6UBv3GbHvb7o6ncSSBx+i5
NJ1Eghg57k/YcFAtG/U4Iexw/Rwwz47ZexCgHbkjWO69aRDR5v+XzWY17voXggqVz3OCn/3sbpB4In3D
YhPL7PVmSGDsaMP8EaE/vv8nxsD4Fhhd58ND6zFGqiYkL5scOk7OC924Lh2Bf1jeF7mYOjJLBlPmO9KZ
+J+7s6El4bWO6/+DBQOuMoJ1s10IcBmqdcULil/89FwztC3IzrQ6AQMlouS36VHLLww9b3Ks9eQ4CI02
i/5INsvRGey6D8OGCrKb85pp8y/sMrXqb6RUY+8b3xpXbrhdR54vUcjxcrUdGjHU7/bRO1O4zO954aP3
LzFGKJiU6PvYBvKETqadZ5tLveXsrywzEqO8m0lZlwbtUDALjuE383Olr/iKNBvfhi+q1/s9lKIxpij4
j17TAlsoDNTK1Hm/f6XHwur3wcQ8J85QiX5yhn7IM8zLfmd79FZ6nS6d4SYngVYrlcyOYRm5xU9COzYL
VkjlOhb2wyM1sT7kYqAeDxlpOJ7aTyd/gg2lzkjE09fX2xwZS6MRGoij2SeHe5R+92+7e9TharwlH4KY
8VOCPOk7h+he3fpqOvIYi+Pd+2hfqcpPQsOXJ7gAQrNa5VzQnyR4Kw+cewzGuH2+6FR8CIzxkwkavRxP
9L/XyX8BAAD//4I5gC2DFAAA
`,
	},

	"/templates/build/ecs-stack.json": {
		local:   "templates/build/ecs-stack.json",
		size:    15295,
		modtime: 1449041232,
		compressed: `
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
`,
	},

	"/templates/build/vpc.json": {
		local:   "templates/build/vpc.json",
		size:    3348,
		modtime: 1449012824,
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
