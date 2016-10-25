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
		modtime: 1477309556,
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
		size:    12854,
		modtime: 1477368077,
		compressed: `
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
AAAA//+ogziBNjIAAA==
`,
	},

	"/templates/build/network-stack.json": {
		local:   "templates/build/network-stack.json",
		size:    4166,
		modtime: 1477309549,
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
