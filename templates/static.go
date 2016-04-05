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
		modtime: 1456986937,
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
		size:    12230,
		modtime: 1456373101,
		compressed: `
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
`,
	},

	"/templates/build/network-stack.json": {
		local:   "templates/build/network-stack.json",
		size:    4166,
		modtime: 1456373101,
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
