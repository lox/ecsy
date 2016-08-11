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

	"/templates/build/ecs-service.json": {
		local:   "templates/build/ecs-service.json",
		size:    8296,
		modtime: 1470892715,
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
		size:    12792,
		modtime: 1470892715,
		compressed: `
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
`,
	},

	"/templates/build/network-stack.json": {
		local:   "templates/build/network-stack.json",
		size:    4166,
		modtime: 1470892716,
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
