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
		size:    9456,
		modtime: 1480647478,
		compressed: `
H4sIAAAJbogA/8xa+2/btvb/3X8F5xUzMFSW7bh5CMu+X89xNmNLYsRuh90HUpqibaISKZBUHs3N/34P
ScnWy0naXuDedU2lww/Pi4fnHFLxPK81+nO+oHESYU3PhYyx/kClYoIHqDPo9Xte7wT+77TOqCKSJdqN
TMZzNI5SpakMEE61UARHjK8R40pjTqhCmIfwFxWQ3U6rNcMSxxReVNBC6ENCpqF5QGjxkNAAgTJBMBkP
guDDbBwE09COlWQvNhSxkHLNVoxKJFYIoEgLJFMO4luO7UyyW7Boni451f19QtzwfjkrJpVGieOFlEWD
DKRhTCWUGBVCdMf0xpnSJHzwtcIVJQJ8+GXSf6cPl+DgYA9LtUGf6EOCmUSpgsngNkxgtZRlSonaLWCj
zsB+BpPtg5FjRI5SvRGSfabhewXr+l5GDdKv7L84Qh4acZTKyIhOqGQiZBA60QMKxR2PBA6tknjL8wb0
VWglRVxQaK4lBFsmZYXTSAeo3TbKTDP1La7ZCRqGTNjklhpNwBloJaT1QpMHmgTqQTdmRApLHUWRuKPh
BxylVDnBCGwFDMdcFN93czKCisH+EoKGLI2LlAjLNd0S4oMqBCg1yH2NMqiShtVZw9qsYdOsYZ3U71Vo
pMqc1JiTOnNSZw6k4yqpai6pmUvq5gKpxvugxltWecsab1nnLeu8ZZ03G1RZAaXKCkhVVkDKWQHtAt/P
YWfsie4Y37M4jRFP46VLj7uUDIEe4ZSTTSGyLy2uHNmHRgxwZpKGY5xgwvTDHnGhQyGSwRBeQXJ/nZgD
aw3jz1nDwI5vtKZvrRHkE5W/pUuTo/g2Qzbs7oIGf4nU5kk7AQmXeR0nBKzKYjqdkpxJjFn0aiHUoBEO
Q2ly8ZdJmmGl7oQMXy0sySbskXMpJmQjIL/JlDYJhoqel/5XSDSr6Ny3ssJMP0DcdMPsD7FWiUj1wkS3
fg3DvJJ0UZTNNUGoGceGjuga7IPCoqXxKOVhIhjX3aZacYY1DsV6lDCoZl8m2oa+m45Gs6kpqragpVC7
QudOegs9iqmnQM+gDVpcpTpJta0Yc43Jp13dsqUEQFCUPShNMbg7MM/KwNpNC5HN+O6argpjBjmnJJWw
O3+VIk3q4NJwtir2eddKVHjngFar9T3C8Wfu4Zh50C8ednsnXQxv+LPgntFXgNNiU8kBOacUbbROAt8H
L6kuvoO/FtolIvZH9hHY+6YbVdoPwYmRgNVcp9Dy+W6730BbpDHjVN7kyaC70XHUusBJAutmnQktyzVd
w2otxOhi6mxIlUex0l4/QI8IqNMz6F1B7f7JYHh01KPoqQQbVGDL8IAeDsPjHeyONnA7WvV6w2V/VYFV
ub07pOHw5OAwg9G0mRs5Pjg4CpfLHYxAWEkc1ZBh2B/Q5XKQIXHicSFhxzVZTI6XsFT4ZIdVsI+asYfh
YHA83HqnhK0adXQw7B2F/R5gW9cUcJK4bgiax7wzm0mxYhGt9cTT0QX8KIMsBp7Nfma7vmqG9SZAfvZ2
LSIYQn93wQlsDAH903aD7qVZlhnZI2CkVBpTA5iJiJEHyJDwznU+bjerpmWSqdST1YoS2Nq2HSyMGCGM
E5bgKCiRzd6Ut4xQYwAlg2w7mI0BOwKsKEJHxOUhBQljp2KjUy4wx2saOvVHkqudWA9hyQOQEDAcB/Yh
sTBfOVU8CVy3m3EwzrdbpikkI6Bayc7FVsjWe1U3OxX2raQdtGlm27oXJuSA/4D/M99VnA8tYCTS8A5r
sglmqb6gkPqJKQzNwJU9IRtGrjAsqc3aeazXJsGKbpELvFYlQD4LkvyP7W80o22KYTCWFNxSYuaG7SjY
B3l7YutSMyJXFWBQAymO9yqch5Cd9qP50y5tx0K0FbdlVrby8lGLGNMO2J+uupjRetBcU22O/4JP+Rl+
gK3fH2YjpbrlBHcsX7tKhpp3MHBoFXN3Z9GsSQFQenlOsQ+z8d8Ep9Pt7UTNDQ03Ey9BBlvIH7b8wXZc
sXUqbRgWDG0YzXNB1mI7XPa2zRP3xTH3lo1VTwAOU6FmWBPcRWsfzU0E9OOg39u8eTDdQF6x31r34TWY
OdJOdddxZmWmxMVEzldxsfsBHOGyyC5+XBjP2ZoX8/GCxRRKW4Bmi/67iy15LFJujxHm5X0CvRwt8ytE
B2hq/nGoHWfweZ7d1JRvM35/B8D3v5gc5NZiR59hOH4YvYxSBZ3+xExf8bIhKmvZW42R8nyA7wuextIo
CDM+SJfghGkycoeWAK1AhTx0Sh3lrjiXyNvilrvmQnCmhWm9C4cPGI6hkk1DCL9zxsMphzYP+JUbvLfF
ze7onbeuM6mJcT5wWalAylE4rnYqeddbbWGyGfmVm4PlF2NuzJw2TTXZhULnnAfBL1hBK9mBKbDH0b8K
Sfb77/wl4/4Sw0HCu78tVpQHOAfbwI8i5MFB+055ZMW9pRBaQUeYFKA+NN0+ACwvA2LgWeTdIs8dH9Cb
x3JafIIRmUVTU/TYYePVfKbz8dNzMpWNSuRR9Ob/Xie4IS8/KxiqNQ4L3rWAsanT59s6DS0lKzQJxFpV
rJ9mJVW5oPpUEx/yjPnbrc/I+Jgi1LCC7j8w5Wb8x/v5YnJ9+uZxdxp7akROLn+dXk5uRu8Xv90s/ppN
Tt0Z8kXs2WgxOn1smyOVgjMVbA9633Vzu0z4t32/HTy28xuMdtB+81i7CHlqv23nNwJlRH6xYBD2gqI8
bG84ntpPVZNiEZoOodfrHfZ61VZE3HHzyUBC0FZG1rYM10f8jYgprMTAM3b4mXlk9cKS/Fxz3v+anzJX
5JY1u2PPqPNxB06avV6nHrtEwol6A5sqevArV+lfGsilhETT2rg5pyPvMwKL6x8Cntrohx8QvYfk06vN
JOYrAOQF2H8cWKz2svgZ+TpOqobU+MW3jbhaACm1eZEX2YCHEcTvV86Gtd0uXfAVPHZ76Ojb9xAcJGPM
w8rSryBQNt5ODatbQ3jYucGLcdUqTczqlLc7XXn2pLCP+8vBV3cy1uinnyZX5yY+jHLqQblc7e/Jnlez
xfTqcn7a9owyXghdNpWnULKMasgRoYihjJIVntNy4WnAWa+7LJ8fQJ6qawa5++q8bh/6/8XV2RWsGI3F
LUUfs7vDRH1Ewn6aMh8hhTkFmo+ry3SNmIKKdU/DqisNMw/lGW7N9CZd2ks1U5IL13HQS3HtM6VSqvyD
45Mam60KtZHseiBHSGo+15XXPb+SrWpnskS2ym2XMTjqvHks3/4+dape+4L4aExOfqqkReQqp9DJhOgf
NSAcNjyT5E/buQXtPShoVjSW+hRHd3D23APaCKXtpffH/OnjHuStiFIQ699i6YN2maJd6LE/BTabFQiN
LNYRC6mM8FL52/vwJlzN2eiHn6uJKWfQNT1lF97Ki5vdYr9ybUu36//lpX3VotkFC0O3Q751YRPUPxp0
+0fdIfwMjvuDd/aHn4ZJ8wSKOovRr/PT7NtIUGoaO3vnjGbTm98nf53W/N08A0pkY6Q1EPcySKQgfuAb
B2TPUuwFE5scczikaH+lMuK+aVmYZbp42xWpx2v+WaUYrg2fOmq/eFH63GEQ9WOuHSt9/DEfrvIDoPnd
DMvCVdz8+sf+Nsv2Fmd7b1OSN+Vre1wu3HJME5CvBRERnHtJ8Rx3LkU8E9J8eusX27yFyKiH794dvCuO
jFkop9AC9Htd+8fvH9a8MqfRCnSkkoIpRRe19/go09lt4TOaUGgmrsAl7RKq/Zwrt46pOx/tcUCT8c2G
z90tSEnlRnn/DgAA//9prxJE8CQAAA==
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
