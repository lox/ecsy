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

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
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

func (dir _escDirectory) Open(name string) (http.File, error) {
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
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
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
		size:    5467,
		modtime: 1491293338,
		compressed: `
H4sIAAAAAAAA/7xYXW/UOBe+z68wERIS6sB0Xl4JfLFSNi3bkdjuqJktF4gL13PSWCR21naoCuK/r2wn
M3Y+5qPsrm+Yxuc8z/H58jGz2SxKPmZrqOqSaHgvZEX0LUjFBMfoxWJ+Pp/N383m715EF6CoZLW2O79E
CCF0mWYoA/mVUcAo6X4iwjeIoDVRX9AF5IwzoxNFKyJJBRqkwlb7tqbLjftp1vqxNigfM4wv0wXGt6sU
4+Vmux/wrwtAbANcs5yBRCJHt6sUaYFkwxHjUUewau5KRrPmjoM+38fmRPYT5kwqjWoLiZRVQIwjXYBl
VzVQY84GPTBduOONGrL4WUMUUME3T7HkMs3SslEaZN+CTEvG76c5TaypUzV+JloTWlhG1YZdiy3HlVA6
A9pIph9/k6Kp95w3EJs8dmIObQXRvZFEuiAaUcIRoRSUsvYxrjThFJQz5IMgm19Jab7If9maXEjri1KQ
DbprOZ0Vpg7ek4qVj6e6PLdaiJMKTIYbfG2KinHUKPiv0HclfCpDWKC6TaIQ0+TSli4VXBPGQV6TCk5l
o52yn57Cz9AeyUpI3Se5bqo7kNMktZAaCVdpAaGogXtqOWlKjdHbeZuJTGn4ByiD9DpMm2UfUpAmBJRo
SORJEUw4Sm6uTeQIN0iI7qDaqE2k/cCeOHb2XIBiEjapaPjJbuB235izrfIgbV3vH3KfR1Eq+MYmW3vx
/KngSuu6C8rOkmeXfzWkVOgTenYD+dB9ZyiO0efIB1EjKNdCG4ijwDq4pbJordlhnjxL+AZ92v75BI6z
UDnU8pPzDMVv3vwv9jUCA4+27yjDTjHr7bxv1R+NrhvdxjTThH6xmbQVuSVlAxjFQNUsF7ICibH53XaD
ePpGbDWtGbv9rbx/rwyVlnnPF2bF/XSJzwYy45rjyTHUtiyFEcOvXz//7tv46uI6Mx31xwjpMWr4+Xc/
GD/iAcrn009z3GGedpafOkoU/tUFvZtzp5Ik8++YsVvZl9/tR9ENKNFICm0ijyeXP6eURGlGd3KM398u
MPYVt3orKWpTeB16t9yU2fto1swZOBydj5Rc9CSDgWqab3JU6/wp70Hvmd0mfOLpHXSJe45sD9U7h+14
bjMYIEIhKbSgosToar1eBVtXQEpdpAXQL0uuQX4lZWaneIXR+XxKdAe4Tifx1qwC0egt3P9HBB/XhQRV
iLK9fNFiOCIPb7Jjsq5VOuhen8lMIsO4T/nbL9gJd9tmM9ZkbSAy90+vtL0LaSQvpxEnLzTrp2th6zzu
kbXdLqHeGOKvWevrXMgHIofFFpTAzn1+fk82quCNk5k3jutUhyLWXYuDezA8mDfSOUn/02QKjPognP97
1Wa+jfglHOcPVeiRvvSFb0Q5bPPmYxRCBi+lXZv3/wckjJAFHovSMvkd44BhKkSJUk1lgVaiZPTxQtCm
An+47lamiYbxLef7yzwHqjFKylI8jMoYMxinrCblOEjL1M+/qTVDQNUrUpFvgpMH9YqKao+OK51jUJVW
eOeYsGMQXWD0utdpSkYHnnVQzqsuG73RccSIQ/5360AUOtpjYuHW8V4B18DNS+2ua+DYvbLu4OXTASTc
m94ol92r7L0U1WRbPwH6pg/8keniVGC6OO2MdIGTRhdCsm8QzCFLfi9Bqb0Y3SSHUfwyjv4OAAD//9ED
oWxbFQAA
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    12394,
		modtime: 1491292177,
		compressed: `
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
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    4278,
		modtime: 1491290422,
		compressed: `
H4sIAAAAAAAA/8RXS2/iOhTe8ytO0ZWyAZqE3ivVu9xAO9FM2wgQlVrNwgS3jRpsK3baQRX/fZQXeUJC
yzDpprI/H3/feeJ+v98x7qczsuIeluSK+Sss58QXLqMIFF3V1L562Vcvlc6ICMd3uYx3xuYU5rYJt0S+
M/8VgXxnwIOF5zq9+H/ffcOSgAgWlEjRA+z4TIhoz3gQA6XTucGcu/RZoA4AwDQCmow+uc/xSvjNbRPB
B5jWaIJAUwfR37n2H2y2kPigWoLpA/Vcv6jAtBJsWA/TS7CLetiwBPs3hXXuAskDmUibc8da5jRhLyAI
zibkKdS3XR//4syXGS78bvEqhE6DBXT/+TDupwhNJXZew/VNPzLc7eT8p9pRDDIjxbAloCRSA6WWU8FS
W3ZKlV3BjpJnqbVhqbVhqR2JpVbDUrfjHN5PU09TfS/P1NaXiaaGCkyHrZgOWzEdHovpMGM6IYIFvkPS
erDNzNxszQmCyMTY1BHKV4TtM0586aYH0890l/7/HnNeEZxduXRp0RvM4bHQQ3rQndtmtwfdsEC78LNg
waJCYuqQGaGYOmsES/KEA08WQGOKFx4ZUTENeOQDkH5A6iHfmJAUr4ioAc3wc0lA+PXhO1mjyJWVvWJg
ukUHJwV/jSV5x+vdnrSoJD4lMgHWeRU+NgVjhpTYeVkRKvcGqILO5R0ndCnuKII915ZCUaBpLRPV5fOw
7aRZ64w2x5ZdZqtkdMeWrTRyGLEVdimCN+7ENm+x3OHfnOkM1HyD4XnMwWE9RgKSyhlb9iC/s1EKh+J8
3iouNuZ9Xb8csxjWSPIG89iSxe/oDxxQ56UmnY037Hp44XquXD8wGjUC4hFHwiOoPTi7JtJ4EKAopZpr
WbWJpF2VW5sCO8bKX3ODdiw3aJ9wQ3VufdIPJwm0/gmF1Xn3BxQeLYbDgxROWCCJaMrkCDULZ0+jxn23
NLnxq9eM4qmadMnIWsNlza2aCOnSqF/mopC+DNQCts1IySRuYfkQlMfBUTTsvrP02+tQwRnPreFsqTAx
4vuaCigjagjBHDdi0Shv79g62PdJIzw54fLT4kDC+skJV94YrVIt31ZPTrny2GhF+XcAAAD//6U5GP+2
EAAA
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/templates": {
		isDir: true,
		local: "templates",
	},

	"/templates/src": {
		isDir: true,
		local: "templates/src",
	},
}
