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
		size:    5086,
		modtime: 1480643635,
		compressed: `
H4sIAAAJbogA/+xXS2/bOBC+61dMjQIFijqvdhe7ujlOujXgzRqRkx6KHmiaiojKpJakGqRF//sOJUoi
JTkvdG99AQln+H0czseZ0XQ6jWYfkzXbFTkx7L1UO2KumdJcihhenRwdH02P/sR/r6Izpqnihakt5/ME
Eqa+cspimDU/vgECa6K/wBlLueDWF4jY4n84X56+iqIVUWTHDBLEEcB1QRdb+wPA+q6wQB+TOD6fn8Tx
9Woex4ttZQuY1xkDvmXC8JQzBTIFdAUjQZUCuIhq2FW5yTlNyo1g5ngfR23eT5NypQ0UFRToyhkZwKDJ
cuqCUXuILdxyk9XBjNCfPJdeMyrx7p7Mj6mZ56XGS/Z5E6O4uBlnssmk9RZ7k8QYQrOKR9d5xVWHnDBa
Km7u/lKyLPbEFbiMhjezwVVOcGO9kIsYoCgTQinTujoRF9oQgb9aaquq92TH87vHBpVW3iBQb1YlNhpj
pYk3WGr2f2J24n8sbiho41ISYtnMOJK5FIZwwdQFHuSxHLTZ5KdY+lkOoFdSGR/6otxtmBqHLtAXZK3M
gEYWTLgtKSlzE8MfR5WQlqfPxc8l2cKG5FYZ93J8YCQ32Txj9MuVyh97S1eXSwuacQO3GROwlegKWYVF
LZaGVMldnaPlaUh8aHmTZDlnymaTYkUNH/8+5pnoSQBfAuIA7YBc+iGVapx8MokizN62EouuWa80+2BM
seT4tEVTDwBenP9bklzDJ3hxydLBid8gFnyOfAQ9gLiQxu5/DJLF+qc0RWmqYyUovS/VfVRY1yQv8Wom
jOopBrdDjtj+7EQ5GStpbk/F2dmc5xIVcuoEErovUvjkzg8w6Yc2edPaXmBthklmrfHh4cvvH9brVYB7
cHaR2Mf3I3753Yn5x+j+dvuDu93mz22hrZvrWLxJ9177Ncz37GxRdMm0LBUW07h+fk+t45XHSuFjw/Qy
3cgAKqOnZUCRh6XdChYJbY3ElyOMlTfG0ABUbQvciese5iwB/0LcKGwMLS9MYVHgeYykMo/B0KK1ALzH
F1rVF3dj9RV7Dmt5r3nOt2pRxHB0UP09PHoi67t3b0fIwtUhh61ZPaEME5MTbTjtfLCcxLG/pdrR1oG4
XwL25LGeRYLrbRLSG6Ue8DhpPYL0DZH7GhzYe7OGszdxBIB+/PcmduFGCs8paHiea5dlmxRn8HpKx78m
6oZZuOrRW2981gHsDywBQTdqX3sDebfOUN6ZzLcxnLS2K5ENrMedGBcCi95Xgkd82y2u+Y7JEo/zm1vC
kwhGrRjOFJ4IBbOSmK67LoBzQTY5Q2yjSjYE+r2V5khp/Rna1L/E+TPEOU4WVp1ncySeZTDhjDf/X2+m
38aDDpvYDlv3cXfdPdk3E89gzKnN2HS5Ytu5LAUSHjcS9AQQyDD8Zuil3q75vSn4CLhfJj6hg7Zj1th4
VQMNCom37i/jJFQzXMp8OPvYxajRUvDJ1c0+3WKYjAqwn5DF7O84blGH2ZhpXe6qrbUUziTF34Xp7hjn
WsPCJXvx52mKUkKWPJe3nsWScCySBcnjYBkgEI3/Zwo4GR+QHfkmBbnVB1TuAp8Z7b46/V0aq2cXgDOv
iMnsN4v7zUblT3a4rY60zqo3kvtFYs9d7L2Ph+5kfxQu/rrJ2K/ATdtk6vlzw14/ZYtiN1aZqimJ2s5v
g271KLDLPtRHbrKHoejJQydHj1lpMqn4NzY2Dg92NWM+fk+9nkT/BQAA//9COchP3hMAAA==
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    9484,
		modtime: 1480643714,
		compressed: `
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
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    3730,
		modtime: 1480643528,
		compressed: `
H4sIAAAJbogA/7xWXW/aMBR951fcokl5gS4f3aTljQHt0LYWAeqkTnswwW2tBieKnXao4r/vxglJHAhf
XWkrpfgeH597fO8l7Xa70fk1ntB56BNJL4NoTuQtjQQLuAuGbVpm2/yCf0ajR4UXsVCmkX53DLfDLlxT
+RJETy7gA8J46jOvlf4fsWdkBBFPOZWiBcSLAiFUrHMnzo1G4ycJQ8YfhNsAGCtYN+D37CH5DAm7C6/Q
HfRGLljmufr9aH2GpQqnG8wKxEaIfaFBrArEWYfYFcjFOsSpQD6tII2bWIaxVEncht5glqknfkxdOBvR
+yQTtdb/GwaRTOMA12SeAJAemh9e8RJcdyyJ95SsL9uKqtnInTGHytt0s34VGSBzH42tnq8xbFdirCvR
dhuFImuXImuXIutNiqw1RfYwrbl6SfaqLGs1rTiOFLXaXlLl7FTl7FTlvE2VU6gaURHEkUfTesUWUzST
RYgkanO/a7vuqmKHURDSSLIUnvx02Sz66gcedvzZJeOzAccuht9a/7agiQRNfCQN04Q/2d4BF5Jwj04o
x8fChRm9J7Evs3Cfk6lPe1yM41BlCDKKaTX4LRCSY25CC0/IQy4RoA3fKdInFuRruq1N3SjVaFfo0AtZ
bHZkwCWNMMMMVHUHXpclio5E4sc55bLW3jVkVh0h5TNxg7VRc1BupSYHp06aV3nXah4VQwiX+4NhWZNR
iMKIUXNaL5gThpqeQy/huCZyg1clqgJQx9jxsYZI0gVKYFa9KOG8HFkaGTytrjwXfaBtnpFVz1NIjRys
4XTvILzhP0jMvUetvDrPhPlkynwmF3cBVw1HfepJrHyzBWdXVOJ3GhhGXut79kkmfL1XNlzdhrF7siSt
tydpHZCkPskPyvKdr8o+IAvnfbL4D3fh7JXFKIglFdvqTSEmyWSuyaOedZs1x9D20u+SbO4ohi3kdYOO
Csm4mj4lP1fvnmaG2j50C/E5oGykPkCP1Fl/Ruk9Yd90CjU5WbFUmq7pCduKuZDVESLwmDq5JoWtI31v
J7PRcgJp1nHS7BNIs/Vb31kcxYA6gTjnQHH/AgAA//9NTWifkg4AAA==
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
