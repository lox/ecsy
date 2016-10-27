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
		modtime: 1477520712,
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
		modtime: 1477531980,
		compressed: `
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
`,
	},

	"/templates/build/network-stack.json": {
		local:   "templates/build/network-stack.json",
		size:    4166,
		modtime: 1477520116,
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
