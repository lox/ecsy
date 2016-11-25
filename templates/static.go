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
		size:    4522,
		modtime: 1480049273,
		compressed: `
H4sIAAAJbogA/+RXS1MbORC++1d0XDnEqQQIyW5t6WZssrjKy7qwIQeKg6yRGRUz0qykCUV+/bbmKc2M
wWT3luSApx9ft/ql1mhFNU255dqQEcCGmoevNBXJk/vC76eME1hbLeR9QZhzw7TIrFCSwCbmsCukQSII
qB1YJFkEASEhN3xUYc75TkhRaB2IKyIurdgJrmvY89m6gwVW1UaQOUtyg+f4bwYqEA95pqSlQnJ9iUc8
FJzVSg6IWktZ7H45K4br74KF0CulrQ99madbroehM5QFd/auGZVxWansaJ5YAn+cFKFZnv0sfqJoBFua
UMlesDFbrznLtbBPf2qVZ76x6bc1IeezU0JCEbKI+h5MMT6lENw7KdgpXSZneYZ/qYWYGqCMcWPqiPp5
W0hjnbfGeXWTsWFHblazYfMuAMhsrNYeQ+Gyw1znW8mt8XFnKk3pnCciFZZHS2FsH7hScy4zzanlzaGE
dLAXnCY2nsWcPVzr5NAyu75aOsRYWHiMuYRIoSjEBRZzWAZ2WqW1rTBzx8Vx1ssZ164PGDq1iA6xPJWd
5qHS4QBrgar+8bMXGh+PRyMs/6jo4yqa14ZfWJu5AHJZNzLAm/N/cpoYuIW377ruTj4gEtyNfH3TA7hU
FrUPwHFIf+c2y8sEr7FvH4pYFEg3NMkxLGPOzEc8WIoWiPtddfS4aoQlNs1Z1TOB4pvFDm4rnwDGXXfH
H1pe7Bjk+Pjtu4vNZhVA3s4v124U3U3I23dVb0+6uo3qHs0vXz6PK427ofFZeYwGGsak6fPitH25ijEZ
DdwjjVzLQLkrblSusVlJOaZeO0IKiZXGoYSp5KbOd9msXskCaRu5nSpl87kGkdZVMR6hBsDBsYjA+Yu/
JjU1sL2Q9xpnUGMTPsIiQ1+sYiohYFnWcAC+YhMWMxjalHnsjXqGORORXmQETo6K/8cnr7SIqR4wFVL7
NtxI6pRPPyEJNVawVganBSG+SqHRtDnpdvie/NUj1vVp+bOOR5ABFLgtQxaQsZWrcvSJWOYlRG07yJzv
8zOpqC+XWiS4v33BNisuiBXDG/Gt7Q3V99yWcqSLiC0cXAutiZL+tImxBmOVRAROG961jHvcT23VLCQ2
83eKvn1uiRuRcpWjH79VJHREcuayNtfoEGZ2pRLBnlrPzyXdJhyxrc55H+j3poYGBuL/UUTm16yiYUNh
O/+khbXH6e0FQ5fmL1zY3YswuKfW7p4qt/zh2qzv2s4FWzLx4hKaRzOVSzT3qa44L+dB1YXvEz/bjhDc
JMFz47m68I2VsMXuMrSzwMCKMqmpAbFpniuVdHYGR6ntd56L1c7QUjp7SAHWTcFi+hchjrMn/lNj8rRQ
LZM/Vwy/pW3jiquf5SHJBft8t8PiQStJoh49jjMicHZlNCEB2U0br0z8fx8Bl8cjmtIfStJHc8RUGshM
Wfte9rUMDrX2ABV7RW3sVvrqy53K34hQrTxpmU1va/WnwZ5Y7I3HSzHZf4rq/OXsd6/MbTP7y71ty9+/
RkXze1eRunn+uf2nd4kcBHbVhfombPwyFDt9yXOUmOY2Vlr84EOrZE+rXo/xyfF+PPo3AAD//2GxL/6q
EQAA
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    8679,
		modtime: 1480048947,
		compressed: `
