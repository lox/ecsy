package core

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
		size:    5033,
		modtime: 1448503017,
		compressed: `
H4sIAAAJbogA/6xXT28bKRS/+1OgOVZp4mS1q+7cXDtpvcpmrdibHqocMMYZlBkYAdMqXeW774MZ28CA
x4ritqoNj/d7f38P/hshlE2+LVe0qkus6Y2QFdYPVComeJaj7Gp8Of44/hP+ZmdGdoElrqgGAdg1p2Ft
hdXzDa5Y+bJfM6svNTUalloy/mRP2/UZVUSyWncAq4KirT2MOGhGYos0LGnQiRhHjaKZPfl6dgCb0S3j
rNPwFkC2oVyzLaNyh3c9XSJfNdKijw5i07JR4P+7Inc6o5BTwTVmnMo7CM8bUclOh0HAWmNSmG8GXlH5
g5EU5kJIHcO8a6o1xOAYZg1HkQmjh4/5BpUCb9Aal5iT1iJRU+7q2uKmNLCfxmHwl5Q0kumXL1I0dcww
qOU8v55e5bknmufzTdLaCQShlUVPRhhthWxTc/sZ/scaFVghTAhVahc2N2tzrrTxRflBfKjJcQsfFtOj
dpkogszenJ1HqPXeA1s2a061igFORVXhGS1ZxTTd3DKlk4idFuMkkRT4YB8Gxju8UYeZ/dPounEQoQ4x
ee5A90Y84LKxVlCiPoIjFZV5br5H6w6iegvF8bmrjZii3QIs3fA8/0swY/n3/Sqs790zH3cH9gqt6/zi
whNBjlJH9ReqJ1oHyrv90NCzvsjsbmk7Nth59H6/+gezfMiwe7ptc+o2qK9yFAPbfXs9gc4ioe5gHfmo
vvgkSOtz5H19+yq7p0o00vTWvs6gHN9KA4fCX0ggHakZVb5hVizsij47tC2BthLyYIgcAuNkzrQ+9HXu
pW/ns6GF/bKT/8yzdM6fJNBNUH1+OWTzGtzQgojSKNakDqonu5GiCig8NOdYJQXaVuL9dE3ZRs5N6rLx
uf1zMXYPvA4Xbooo/BIosdKMHGRhUuZ5tHWTJdGnVtfn3e5wSgeSuWuxsLpHiSAmjoczMh5Ux0YzEChv
L3RHzHNj9n5lsBud76fR7Yivq9ViOABfKS51MS0oeQ5zvMLyifYNS00eu9ebBZHxYQwLqd58Qv9PjIH1
LTKCLo4Pn8dUaXUheVkVwEGFKA2VXXkC//KiL3I59mTmHObED2wy8Zu/s2IVFY2J6+/RtgFXOSWGfmcS
XIaeXYiSkZcwPdccr0u6seQnYcQkQP4YH1BSbLLsriPpWbI0s6QVGqSM/lC1y8kp6rsP44dJupmKhhvz
L91mdbpwoGFT74bQGl/uOGkHz4JP47M0SQzAhgyevMPEm/ZelKF2jwJtoqxQNMTJV6SrKBA6uYgC2/xC
mk/+znMrMVhFE6Waympry38mCPzmIR+Zi7em3cb349fH6+0WGssaU5biZ4+CwBYGQ7K2Xdtno+zQJn1W
Q/aSf44r/Etw/FOdE1H1eeoxWOnxVjYhXU4ixKm0yg9hGbhbRy/CbhUssC5MLNznQGZj3ebiSHe1Gelq
PHMfNOE8OpY6K5FOXx+3OzKURit0JI52n7Z3I/MaX+/vRu3Vd00/RHWmT0n6ZO4Rcv8WNlfPgSdSWt99
qO0b08VJ2sjVCS6A0KTRhZDsF43euiPnHqMx3j1PTCo+RIbyyQWavPCOzL/X0f8BAAD//w1ilgapEwAA
`,
	},

	"/templates/build/ecs-stack.json": {
		local:   "templates/build/ecs-stack.json",
		size:    11931,
		modtime: 1448503017,
		compressed: `
H4sIAAAJbogA/9RabXPbNhL+7l+B0930w41lSfRbq5m7G0Wxc7rWicZS3LmeMx0IhCxMSIADgraVjv77
LUBSIkiQop04Sd26loDdxYN9wy7QPw4Q6ox+nc1pGAVY0UshQ6xuqIyZ4J0h6nj9Qb/b/wn+7Rxq2imW
OKQKCGBWc8PYz3T9Fga3AzD0msZEskhlUuYriuJ4hT7SdYSZRElMfaQEwoTQOEYKpimJEeOxwhyGzFpG
0HwdacEa43B4MfaGQ1htCjLMB7OsId2kHJ1RolZCsk/Ufx8DyPcyaID1znzCAeqiEUeJDDSmiEomfEZw
EKyRLx54ILBv0OOt7N9hIzFaShFWkc6UZPxuN/6aLnESKD1lQ51ku80Ym3SngASJ5VZBGifoEC2FNMpr
UFwDHOUdhYxIsZscBYF4oP4NDhKqDfy/bAKmyMlRgOUd3RKnY4+uQc85euIc/bE6WsGVjsUhmKRMR32W
hMXB8Ng9WFkFxqpLw2AOPhv8YNnsCj/OwAH2mCvEjyxMQsSTcEFl0XCxtlyAE05WVVO9NeQuU51ZIGBF
Jqk/xhEmTK33gPFTakQycoSXEL7PAHFsa4LxNppgoIMvqYmBrQlBPlL572Shg52XclB9BFgw/ysSk5AM
PxLcBFQqGIHkFoG8RXERYhY8FwLVzAj7vtQp8bNwTHEcPwjpPxdKlPE3oXgrLshKmDQiE9oC3cV4Ng6S
GHzvGbC0L6X2WRpEIAyRTFq6ykG2UuddoqJEFY6nmcLkYznJmhSnJcO504U0GlI4UfTnWFO3QZ5L+GOX
Pa7pUoss0GdzG0vejJJEQiS+kSKJ2om0WZxSbyLSTpYmdONKFpwWNOeWcsmHw/8Ixq3jAcYPC3kUWVOo
wF7allmx37GmN4dP4B18Bq9X4i18+3BQ/rQpudkVjiLw2IKfQYlyTe/AYedidDUpajGJuxTHqjuwNQlU
k9caDg5Z1/fJ+eJscb41zWGB+4ECt9fAjY9/8pcDclzlpknK3bT2crA4OVv86FW5cdTlQkLE7YN/3D8/
P/VOqVNELJJMRNMevOPFCV2c9Ds1Gr+mIEcSWlA5FIR5FTWVYskC5wFgasfJ6Ar+UyLegoQBqPsUo7EN
cIrVSovoFYuEaxGUqqOyp+V+Botq4qKfbeo8K68KM5bGfRiKveBHcZyEVNNORcDIGtI4fOfKokoTpKLZ
xJ64vVguKTHJ3ZSKncMywRTSOGERDkqL5EtRec8IrSyUTVPiHeEQf4Ki/CE+IlBeV6g+lEY2FQgjkp8a
sYqHOyW0ifaCNLftrzDHd9RPFTqSvOwHHSz5EMAPGQ6H5kNkSHtxuvOuBCS9kdkjeO9YcIUZpzJTDJxE
MGqjrfUTA8K2usNXUqgtXN3QZY3ctjkpcW/pvoYvbQ1Z9ZUOCUTiP2BFVsNpoq4olA/kNVa4IiWnXZrO
FuQN09JiQU1dkCcVFx9445Z4jiHVl0g+VBBvpcGm/v5Uf3vRtAI1CbTGYgZtLZxalcrDcpsCofUl5drr
SDfT8W+C04kPXsCWzFRN+/fkKAWs0G7iGjyLy6vRX8Ekv5j+BIJ0ye4SadyndNNRFOug7jgtXe2dimLy
WTdrpQG1WLNZJ2t921gUUaZyijLR0GjVn6leoGMdVNmUVXrn7WCZSDsXvoNcMlKpWredRrPTb71zLKmx
QJbASoVwGqYzdscrJ1VnzkIKBYtecDofnF5Z0DpjkfBdI1pEsVv7feQDctfKhWgC1eg/KW0ZA/hAnoHj
Cd+dmgMLC5j7lU6BmUfYk1MMTa3eS7qR8j5+xUy947Ym4qqSN85s4vL0VtnExdimjhGEaY0mC1DpJBql
TXKl+US7c+tKcKaE6SkdVFY71S7nOjuwoobszDEJoUyY+GW76ubpknF/wqF/qJxs5RaiXTNjtJyyNTZD
WbW99xyquZcsLmmRuIXgsL4yL4pylPBOgdVL5qIU6y64zKovh0xt4DDGKxzTs5NKsVrX5Zq5Sq3gqE/+
+pfegvHeAscr1H28p7d8nYTp5VcQoO4aQWnYJUveXQihYiVxdMt7IlI9GDeceo6BD6PuPeqmtxLIUaRU
a+ySY5gqx9aNQ0dbVlhLZjkBOYLVTGtPexYWl5PWAXGoIzZJCnUp+tu/vpVOHHXU19RJpQQ9qPu2aTif
oFzGfjkgUjhjXStfbmtl6JhZpcYnxiOqMaODN3b0fZ0eVaQHp73+ParhziWralOxnW4IyozCYQH946ZO
y+Lfx7+8n80vrv9Rw+u2YcZfe+NX/nGYMxNRNWr6U2503WMuNwmFb47ffr9/1u+7OhvxwE1Z3pGQfVwE
d1mDkBIc7FkzNTCRgh+tIEyCda/0VPadWNtKyjS5heVjSLCf0G2dqDbGdzw5Pt0JOuiHHxB9hHzfv+VE
P0VCgoN44oBv+V3g+yfqqTAqW/aWh/fOCdRbiZBCxHtd/bjSO4rjVZWZrMBZEXhpW3LwXJQTDdsxvWB4
nb9IeDUWcUSEIea+M88uIQxX3Z0CjFZqg8/IMRdc+6K3GWDzoWP1DbUvH9Vnfpt0b49gyEqPRjfTcQ9k
oVwUsmUh8xRSrdGrzyR2OWkhm/C7rBVpaCEmEcBWgojANCMkKne7l1KEUyF1KvS80txc1M2MmS8nxo36
R+afXv8pV0I1+2hrk5yjnWlq1VzTV1m9RIP6iqrrF28ocq2dnZ4enxatl/a71k6eis2pztL7W1V74I77
laVt+ioQ5KPmHeR2HZy5etw55fDH3LP42avrc+5oTCFcc0fT2HPWlNHuhOD2wTfQ1T/gdbPiJnBOS05V
TtysROcCIwVYV6ULa7eNqhx7jVYCWOtP+QacXt4+F7mjObs/bQ5gQ9TixuUeswAvWADOr+9yXZ3zjAbp
9b1dbvX3XFpo1jdUjX6rOZf290b1x4/7PqMaUgMIKe/keQdBg/IH31z5gz+B8r2XUb73zZXv/QmUf/yl
lX8tEtXwAml0b2jmeNHmzfpLwNn9L0j7QO3HAyWlYtxchlja3FVcBV1+Vv7faamWP1N2u5OgjWVS0t3K
+U17q2v5lLm+eCo/6L3gZgffx2YHX2Wz3vexWe+lNnugfzcH/w8AAP//SsIMBJsuAAA=
`,
	},

	"/templates/build/vpc.json": {
		local:   "templates/build/vpc.json",
		size:    3348,
		modtime: 1448503017,
		compressed: `
H4sIAAAJbogA/9SW32/TMBDH3/tXRH4OEBuERN5Cx6aCBFMbFWmIBze9DWupHcUOqEL537GTtHXcpCla
ux/aqtT2fe27z93F/TvyPBR9n8WwylKq4FLkK6rmkEsmOAo9RAIcvAo+6H/kG9tvhcoKJfWSkeqJeZZs
B2ZI0wKsCT01hVuzkzFsJsvqWfr1DrNiwcHasnuXSx6GnwUzTv3Yzup5H/nW0F7yLHnblfrEALWWS/8/
tPgBWuJordHPkfutQTVqztA7SVHkCfRnIF5nBp3Jahh+GpMwnF+Pt5DQdS4yyBUD2eY7Zsv8YyqSe6PF
wevq7w1+b+FFEy4V5QnEwPVjbSyXcEuLVNlWMb2TTpbaQNAXqLRf6QpQm1xH5tsIq6hmiib3lbqXZNlH
sim6K13sf+j6MLgJV5DrjG2MD0PsPCBS2tdfK+BqMEf7isGkOQ5Olg65DbVNANslq2BN/fQKrabdMe1q
3+BweLXRcEDRb8pSumApU+sbwd1KqN4CM0ghUU6FeV4w0INGegUqupEd5eUU2BTuzPvPMSqP6Fsb7H5L
Yd1S5B06NXz85PDxC4BPzgOfPDl88gLgvz01/KkoVIthB/vKJqaLFIb5n8Kdi+YyHHZq2J8LkIpxqnQ2
WjQ3N3Ngs3zQ+39HqVffwD7uJjgmM7Xp7uRISpGwKtphNLW419m9X3dnDBY/j2DxowRLnkew5FzBjsyn
HP0LAAD//1i9GBwUDQAA
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
