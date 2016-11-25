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
		size:    4938,
		modtime: 1480055179,
		compressed: `
H4sIAAAJbogA/+xX308jtxN/379iiE6KdCKQ4+77Veu3EEIvUkojNtw9nO7B2Tisxa69tb0gDvG/d7w/
7d0AgbaqVPUOROKZ+cx4fns0GgWTr+GKpVlCDTuXKqXmC1OaS0FgeDL+MB6Nf8afYXDGdKR4ZkrKbBpC
yNQtjxiBSf3xECisqL6BM7blglteoGKDvzBbnA6DYEkVTZlBBSQAuGDmTqqb0NDoxn4HWN1niBcaxcV1
ceBpXcUMBMqD3IKxn0t50BYgQH6r+5ymPLnfF25bcHuoxl6AC8g1+zsxWxfti8s3TBi+5UzVsDYKPhYY
WStB4jTJNTr7zymoQBzkqRSGcsHUBV5xX/CoFrJA1GDIYvvJatFl9njQS6mMC32Rp2umdkNnyAv27l01
MmOiEtnSPDEEfhoXrlmcvhU/kXQDa5pQEb2gYxqGLMoVN/e/KJlnrjIsOUJm0xNCfBYy3/QtmKB/Sia4
tlywlaoMzuIU/1IDMdVAo4hpXXvUjdtcYH2gtdpa9ZnRxMTTmEU3VyrZN3ZXlwuLHHMDdzETsJHICnGB
FVksDVsl09oq3x3HVm8YLqZM2eSKsM/MN/tonohORmIbQRyIWqAqKV2X+MoHgyDAnNoUxaFLrVeafTYm
W3B0j6irA+Bg9ntOEw3f4OCSbXsWHyIWfA9cBN2DuJDGyu+DZLF+y02Wm8KsogkW/iiwvtAkR9cMWKRH
eLkUdRD7uSqVQZVhC8zG0yoZPcGD+Ra+VVYBDLoGDw4b2kGYr2EQWyo5Pn738Hm1Wnq4R2cXoS30R/Lu
oSqcx53yjfiL0pXw911tqr6A9VxLayqqHDm7OMO2i3R7tsvZ0oLgkmmZKywNUjaF1xZswbFU2AIwvEzX
aQAF0cllwCTv17CtX5wJWDnC2PTGO9QAX7JovkGpBxieC2wLqW1BxS2GpPT48N2DOz4fR4XIEB5rCM/Q
ubhW2B4aA2EE8wwNNzKSCQETZQ0F4BxLuWiPlWvLqDkMK/kseco3ap4RGB8V/4/Hr9T66dPHHcr8074O
29w6udePYEK14VHLg32HEFekkGgaBun2iicCjgHBXcRz7ysiV0qPl/k64VEbQS+AmtTdpJumh04BOMdY
XCVMbbpnnXvlZ2NZDw+HyRvRDmsbWBuHiuDMm1b/iqprZip3DCw3NgcP9hEbiTepmp5RQ96vYszoWCYb
AicN7UrEPeqHNv/mAlvJLUUTP7aHK54ymaM5/6uO0BLBIhv/M4UWYY4sJUbmvr3ATNB1whDbqJz1gf7f
ZOOOBv1XpKP+J/LxTSgnS8Vvcer9y9J6tzK/Rb1ZR+hQenvT7pXiv2rrLgfe3A7t3C63g8rdnYKpN5De
2lGScZRzxTZTmQtU+KFOQScBvDT030ed0Nszd5B5D57n08RVWEHbNW/XelcC9VqQc+4eN1V1KZP+RmUP
gzqXvIdru1G1h34wCsBuQOaTXwlpUPvRmGidp4VomQpnMsLvwrQ+xh5jmH9kHT/bbjGVUEuSyDuHYpVw
bK8ZTYh3DOAljftvBLhvH9GU/pCC3umjSKYezyRq3+6ulMa+216gIi+pie1LqPpmb+XuiyhW3rSMqrPo
u03iCV886Y+XfPL0Lar7l+PJvnjXzXgqt9o1e/8aEcWubWaq5ilql73enNsL7LIL9ZWb+GWo6OQly5Fj
kptYKv6D7dqde1L14wFfae8HwR8BAAD//3OiUDhKEwAA
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    9295,
		modtime: 1480054129,
		compressed: `
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
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    3960,
		modtime: 1480054203,
		compressed: `
H4sIAAAJbogA/7xW3W7jNhO911NMjA/QTZzP0m4LVHeu7d012iaGbaTAFr1gpHFCrEwKIuWtG+TdO6Ko
X1t2nMX+BLDNORyeOZo51HA4dMZ/rta4TWKm8YNMt0zfY6q4FAG4/sgbDUe/0J/rTFGFKU90EZlNVnC/
mMAt6q8y/RKAFAhJ9hDzEFT2IFADExFQEJKU7yi3XVbXwMJUKmVi48/qxnWcBUvZFjWdGzgAnynXLf3O
vwOs9wkGsNIpF49moUVk/YQgCAtyA5q+c0FZBIvhX0riOH+wJKF9Ju3KnD+RYsMfi9RUQADPMJlPlwF4
oxvz///ez/BiwsWGUQfiE8R/34J4Hci7Q4jfgbwvIc5dppNMG4b3STiPLDUWZ1T21RI3+bJZm/2TyFQX
cQCjEFxRehj875keYhCsNAu/5OsvQ5Nq4BB2bhXJVW3nPr71ppRwcPGhzaPM2VbBhemL4PDpuRZgO4da
oVt7K8NpQu4hodZut2bkLYqe7KfklW3by6nM8UZS5fYGK/8sK/8sK//bWPk1qyUqmaUhlp3ZnEazeTbx
g4BGyKwvUplgqnkBz/9NeJT+GsuQvOHqAxfRXNAwwl+tMbyGASUY0Ec+GgP42+6dC6WZCHGNgj72AUS4
YVmsbXgm2EOMU6FWWWIqBJ1m2A1+kkrn1qBa4TV7rCgCDOE3pPS5BNVaW9ZBWyjT1nlmjOqBaoiylJnG
n94FQY3p0acGWEvKDUJutyioniJonBA2MoUVpjseIky5CuUO0701lzMT2Zll43it6uk3OU7bZirkEh9N
8zV0KJaMCNNPk8WdaU91vDUagB4JpnLLuLisgnpProm9McpyxltGvk8H7XiE0fR2VRIdKyVDzsww9fVx
g28D3kf9KUwsuBKwW3Bp6LW8tPyRpusr2x9nURgoagvqng3PL40UY00qPZl+6SvpAGmdJUERqTt6tD0H
VWPYolOV0tx1tMjZfNHk5NakKOKe7IYAdkWOW6aPaNVIVQP6Mo5j8h/zFA1B63xE4aYZeXEtvHCmqpb2
1XP8NutqXkB66JD/FXvnyZ34nWUifGpZ03jHeMweeMz13pgLUcYYQ02uObqGq4+o6X0JXLfyyVd6rCV+
6LNHHt2xC/KiKr9zFd4FVfjfpwrv26vwX1WFuUzUqX4ziHV+4fXU0Z/1lDRvSTstrmg7kibDieR9HoBK
c2EGs6Fn+WY+sqjTflSTrwBNIdve8kae/Wc0Xr9eW07NpkpWLzWMpzjhVDPXtM7fYCfd7tVK2qH8AdS8
trRnn0DtAj+AnH8huf8CAAD//0kNmPF4DwAA
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
