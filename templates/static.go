// Code generated by "esc -o templates/static.go -pkg templates templates/src"; DO NOT EDIT.

package templates

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
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
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
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
		_ = f.Close()
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

	"/templates/src/ecs-service.yml": {
		name:    "ecs-service.yml",
		local:   "templates/src/ecs-service.yml",
		size:    6040,
		modtime: 1634116747,
		compressed: `
H4sIAAAAAAAC/+xYT2/buBO961NMjQAFijpx3P5+2PKwgCO7WwPerBE56aHogaZGFlGJ1JJUgzTod19Q
km1Rfxy72V3socrFEGfeG3LejIYZDofe5GOwwjRLqMH3UqXU3KHSXAoCL8ejy9Fw9G44evfSm6Jmimem
WPnVAwCY+QEEqL5yhgQm259ARQgUVlR/gSlGXHDr43lLqmiKBpUmhfddxuZh+dM+q4fMonwMCJn5Y0Lu
lj4h83C37vCvYgQeojA84qhARnC39MFIULkALrwtwTJfJ5wF+VqguTzEVpocJoy40gayAhJ04QBcgImx
YNcZMhtOCPfcxOX2OgMZPzcQjUyK8EcimfmBn+TaoGpGEBjFxaaf0+aala72nKkxlMUFo67SbuSOI0CW
K24eflMyzw7s1THr3fLEbrgwhI21BBNTA4wKoIyh1kVsXGhDBUNdBmHV956mPHk4daNR4QWCpmh1ZXdo
rJS5gFzjv4W+L5xTGdyyMFXqXEybwR2dL4WhXKC6pimeysa2znVRyLouGiRLqUyT5DpP16j6STKpDMhS
3w6hzFDU3CKaJ4bAL6NKiIur57IlkoawpolV1hGMH5AmJvZjZF9uVXLqWd7eLCxFzA3cxygglFxsIC4w
mcXUECmZljldXLXDuCijCIKFj8pqgFGD7WZzKI6JaAiICosHbA9YiQciqfpDGQw8z5ciLORWNfxbjR+M
yRZcGxT1DvRi9mdOEw2f4MUNRq34X8NgAJ+9OobuALmWxiIcg2XR/shNlpsqtMBQ9qU4nR3eHU1yJDBA
poeRVCkqQuzvStaD/oZaeRb8+/Wd/ULS8KqSVNtpHsGn3cvtM2hue/C6ZfMiyNcwiK0Vubg4e/ywWi0d
rvPpdWBL/Ds5e6xK4/tBnB3MkygOyOfaZ6CcDfpOJqh3iK6eWrffr3veDWqZK4ZV9maLqx/54uyslkpm
ViRbvN1TmNXqAwgE7qfIFsFscWX7d6SkMLZkZn7gwhQfYKj2UX6NnXUnqrnYKNTNSABgCPNsqaSRTCYE
DMtaFgDvlUyLnledcZmgDsOVPMrM56GaZwRG58XfxeiZUb19++ZAMN2r7RjKXtsQZnfSE6oNZ3s7LjaE
1N12XrtuRZqN6kmdlDNbZ8q2KW8MokdajhuWjlD6+Zr10GvXmNQcu+3uO2nqJ3iUkObVeFYzdkaCDpe9
qmyqHYPad7Yd3YqqDVqSopFZX3L26JB9P3t0v9SNDraneFjFCnUsk5DAuGVzK+KW1WW7RObCoPpKEwJv
2osrnqLMDYH/OUu+FAKZleNUUS642CxlwtlDe7szQdcJhgSMyrEf/v+1sgn+2brRPwvnv1o4fUF0d95n
cwcdFq3ptHtQ+1nvf0+9d41hzlQU2KmonMOeKtftpNsabetGU9RcYejLXBgCl25p1ITXWR7uXbQhOfuu
az5wrpbHybQeSEVlB++uQbsEbLXN2vv6a/jssN3IpD3v2peeq2Dnwr+fd+v/PnOTWQB3JXQ++Z0Qh6Ev
mxOt87QAKqU2lSxPUZh2XgJDDXYvlWmbRREyQ2CSJPK+08aGwQXjGU1Ij0HZpxtS7XuGgEyf05R+k4Le
63Mm0wM+E+b+J6UfVRtN9gfjOCypie0d23lnT659byigylMt1VW7OHY1zSfO/6gsnJKLU08Fy3EgkTRc
78aB8l60xlc/DqBwY+tMbT802t4UOqeME6FvmsAfuYlPBWbj0/bIxmSSm1gq/g27rnQHMbZXWgKDVwPv
rwAAAP//5GQuIpgXAAA=
`,
	},

	"/templates/src/ecs-stack.yml": {
		name:    "ecs-stack.yml",
		local:   "templates/src/ecs-stack.yml",
		size:    11590,
		modtime: 1634116747,
		compressed: `
H4sIAAAAAAAC/9RafW/bOJP/359iqi3WwIPKsp00aYUne+d13F2jTWPEbou9vUVLU2ObqEQKJBXHzeW7
H0hKjmVLttN9wZ2DJDL547xxyBkO5ft+o/dpPMEkjYnGN0ImRH9EqZjgITS77U7bb7/226+bjUtUVLJU
u55Bfwz9OFMaZQgk00JREjM+B8aVJpyiAsIjIBw2kK1mozEikiSoUaqwAQDwMaXDyD2az2SVYgi9T+Mw
HPS7Yfhx1A/DYbTuL0kxWSCwCLlmM4YSxAw+jvqgBciMA+ONgsFIsluicZxNOerOPnYOsp/jjEmlIXU0
QdkRwDjoBYJKkRphIlgyvXDKVYvR/bNiKKSCR0+W4y2u3pMEwz2E1QK+4iolTEKmMDIGJZSiUpY0UvU4
ybUavMXViDBpHww/x7yX6YWQ7BtGHxRK9UHGNXJc2/8kBh96HDIZGyFSlExEjJI4XkEkljwWJLLikjXd
z19xpWAmRbIl2lhLxucb3GYki3UInudEG+YqWXS9cfQqReNohQWMXJlCmAlprVNnmTr2uttKGJVi3dOL
Y7HE6COJM1SPgpiPb9CccLHbWqaxblYJieMKNEYsS3bbYyLnuNWcnFTDk5M6+F1Ne7e647RV01wHr6Vz
WtfRaVf20GrWtIY1rWNN61jT09ar6o5q49Ea49E649GTOs4nNZxlNWdZw1nWcZZ1nGUdZ9Zt1bV36zpO
6zoKFrbnityN2bd9qzYhdyzJEuBZMnWB4jFMaQExyThdbK3Y9xa7u2LPHNNLVExi1CcpoUyv9jCPHBJo
DgUy0yifxvQk15TxQ5oynv1VmnZyTQX9ivLXbGr2bF6KHTV724ZMv4nMRhE7EISLTo4i/JpNd5k2m1tc
Bwlh8ZNZohkFJIqkiVvfw3dElFoKGT2ZdZoP3MP1vRjQhQhBywzrRRn0x0WK9QQZjBc4Y88se5N/UUfG
kX0n5ioVmZ6YFaSfQrqIyi2IcxrGuTXjxLQDziElKUotje2RR6lgXLfqY+4l0SQS817K3uLq+wSxC8yR
gd5oaBIXmypkagGRMzzeItcmZ9GigNbIdJ3pNNN51B1rQr+WswEblEPwkCp/JmSCMgzNszJQr37S8nHP
bnC20e/wY6SZZHr1ixRZWj2kBFnPof1WTuW2+BSgRqPxA5DkG/dJwvxuu3PWar9uEZ8k5JvgvtFApJol
Jn9q/ABjRFhonYZBEAmqWmSpWg7aoiIJevZx0B8H5sSgdBDhLcYiRTnPWISB22I+U8E1YRzl52IDai10
EjeuSJoyPs+N3Ps0vsE5E3wielfDR00y5SNR2u+EcA+9q+HwMgQjfOd19/T8vI3wsAPtbkGn0QmenUav
ytAlVlA9n7Xbp9POrAK6TfXlGUanr0/ONqCYVVOlr05OzqPptAylyLUk8Q46ijpdnE67G2iS+lxIvai0
BH01bXfOyOsyXomsBn8WdbuvTkuWK+G3FT0/OW2fR502PDQaN6hEJmmRjw763SJbHkkxYzFWnmSGvasw
3AKucSNpPEaz7Rx3RPQihKDUdiNiVCH87hx72LsyDfBHnre7r/USmN6DbHtKZQka6EjEjK4uBc0S5LqM
yrcFjdVdLjsZzGZIdegS+UqMEYNxylIShzUAtzPIW0bRKI60my9BsxipSOCPmoE96nZIpVX4qNRBE18R
TuYYOeV7kqtdwXwgkodkqUJGktA+pBYeKCeoL0WM6+2h2y82gFyPmZCDftdKU0ycZVeai+3JcwId9hoL
s5vh+iC3NfQR9jfPbT4DeybWBxqLLFoSTRfhKNNXqCWjJhgeHjSzFRrDwAXEKdo4VazQvQSQdtejJmSu
asAFrRC8f3l/syE8k0iEfYlEYw2zAmqRo0y/E/OBjeiH0YWy78R8rCWS5AiVCye3BP5lfryKzahidWxu
TutUoAjBlR5uEjH710XpQ05+gxq5MemQX5KVCqFzWuov5QROnKblYz3EtD7mlb1Mi7Gr2NXLtwEqfTlO
3I+j/n8JjsN1fa7WaBUVumOh3S3oO5t89AWfsXkm7ULZMEdFb3kXzI9XDp1/29on7zYR7lsJsX0udMit
1tIIsxSrbHMPJisGI/+LIrUzuVqRT72wpidzorGnnWruRLER4neoGe/8bmp2nTLB3T667Z1uCY3ZnFfF
tQlLUGQ6hNGk8/Jqp7svMm6PnEXDhzQiGqs4bXjijYjNP4fd5XnFeBEN1JCv42lnF0jufja7sZvb3f4R
yRQaDYz4FdJ/Ikxf87IJVH68a9T45eElt89d9yQygjJjuWwaMzpMe+4IHMKMxKrsqqVzxWOCVWreSjQK
c14JzrQwB7StQ6wFJWSOwyiEZ28Yj4b8iqTw+1aq/2Jzg3LtzRcu/6xh6azkdtmNpjKWJNuZaXEa6tZl
optFcQcuStabiA8KpYnQu07WfMPD8Gei8Oy0GcKzcTaF/6kMMz88C6aMB1OiFuDf3VbH6lWWuIpNHIO/
ArJUPp1xfyqEVlqStHJQIFIdkKWy9A2ccabBvwXfHU3h+X05EDyA78vcW6t803abOSlGuhl6OI67sv4P
PsLz/zhOhIqIdFCEK9Qk2pkRC+2bTOnNOlMKh5xVpHTUalydnxgPUfWpS4CaBkiV+W3to7PBywTvvf6x
+Rn0x5/77z6MJ4Obi+f3jzWDh6NGDt7/Mnw/+Nz7MPn18+S30eDCVUKePPayN+ld3HsLrVMVBgHjEd61
HK0WE8FtJ/DCe6+o7nmh9/x+p1j44L3wiopYGVEU2AzCluvK3bbu9+A97Fc5EZHJ2trt9lm77e2FiiVH
GYIUQu/FzW1CtB8XLESCAdKub7QPcqPQ2ZPc4KeDE/J/2/a5QQsrHGPUo7BuTpvt9mm73dy/BqkUvLUQ
mYxXwdYl4F+7IEubN2YH8RqVBv8beM/vdy89Hzz48UfAO6ahfZASzWRsdlEWI9fgz2pJ/gSBTtJtMxyk
n9xWjttxcaUWT6ZNF4mI4Kzd/ouoiSVfu1D4p2k+7h3n/9TeQUWSEB7tcc4ZarrwH9Wwuh10Zks1PLgm
GrVk8pzDf6xM+PYEfBzfpy+iw5NNNPz734PrN8avjVpqpVykDY6MZtejyfD6/fjC840qfiTZLcoLslRG
MXCNItWQt+TpxkU53ajA2Tl2Ubk4cj94h2Pr9ZvD9oH/nFxfXocgMRG3CF/yG4xUfQFhXzZYIMxEHIsl
43OYZnNgCmbsDqPwCOI+FNFkzvQim9pyvkncNi4CyBy5DphSGarg5NXrg2TXIh5E5mXCYoTEWJCo3iOL
K6b9mpl9Nvc/z+25HJrP78t3XA9N7+/y46OCQZApaUcUqmcc/Aj+++BAe3T3TSi/8Ap7eEePk6g0kfqC
xEuyUkcPWwil7QXil+Lpy9Fjb0WcJXgR3BIZyKxQuKUE/RraKLPRcCTRecwilDGZqmB963jcyB03gB9/
2g4aBcmWOWG1YjGvd8j87vC7/LF00/n/3B2/062sS0WR22H+CWdMoXPebXXOW6fdVid81em+tH+CLEqP
JYHQnPR+GV/k9+dh6SjWfAKV3mj4+e3gt4sdTziWxi1Ur6mKxieQTKWgQRgY0+bPUjxhOLWhsCCgViqY
qbzxeEL5sso18NcesrtWi8v78lKtuUDfeaWydIl+qIxmUaVXDgb98bruBCon5vK9ctnbvs+6rlZv1adL
Ugz53JbmKqq1w3QkhRZUxCFoWlX1eSNFMhJSh9DsVB2SJiLvPXv58uRlFaLPIjlMQ+i0W/Yn6JxVWHSM
8ewGZyiR051bVq/Gvrlm3sZrFinySF3zELwS0jtuKtYmrZ5GW/vbY7J6Y+0z09hVc0uKVcrxvwEAAP//
Myanp0YtAAA=
`,
	},

	"/templates/src/network-stack.yml": {
		name:    "network-stack.yml",
		local:   "templates/src/network-stack.yml",
		size:    4278,
		modtime: 1634116747,
		compressed: `
H4sIAAAAAAAC/8RXS2/yOBTd8ytu0UjZAM2jM1K9ywTaiWbaRoCo1GoWJrht1GBbsdMOqvjvo7xIHAKh
FPGZDbKPr8+5T+j3+x37cTIlSx5iSW5YtMRyRiIRMIpAM3VD7+vXff1a6wyJ8KOAy+xk5Exg5jlwT+Qn
i94RyE8GPJ6Hgd/LvkfBB5YERDynRIoeYD9iQqRn9pMYaJ3OHeY8oK8CdQAAJinQYfQleM12kjXzHARf
4LjDMQJDH6SfS+MPWG8g2UW9BjMH+qV5tQUzajCrGWbWYFfNMKsG+72AdR5iyWOZS5tx311UNOEwJggu
xuQl0bfZH/3HWSRLXLLu8TKBTuI5dH/7sh8nCE0k9t+T/XU/NdztVPyne2kMSiNq2HJQHqmB1shJsXQo
O22bnWJHq7I0DmFpHMLSOBFLo4Gl6WU5vJ+mWaT6Xp6FrR8TLQwpTK2DmFoHMbVOxdQqmY6JYHHkk6Ie
PKc0N11xgiA1MXJMhKoV4UWMk0gGxcViOcEi+jNk/juCi5uALlx6hzk8Kz2kB92Z53R70E0KtAv/KhZc
KiSmPpkSiqm/QrAgLzgOpQIaUTwPyZCKScxTH4CMYtIM+YsJSfGSiAbQFL/WBCSrD3+TFUpduXWmBqar
Ojgv+FssySde7fakSyWJKJE5sMmr8LVWjNlSYv9tSajcG6AtdCXvOKEL8UAR7Hm2FgqFprvIVdfvl520
bJ3p4cj16my1ku7I9bRWDkO2xAFF8MH9zOY9ljv8WzFdgtpfsMOQ+Tipx1RAXjkj1xtUT9aacinL541i
tTHv6/r1mGWwVpJ3mGeWXP5A/8Ex9d8a0tn+wEGI50EYyNUTo2kjICHxJTyD3oOLWyLtJwGaVqu5A6s2
l7SrchtTYMdY+WVuME7lBuMIN2zPrSP9cJZAm0cotM6h8GQxtL6lcMxiSURbJqeoaTJ7WjXue6XNjT99
ZphN1bxLptZaHmtv1UTIgKb9shKF4p+BrmAPGSmlxA2sGoL6ODiJht1v1n57fVdwyXNjuNxSJkb2XlsB
lURtIZgfpCxa5e0dW9/2fd4Iz07Y+Blh8+yEzabsaU21als9O2XrKMr/BwAA//+lORj/thAAAA==
`,
	},

	"/templates/src": {
		name:  "src",
		local: `templates/src`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"templates/src": {
		_escData["/templates/src/ecs-service.yml"],
		_escData["/templates/src/ecs-stack.yml"],
		_escData["/templates/src/network-stack.yml"],
	},
}
