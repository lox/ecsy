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
		size:    12230,
		modtime: 1449462966,
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