H4sIAAAJbogA/8xZe3PbNhL/358CVTpR3QlFSVZsh1PlTrWVVNfY1lhKO7k040AgJGFCEjwA9CO5fPdb
ACTFp+MkN3PXTmQS+GHf2F2Ae48eoRNBsaISTUL8kUdoGmCpGEEnPFKYRVSgBRXXjFBEgkQqKvSarVKx
57o+J7KHb+CfWdsjPHQnfy5OAp74L7gIsWI8cgNNXrmvJRUvE+ZT918JIx8EXTsUlm9VGOzNscAhBeLS
20NokawiqswjQqdUEsFiTclDEwQsQowkjWGFoj4KmFSIr5G0a5DiKE4UUluKpidDxCKpcERAPRYpbggu
72LqgXpA55QGLGRA5hVQ2dOcKUkEU3cvBU/iBv5LICtTDNpo0IMZLpRg0UYz+Z3enYO2beTlFn2gdzFm
AiUSNAQGmABBaXiAyXY8CuTB7J4HAngekJ/DYvOg+WiWk0RtuWAfqa+9IF+LoIH7hfmLA+SgSYQSERjd
qGDcZwQHwR3y+U0UcOwbIXFO8wrklWgteFjX13JZ4yRQHup0tDCzVHyDazaCgint1UxTLQkYA625MFZo
skATQzXshYwI64dJEPAb6v+Bg4SmwYVAV8BEOOLF992adECGoH8JQX2WhMWRAIsNzQfCgyoERmqQ29rI
sDo0qq4a1VaNmlaN6kODfmWMVImTGnFSJ07qxGHouDpUVZfU1CV1dWGoRvugRltUaYsabVGnLeq0RZ02
G1ZJwUiVFAxVScFQRgrGzvDtAnZGS3SH+JaFSYiiJFxBfi3EuclfAU4isi1E9rnBlSP7ULMBykxQ/wTH
mEBGamHnWxQiKQzhNWTah7E5MNqw6D5tGOjxndoMjDacfKDit2Slc1SUZ8iG3V2Q4A1PTJ40CxDUL50e
LCUEpMpsut0Sn2mIWfBgJlSjEfZ9oXPx13GaYylvuPAfzCxOF7TwOedTsuWQ30RCmxhPT6AIm2r9EI7a
i9Z867SOLfJiD+hXfCNjnqiljm71EIJZJemhIF2rg1CxyDQEiG5APygsSmiL0siPORTLXlOtOMUK+3wz
iRlUs69jbULfLkeT+UwXVVusoXb51pz0mkbKVuoM2iDFRaKgwtvORGHyYVe3TCkBEBRlB0pTCOb29LPU
sE6TI9IVP/60m9hvaTxyaGluP/WIedm1EUW62SxA9x4hHH6MHBwyZ9gfHPb6z3rYsd2a7r4cDuYKdQ0H
5ILSexs78wj0s4bOB/MFHPy4MV2d3ehXJOsbr7I0YFu8MxzH4DFjRmhWLukG/LTkk7OZ1SCRDoXO0xl4
6BOC0dmph7TYg2fD0dFRn6LPJdiwAlv5B/Rw5B/vYDe0gdrRut8frQbrCqxK7ekh9UfPDg5TGE2aqZHj
g4Mjf7XawQgElMBBDen7gyFdrYYpEsdOxAXstSaNyfEKXIWf7bASdlAz9tAfDo9HuXVK2KpSRwej/pE/
6AN275ICThDbB0HbmPVkc8HXLChlXtNZziZn8FMGGQw8653Mdh3VHKuth9z07ZIHMIXeQmgCDf22j96Z
LtC+NXPSMy3kJ1ImIdWAOQ8YuYPMCO+RyubNJlW0PKQr9HS9pgS2tGkDCzOaCYsIi3HglYZRdvLR4lMy
TDeD3hawH0CLInRCbP6RkCh2Ijaa5AxHeEN9K/5ERHLH1kFYRB5w8BgOPfMQG5grrSiOAKr5Vhzmh7RU
UkhCMGo4WxMbJrn1qma2IrT50UyaFJO37IUFGeC/YP/UdhXjQ+unz5E3WJGtN0/UGYWUT3RBaAauswOn
ZwvCippsnUV6bRF4NEcu8UaWANkqSO4/d75TjY4ugp49ZpeI2WkzC/pB2p6aetSMyEQFGNQ+isNWgbMQ
Mst+1v93SpuxEG27TZkWq6x21OJFNwHm10y3hMwlHOEjbYRZdIrvYNsPRulMqWKhH/4BBR+97TrdJ+ht
gesT9MMlXaOu4Wi8p/Hdd+9S6eAYyxdwFIU60ixlAVB6uU/oP+Yn/+QRnfla9DWDSq1rrr1O2M/EN7UN
dtuabRJhosxq8uNPDVPZqqxvBlT6mM+k5wM9Yx+zmWpLD4jKUIbUMVv05Sd9sQDtNcj1JOsHdInPyvAT
ozveQBROlJXaNpBp7ShR0THxTVRMmIMNbHLYBYaNzgXbRMU0u2QhhXrlofly8PQsHz7hSWROBfrldQyt
GS3TK7gWJNV/LGpHGQyeJS05i/JEPtgB8O2vOrVYP+zG5xhOE1ouLVRBpj8xUxdRWRGZduB7jRFyf3Q2
LGiveJwwbYNkBUaYxRN7Bim0/7tLlTMeMcV1b1yaLvWQaT0u95V5PZuFUJ1mvo68cp/21gxAcTZD++/e
mt7i3X5FAqswZJbCe47BYbXTMD1rtf/I8NlVGYDSx2xGnxF1LYBk8iuW0Pqhf+euevSDu2KRu8LQ6Tu3
17vUfwfHVBPIQYAcOAffSIesI2fFuZLQtsU50IW+2IVpQ0dDGFgVOdfIsb09Sk2R56h9mBFpbDTFgpnW
ZkNlI7ZzlCbCkEPRXz/+7WF8GzLkvXyhomJ97MkDTQPKl7e66WOFQk6MUsUap90ly0XPpYq4kDT0v159
RUpHlwqv4LXsP9Di6uTV68Vyejkun5OakNPzl7Pz6dXk9fK3q+Wb+XRsj3ZfxJ5OlpPxp44+70g48LDI
p7c9u7bHuHs9cDvep052sdDxOpCHq/cT+50nneygXkZk532NMPcG5Wlz8bDf+fy5ImbIfV3A+/3+Yb9f
7RT4TaRrk4BgrcyYi+iGGXfLQwpOGDpaDzdVj6zbvfG8Zrf/NxOlVsiUarZEy6w1bxdOgP1+tx6xRMBJ
dwtbKbhzK5fbXxG+peRDk9q8Pjoj5yMCZeu38vsd9PgxoreQbPq1lURfyUMigA0XAYl1K4nnyFVhXNWh
Ri+8bsTVwkbK7RdpkS0YF0HUfuNqcGvuNe8baOx2ztH37xzzjSnyK15fQ4xsnZ0YRraGyDBrvS+G1F5p
YVqVnN2RxzHtexv1e+Oubl+s0C+/TC9e6NDQcsk7afOy25IuL+bL2cX5YtxxtByOL9g1FWOoT1oqZAeh
YqF0JC0y43KRacAZg4/L11RVd0GyvnhR1w/9fXlxegHOoiG/puh9eocXy/eIm09EW/19SJ/KoPihVbJB
TEJ1uqV+1YqamIOyvLZhapuszBWXrr+FyzHogyLlMjjNU+keHD+rkclFqM2kx/UMIaj+bFZ2eXY1WpVO
JwgIZZsnItSFI0bpAna/WzXYw6KiMRu5iRQGkQmaQKvio79qQDgbODqhjzuZ3J0WFLQjCgs1xsENnAFb
QFsulblyfp89vW9BXvMgAbbuNRYuSJcK2oOW+INn0ldhoJHEJmA+FQFeSTe/jW7C1eyMHj+vZqKMQE+3
jD14K7s0vUP+skdL19r/O4c+yFXGTb5vd8P3ujNGg6Nhb3DUG8GvdzwYPjU/buLHzQso6i4nLxfj9HuE
V+oIu61rJvPZ1e/TN+OaqZtXQCVsjK+GwVYCseDE9VxtgPRZ8FYwMYkwg0M6dtcyHWxblgZXKouTe6Qe
pdmnjDxI/xMAAP//gwz5PuchAAA=
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    2285,
		modtime: 1480049273,
		compressed: `
H4sIAAAJbogA/7yV32/aMBDH3/NXXNuHalK3knT0IW8Z3VA2qUOAmNSKB2OuzCLYkeO0YtP+9znOb4JT
HqqKByd357vPfX3BzsUFjCQShUA4BDvyR3BYMKlSEsFEsufMc4/qRcit4/xMVZyqxHcAFjHNFv1AohR9
OJviU2Z0tHGWrjjmYaXfPAKcfReMwyOcX51f6cVsyqMHV803t/XmwRKWJjHSVDK1H0uRxt3yLbfjTDER
qaTY5p3vYx0f/Jr5/teR5/uLycjYJ1LEKBXDpIQdsbX8Egm69cEdfDK/a/e2cIY8UYRTnCPXy96HNT6R
NFKFe042VR6Aj/ADdcg92WFla6MbnpkidGuCdNRYK/9C9sehQ65QamWKoMMG4O+/RopA6by/d8iVVYFO
pEWSg7rhusBvchitK09rJAbHAXKnpWTwTFhEVizSB/sgeCbYDCOkSs9PNjRjVMFDAo+1jFPcMMGXS8sx
uvoYvc+vobpvjOr2oIKV1TuJ1XtjVq/BenkI++HSSnvzCu1UpKokOIQ1vjlZRWgBtuW7y7+7nqyWhHeY
KMaJ0j01Gik/9EER1T/pNXUVkDfZmPm+rvOQOkuQJIIyw2ShzndU1Yoap+O474Djno7jvQOOdxJO92rp
4LSul+MAxqcni0oWZ5w+6H/Xa729KgD19u5QF100C4V8IzFpXSZhrAsrQUXkg6Jx40r5JsVuIqTywfMa
5rk4YsyGPoybE38oQ6t0rxpFZJ8o9aF0ZARLR3U35YCXjdwOhzfDUi1zybdxjtf6HwAA//9iMYjb7QgA
AA==
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
