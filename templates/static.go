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
		modtime: 1448600177,
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
		size:    12360,
		modtime: 1448595006,
		compressed: `
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
`,
	},

	"/templates/build/vpc.json": {
		local:   "templates/build/vpc.json",
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

	"/templates": {
		isDir: true,
		local: "/templates",
	},

	"/templates/build": {
		isDir: true,
		local: "/templates/build",
	},
}
