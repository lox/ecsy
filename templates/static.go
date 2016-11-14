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
		size:    13715,
		modtime: 1479088482,
		compressed: `
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
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
