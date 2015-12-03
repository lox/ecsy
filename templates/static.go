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
		size:    11355,
		modtime: 1449104276,
		compressed: `
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
`,
	},

	"/templates/build/network-stack.json": {
		local:   "templates/build/network-stack.json",
		size:    4166,
		modtime: 1449109612,
		compressed: `
H4sIAAAJbogA/9SX0W+bPhDH3/NXID/n9yvQtdJ4y9olYpO2KEGZ1GkPDrlmqAQj22yKJv732SQktgMx
UZO1VVtR7Psed5+zD/On5zho8G0awSpPMYchoSvMZ0BZQjIUOMh3Pfc/9734RX1p+7XgecGZmJJSMTDL
492NvMVpAcqAGJrAo/QkDbeDZXUt+xsP02KegeKy2cswC4JPJJFBfd+NivE+6iu36pSjyPVQNk90kTZd
9k/Qes/Q+oZWufvRM/8zUEFc0ISvR5QUeTfsukT32tt6FraMFDSG9rpG61z6l2slCD7e+UEwG9/t0KMx
JTlQngDTg7hLFvRDSuInqfXc/6ufK+9WKRoKM8ZxFkMEmbispeUCHnGRctUqwktm1F7HjD5Dpf2CV4D0
ejTg0SFVWU05jp8qdWt9Skt9RmIL/cbr4+DCjAMV66A2Pg6x8QEDLmL9uYKMW2t0qLAWzQgwXBjkamp1
ArspZRvI9dMqVFrBnmlTU3CPp7cxsic0+IWTFM+TVOyDB5KZK6HqLVNIIebGCnMc17KzpXQEfPDAGpaX
scAmsJRd1TAqO3QDFezhlvLElvLfoXPD914cvvcG4PuXge+/OHz/DcC/Pjf8CSm4xrCBfWUT4XkKdv7n
COd++zK0B2WP5x4YTzLMRTU0mvWb2VVZPqv/7ym16rewu70JulRmY7p/8oAxEidVtnY0G3FrsAdnxgsm
672OZL1/kqz/OpL1L51s29G9IVnN1JpeZSY2dkyTnG8/28TB70r4cmpXju7rpM6kiLTIwmxJgVkO5WEu
wuYkJql0zOPcPJoPKVmNCZXtzTdeOSgibTOyeYW53rkUC9s5vSWPrjWpFd1K077wmj7NdODH8KnoXPVL
qaZ2e3NzfaNWr/rG0zM5NTYdZ0/+lb2/AQAA//8njafhRhAAAA==
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
