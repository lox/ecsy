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
		size:    8296,
		modtime: 1479088482,
		compressed: `
H4sIAAAJbogA/+xZTW/bOBO++1cIOhZp46Z9X+zq5jpJ60U2a8ROeljkwNB0RFQivSTVIl3kv++QlmR+
SXJSdw+LoC1qk8P5nofD8d+jJEknnxdLUm4KpMg5FyVSN0RIylmaJenJ+O349fhX+Jseado5EqgkCghg
V5+GtSWSX85RSYuHdk2vPmyI5rBQgrJ7c9qsnxKJBd2oWsAyJ8naHE4YcE74OlGwpIBnQllSSZKak49H
O2GnZE0ZrTk8RyBdEabomhLRyDubLhKXdaJ4KB3IpkUlwf6DSq55RkVOOVOIMiIuwT3PlIobHloCUgrh
XH/S4iURXynukjnnQsVkXlblHfigT+YGjibajb58viHMPrlGVaGF/DJ2XX3x4UDSC45WyR0qEMNP02C6
WBBcCaoePgpebWKqQO1k2dn0JMsc0iybrTr1m4DTt7TJvSZO1lxsU+HiA/yPVJIjmSCMiZRNmOwsmTGp
tC3SDdrNBvdreDOf9uql/QY0rTqNRcnWekfYorpjRMmYwCkvS3RKClpSRVYXVKpOiTUXbSQWBPCndQNl
rrxPBBUqn+YEf7kWxTPL4PrqQovKqUq+5YQlKw4nktywxpq1TNaCl40SsRRJjz0/LC6mROiKxqA+OPfp
mk2YhwqIJcA2wTu+NTDYiRJVrtZtVOun63hlwMxC62tJPim10XEhzMWxcwaZ8leFCk3+Z72atPuG5oqs
jU2+3S3N41H7MW1Wbx2f1RrILhUuueqW36GjT7efrp6+js47vRufWpa0Hv6jUpvKqgMINsJf6si3Vt2g
ojKpQLB8DUEsicgy/TmKvlDrFwBZH2rEijHamWocMlt7zgh9bFvpOspw+A1KIXCo9saRt+JTAE2uBWXH
xwFtGJJW3keiJkpFJNY0n5bLueuFkLkhPL1cmIsxsnsbrD2GTNJsX7XrbGrupZD5qE+8/e3xJwfjJ8Ti
3wnF+/fv/NM9Xhz5q4979GqRCmoCu6Pv4reoy3Vffgu7vH1+8ba5m59F7/JrweiKSF4J3Ri0cATZ+twe
ZnfDzAV0TACgRLqKGTL/Sg9bm+19rq9WpvT9Bo6xUkj3Lc69aduse5rozeJoOmP3AnolL4O9wpptwAzF
MdfdQ6rwxkvj9Byufq/jDEIaq33/AlnyQ3CZ0pWY6XCl4zfmz/E43RNCnmIpVFyX+uHWkE4DBRngSGcu
Fkgqine00D5lWRSCdi2OPu33NoMZHLaxdqCa3eEMHMi9XeDdCtszoDagdB3f+d7SsXHEgHq2Zw+Ru80D
ZYCX+8rs5WgntE6jYdOtB4Mf3SUS9yRU7IevXq3YU7uJHh8Y257GzXslHapFqfk+LHOA2ZwXGq1PHIJr
lockb8cOzYzB1foV6Ri+c3eWtCS80hH5X7TUwEmMYF3kpwKcBWgw5wXFD35gzxi6K8jKoJ6AW7RDyP/H
I98hIVB1tuAHRaqwP3+Bqv8yVPV4LmJZeAP/fHUXfrPQM+HwxQ48tF8wO8btBbMPgtnx15n3xFnoJ86W
aBBvw7ejWe58LLrmw6uICrKa8opp9d/aSGcV+gDadc3+fW1cut567xrtd7IchhDboC5Fo1Oqeq93VhXz
i6tpeFv/cEn6D5UnFuUw0l3xwndTODYwRNHs6vwRzGbkEe1dP55ubg3NJr9nmaEYLKCJlFVpuG0r/5Rj
+M78pNMTU0XqjYGB7tl6DZhilCkK/i1IFNCFQg+2QUUkB00j1CBEbNSlp7NvUIm+c4a+yTeYl2Hg/UFW
kFnpBNcxiWS6VDLbuaV/1GUn0W00C+ZI5an5ScJe1L7exqIHWLYRqYs1tSfRfnvQFzpD0R2+UG59ZCiM
hqjHj2afbFtv/ePeXdt6b4dRd+RV14yy45Qg9xp8RPvTmh6R7DPzjPK78rl9pirfixs+2cMEIJpUKueC
fifROVjk3G3Ux83AUIfi1YFQzkWXkf73OPonAAD//47tE9FoIAAA
`,
	},

	"/templates/build/ecs-stack.json": {
		local:   "templates/build/ecs-stack.json",
		size:    14118,
		modtime: 1479095315,
		compressed: `
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
Tb8sZ957+fdyb7n3vwAAAP//ZShhqiY3AAA=
`,
	},

	"/templates/build/network-stack.json": {
		local:   "templates/build/network-stack.json",
		size:    4166,
		modtime: 1479088482,
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
