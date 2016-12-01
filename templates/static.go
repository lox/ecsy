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
		size:    5132,
		modtime: 1480579508,
		compressed: `
H4sIAAAJbogA/+xXW2/bOBN916+YGgUKFHVu7fdhV2+O424NeLNG5KQPRR9omoqIyqSWpBqkRf77DkXJ
IiU5cYLdt94AmzM8h+ScuXg8HkeTz8mKbYucGPZRqi0xN0xpLkUMb85OTk/GJ7/jvzfRBdNU8cI4y2ya
QMLUd05ZDJPm4zsgsCL6G1ywlAtufYGIDf6H2eL8TRQtiSJbZpAgjgBuCjrf2A8Aq/vCAn1O4ng2PYvj
m+U0juebyhYwrzIGfMOE4SlnCmQK6ApGgioFcBE52GW5zjlNyrVgZh+Fs+5nKSoM0JUbQoPBRUumC0Yt
+wbuuMncLRpexb/jOzro05cyp1xpA4XD6hzgUPKzl5JrRiWG7HnsKIdpXmoMrE+aGMXF7TCNFRB1W2z0
iDGEZo7DaQlXa+SE0VJxc/+HkmWx51KBy+DdJvZmlRPcWi/kIgYoSpNQyrSuTsSFNkTgV0ttlfyRbHl+
f+il0sobBGrcKtPexth0wNcrNfsvMduEOxQ3TCJThyTEspGpSaZSGMIFU5d4kEM5aLPJD7H0oxxAL6Uy
PvRluV0ztSc70RekU2VAIwsm6i0pKXMTw28nlZAW5y/FzyXZwJrkVhmPcnxiJDfZNGP027XKD32l66uF
Bc24gbuMCdhIdIWswqIWS0Oq5NbFaHEeEh9b3iRZTJmy0aSYs2Hm72OeiI4EMBMQB2gLVIcfUqmGyUej
KMLobSqxaMd6rdknY4oFx9QWTT0AeDX7uyS5hi/w6oqlvRO/Qyz4GvkIugdxKY3dfwiSxfqrNEVpqmMl
KL1v1XtUWDckL/FpRozqMV5uixyx/VyLcjRU0uo9FWdrqz0XqJDzWiCh+zyFL/X5AUbdq43e7WyvsDDD
KLPW+Pj49c9Pq9UywD26uExs8j3Er3/WYn4Y3L/b/uTuevPXXaF1DX3ovkmbr90a5nu2tii6YlqWCotp
7NLvuXW88lgqTDYML9ONDKAyeloGFHlY2q1gkdDWSMwcYay88Q4NQNW2oD6x62G1JeCfi1uFjWHHC2OY
F3geI6nMYzC02FkAPmKGVvWlfjH3xJ7DSj5qnvKNmhcxnBxVf49Pnsn64cP7AbJwtc9ha1ZHKP3A5EQb
TlsfLCdx7G+pduzqQNwtAXvi6AaR4HmbgITjWzQQm/62rsB69s4gUdubQwaA/uUejdq8nhc8p6Cbea5t
CO2L1wavYbT8K6JumYWrMtp6Y84GsA+Y30Gr2aVyA3m/ylC7mcw3MZztbNci61lPW6XNBVa07wSP+L5d
XPEtkyUe53/1Ep5EMGojfaHwRKiGpcRo3bcXmAmyzhliG1WyPtD/d7obqJv/hvD0L+U9qbxhsrBevJgj
8Sy92WS4bf9KiG4DDnpjYnuj68D1c3c03cwqvQHFmbFdcsU2U1kKJDxtJOgJIJBhOO13Qm/X/K4SjO+P
y8QnrKHtgDQ0GDmgXpXw1v1lnGEcw5XM+1OLXYwaLQU/ltqppV0Mg1EBdgMyn/wZxzvUfjQmWpfbaquT
woWk+F2Y9o1xIjUsXLIPP0tTlBKy5Lm88yyWhGMFLEgeB8sAgWj8P2PAmfaIbMkPKcidPqJyG/hMaPt7
0d+lsTS2F6jNS2Iy+2uj/mZv5c9kuM3d1EXVG6b9IrHnLfa+x1Nvsv8W9f1dB7G/39a7DuImxzV7+5wt
it1aZaqmJGo7efVa0UFgV12oz9xkT0PRs6dOjh6T0mRS8R9saJDt7WoGdPwl9HYU/RMAAP//84MFVwwU
AAA=
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    9604,
		modtime: 1480310754,
		compressed: `
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
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    3175,
		modtime: 1480310754,
		compressed: `
H4sIAAAJbogA/7xW3U7bTBC9z1MM0Sf5hvDZLqpU36VJoFFbiJKISlS92DgLrHB2V941NEJ5947XG/8l
NgHEnxSxc+bsmeOZcXq9Xqf/azanKxkRTc9EvCL6isaKCR6A47ue23O/4J/TGVIVxkzqLDIazOBqMoAL
qh9FfB+A4BRksohYCCpZcKqB8CVgEGTMHpDbHqtjIGEslDKx/rU6cTqdn0RKxm9V0AGYGdhA8Bt2m/4P
6T0BPMFgPJwG4Lkn5vd/7zNsTDhLcGsQHyH+aQXi1SCfdiF+DXK6hXQuEy0TbRReyXC8tNJIlNAAjqb0
JpVpzkZ/pYh1Fge4IKsUgPTQ/e8JvQ6CmSbhfXq+6Rmqbicv250YC7PkquMWYE1G1+r3VxjalTi7SirZ
TqHIm2SPr1mSt33CjZq2HK8UtU0vqfKfVeU/q8p/myq/UDWlSiRxSLPuwG41NPO1RBKTPBr4QbDtj0ks
JI01y+Dpz4At46+RCHGMjs4YX445DgT8rozCMXSRoIsfaXt24Y/NHXOlCQ/pnHL8WAewpDckibQNjzhZ
RHTI1SyRpkLQcULrwW9CaY61qUp4Tm5ziQA9+E6RPrUgP6va2q0aZdr6HB16JOv9joy5pjFWaEF1d+Bp
U6LoayS+W1GuG+3dQdrukJQv1SX2RsNFuZUVOTjjWV3lrO30FyOPx6PxpKzJKURhxGm4bShWhKGmBxmm
HBdE7/GqRFUAmhj7EfYQSafACLTdixJOypGNY+FZd+W1VNfH/o1U9zyDNMjBHs5yx/KS/yAJD+8q7dV/
ICwiCxYxvb7G10cqmUY01Nj57jEcnVONrwdwnLzXD5wTK3x3VvY8un1L7kVVvnMV3guq8N+nCu/tVfgH
VTEViaaqrd8MYp4urYY6mlnbrHkN7TBbs3YkDUMLedMOoEozbgaz5Of2G45rUe37qBCfA8pGVnfLK3U2
31F6hR5aTqEmJyuOSosnu6GtmQtZfaVEyMzNDSW0bruDnbRD+QHSvKq1zz6BYgt8gDj/heL+BQAA//9Q
AonhZwwAAA==
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
