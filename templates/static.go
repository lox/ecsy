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
		size:    5113,
		modtime: 1480305617,
		compressed: `
H4sIAAAJbogA/+xXTW/bOBO+61dMjAIFijpf7ftiVzfHcbcGvFkjctJD0QNN0xFRmdSSVIO06H/foUhZ
pCUnTrB76xeQcIbPw+E8MxwNh8Nk9ClbsE1ZEMM+SLUh5pYpzaVI4fX56dnp8PR3/Pc6uWSaKl4aZ5mM
M8iY+sYpS2HU/PgWCCyI/gqXbM0Ft75AxAr/w2R28TpJ5kSRDTNIkCYAtyWdruwPAIuHEoE+ZWk6GZ+n
6e18nKbTVW2KiBc5A75iwvA1ZwrkGtAVjARVCeAicajzallwmlVLwUzMMNpSOOt+lrLGAF27ITQYXLRk
umTUsq/gnpvcBdHwKv4Nr9FBn72Uec2VNlA6rJ0DHEp+/lJyzajEjD2PHdUwLiqNeQ1JM6O4uOunsfqh
bovNHjGG0NxxOCnhqkeeCm2IoCxjtFLcPPyhZFXuCS5y6Y1xZCOsneDOeiEnMUBRoYRSpnV9Mu4ptT2C
FfQHsuHFw6HBrWtvECh1q1AblbFVgbdYafZfYrZ1dyhuXEzGpybGshnyJGMpDOGCqSs8yKEctNkUplqG
2Y6g51KZEPqq2iyZ2lOl6AvSqTOikSUTfsuaVIVJ4bfTWlCzi5fiF5KsYEkKq4xHOT4yUph8nDP69UYV
h97SzfXMgubcwH3OBKwkukJeY1GLpWGt5MblaHYRE59Y3iybjZmy2aRYu3EH2Mc8EjsSwEpAHKAtkE8/
rKXqJx8MkgSzt6rFoh3rjWYfjSlnHEtcNH0B4Gjyd0UKDZ/h6JqtOyd+i1jwJQkRdAfiShq7/xAki/VX
ZcrK1MfKUHpf6/uosW5JUeHVDBjVQwxugxyp/dmLcuAb0AzzfuHTHm08mq7hsz8VwGD3wIO3W9sRtl0Y
5Naanpy8+vFxsZhHuMeXV5ktqZ/pqx9eoj9792+3P7nbb/7io2he6ygAe3Otra8zhZ6tLUmumZaVwhaZ
uqJ6bneuPeYKSwiTxnSTXKiNgUIBpRs3bCtDJLSdD+tBGCtajKEBqB8l8Cd2L5S3RPxTcaew3W95YQjT
Es9jJJVFCoaWWwvAB6y7umv4G3NXHDgs5KPmMV+paZnC6XH99+T0mazv37/rIYtXuxy2E+0IpZuYgmjD
aeuDTSJNwy31jm11p7uFvSePbsyIrrdJSDycJT256W7bFVjHvmdM8H7NYSPgMMhHs9cgB07RWxW4tqm0
N+8NwXPQ8i+IumMWrq5s6421G8H+xDqPHpJtSTeQD4scNZzLYpXC+dZ2I/KO9axV3FTgzPWN4BHftYsL
vmGywuP8zy/hSQSjNuOXCk+EqphLzNpDG8BEkGXBENuoinWB/r/VX0///DcEqH8p8GAF9pPF/ePFHFlg
6Uwg/Y/zr8LYfZCjtzKzb6V7kf1172i7+djaSs8veDM+n1yx1VhWAgnPGgkGAohkGM/0O6m3a+ErEw3p
j8skJPTQdmDqG5QcUKdbBOvhMs40juFaFt0pxi4mjZaiT6J2imkX42TUgLsJmY7+TNMtajcbI62rTb3V
SeFSUvxdmPaOce40LF6yFz9Zr1FKyFIU8j6wWBKOnbAkRRotA0SiCf8MASfXY7Ih36Ug9/qYyk3kM6Lt
V2G4S2OLbAPw5jkxuf2m8L/ZqMIZDbe5SF1Wg5E5bBJ77mLvfTx1J/uj8PG7l8R+pS23L4mbJJfszXO2
KHZnlamalqjtJNZ5kg4Cu96F+sRN/jQUPX/q5OgxqkwuFf/O+gbbzq5mYMfvnTeD5J8AAAD///jNgd75
EwAA
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    9603,
		modtime: 1480305502,
		compressed: `
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
AA==
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    3175,
		modtime: 1480305481,
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
