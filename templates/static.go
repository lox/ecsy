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

	"/templates/src/ecs-service.yml": {
		local:   "templates/src/ecs-service.yml",
		size:    5086,
		modtime: 1485155260,
		compressed: `
H4sIAAAAAAAA/+xXTW/bPBK++1dMjQAFijpx3O5iy5ujuFsD3qwROemh6IGmRhZRidRLjhqkQf/7C+rD
Ei05H0XfW9uLw5l5Hg3n0cxoMpmM5p/DDWZ5ygk/apNxukVjpVYMXs+m59PJ9MNk+uH16BKtMDKnyrII
QgjRfJcCGcybn2+Bw4bbb3CJsVTS+QJXEXAFi9XF69FozQ3PkNBYNgK4zcUycj8ANve5A/ocMrYIZozd
rgPGllFp85g3CYKMUJGMJRrQMdyuAyANplAg1aiCXRfbVIqw2Cqk82Mclfk4TSyNJchLKLClM0gFlGDJ
aXMU7iEiuJOUVMkM0M9+ld6i0Cp6Of8iCIO0sISmyxuSkWo3zOSKKaoQd5OciIuk5LFVXYF0jRyiKIyk
+/8aXeRH8vJcBtObu+RKJ9g5L6CEEwiugAuB1pZPJJUlrgRaR+1U9ZFnMr1/blJx6Q2KZ+hU4rIhJ02p
oLD4T2K24n8uri9oqkviY7nK1CSBVsSlQnPFM3wuh2iCuiXW3Sp70GttqAt9VWRbNMPQuTYEulKmR6Nz
VHVIzIuUGPxnWgppdfGr+KnmEWx56pTxKMcn5CklQYLi241Jn3tLN9crB5pIgrsEFURaqh0kJZZwWBZi
o7OqRqsLn/jM8YbhKkDjqik4of/yH2OeqwMJcOVwQLRAdfkh1maYfDwejQKtolIstmK9sfiJKF9JS6ia
fgDwavFXwVMLX+DVNca9J34L4zF8HXURbA/iSpOLfw6Sw/p/QXlB5WOFxMW38j5KrFueFshgjMJOYm0y
NIy537Uox0MtrY4pOVtb7bnSPLqoBeK7L2P4Uj8/wPgwtfHbve1VWGxhnDgrOzs7efi02aw93NPLq9C9
fD/ZyUMt5p+D8fvwJ6Pr4K/7RlsN16F8w/Z9PexhXc/WNhpdo9WFEVhWYLG6eGkfLz3WRueuvGgbGUBp
7GgZGIR+a3eCXawuXI+MjVbk5L0IwgagHFtQP3E1w2qLx79UO4O25YUJLPO10aSFThmQyPcWgI9GZ2V/
qW+suuKOw0Y/ag5kZJY5g+lp+f9s+kLW9+/fDZD5p30O17MOhNIvTMotSdH6SLVjrBtSRuz7ADtsAUfq
WO0i3vU2BTlYpZ7wmO09vPL1kQ812LMf7Bq1vcnDA+zm/2hhl/VK0XHyBl7Hta2yK0pt6MyUln/DzQ4d
XPnSO2928uDB/jx58KfR/m1vIO83iUGb6DRiMNvbblTSs563YlwqQvOdpwzetYcbmaEuiMG/6qNAK4XC
ieHScKmk2q11KsV9m8BC8W2KEQMyBfaB/r2X5kBr/R3atH/E+TvEOUzmd51f5gg7lt6GMzz8/7wzh2Pc
m7Chm7DVHB+WfbPx9NacynyJVhqMAl0oYnDeSLAjAE+G/jfDQendWXc2eR8Bj8ukS1hDuzVraL2qgHqN
pHPePYavNcu1Tvu7jzscNVryPrna3ac99ItRAh4WZDn/H2N71H415tYWWRlaSeFSiyJDRe0dh8QJ/SN3
8Ys4RkEM5mmq7zoWRyKVkDlPmXcM4Imm+28CKOwpz/gPrfidPRU683zmov3q7EZZsqxNoDavOSXum6X+
y2XV3exgUp3dV1XtrOTdJnHkLo7ex1N3cjyLOv9qyLivwO1+yFT75xbfvCTE4M4p0zQt0br9rTetngV2
fQj1WVLyNJSYPfXkYsbmBSXayB84tA73opo1n8H4zXj0dwAAAP//QjnIT94TAAA=
`,
	},

	"/templates/src/ecs-stack.yml": {
		local:   "templates/src/ecs-stack.yml",
		size:    9456,
		modtime: 1485155260,
		compressed: `
H4sIAAAAAAAA/8xafW/bOJP/359i6i3WwKKybMdNWmGzd17H3TXaNEbsttjbW6Q0NbaJSqRAUnHcXL77
gaRk681t2j7A82yxiUT+OG8czgxH8TyvNfowX2CcRETjKyFjot+jVEzwADqDXr/n9V56vZed1gUqKlmi
3cxkPIdxlCqNMgCSaqEoiRhfA+NKE05RAeEhEA4FZLfTas2IJDFqlCpoAbxP6DQ0DwCLXYIBjD7Mg2Ay
HgTB+9k4CKahnSvxXmwQWIhcsxVDCWIF72dj0AJkyoHxliM7k+yWaJynS466f4yJmz7OZ8Wk0pA4WqAs
GhgHvUFQCVIjQghbpjdOlSbmg+9lrpAKHn4j99e4e0tiDI6QVBv4hLuEMAmpwtCYjVCKSlmiSNVhAxtl
fo27GWHSPhg+huUo1Rsh2WcM3ymU6p2MGrhf2d8kAg9GHFIZGdYJSiZCRkkU7SAUWx4JElohyZ7mzSfc
KVhJERcEmmvJ+DrjsiJppANot40w00x8i2s2gt4laNwm19RIkiqElZDWCk0WaGKoB92YUSns6CiKxBbD
9yRKUTnGAJ7BcMJF8f2wJhtQMYmiEgJDlsbFkYjINe4H4pMqJD6pQ+5qI4Pq0LC6alhbNWxaNawP9XuV
MVolTmvEaZ04rROnw+6L6lBVXVpTl9bVpSd12ic12rJKW9ZoyzptWact67TZoEqKDWqk2KBGig32pFoA
l+Ruzj4f8+6Y3LE4jYGn8dKFx0NI1gIiknK6KXj2W4sre/apYXOBikkMxyQhlOndEXahQwHNYEBWGuXj
2JxYbRj/kjaMpz+qTd9qI+gnlH+mSxOj+D5CNpzuggR/idTGSbsAhIu8jhL8mS7LbDqdEp9JTFj0aCZo
0EDCUJpY/G2cZkSprZDho5kl2YIjfN6KCd2IALRMsYnxZDzPU/8jOJpddOZbWWamHqBuuSH2RqxVIlK9
MN6tH0MwzyRdiLK1xgk148SMA64hIQlKLY1FkYeJYFx3m3LFBdEkFOtRwl7j7ttYW9d3y2E0m5qkahNa
qjYQOnPiLXJt8qkWObRBiqtUJ6m2GWOuCf10yFs2lQTQRqq8lZAxyiAwz8rA2k0bka14co2rwpxBzpGm
kundH1KkSR1cms52xT4fSokK7RzQarV+AhJ/5h6JmTfo9U+7vZdd4pGYfBbcM/KKRLPYZPLWTzBHhI3W
SeD7oaCqS7aq66BdKmJ/ZB8n47lvqlGl/RBvMRIJynXKQvTdcb+hgmvCOMqbPBh0NzqOWpckSRhfW2OO
Psyvcc0EX4jR5dTpkCoPidJeP4B7GF1OpxcBGLH7LwfDs7MewkMJNqjAluEJng7DFwfYFhuona16veGy
v6rAqtSen2I4fHlymsEwbaZGX5ycnIXL5QFGkWtJohoyDPsDXC4HGZIkHhdSbxo1pi+Wvf4peXnAKpEe
wZ6Gg8GL4d46JWxVqbOTYe8s7PfgodW6RiVSSV01NBkP8spsJsWKRViriaejyyCogCxmJo0HaHaoq2ZE
bwLws7drEaEK4G/nnNPRpRmAf2w16F6aeZmZIwxGSqUxGsBMRIzuLgRNY+Q6n7eHVWN5yGTqyWqFVAeu
HCzMGCaMU5aQKCgNm7MpbxlFowDSQXYczMGgIoZ/SuARdXFIaRUcRGw0yiXhZI2hE38kuTqw9YBIHpCt
ChiJA/uQWJivnCieFBHuD+NgnB+3TNKVkJPxwHJ2JrZM9tarmtmJcGwn7aQNM/vSvbAgB/wL7J/ZrmJ8
D2gk0nBLNN0Es1RfopaMmsTQDFzZG7Ih5BLDEm3Uzn29tgjpYI9ckLUqAfJVAbR/af+gGm2TDIOxRKKx
RMxN29lZqt+I9cTmpWZELuobsZ5riSQ+KnDuQnbZL+Zfu3QcC95WPJZZ2srTR81jTDlgf7rsAo1Oc43a
XP8Fn/ILslMB9IfZTClvOcYdS9fukhnNK5hRqsXc9SyaJSkASi9fEuz9bPw/guN0352omaGhM/E1yGAP
eWPT31jwFVun0rphQdGG2TwWZCW2w2Vv+zhxV5xzb9lc9QbgMJXRDGucu6jtPZiaCox8z/LiwVQDecZ+
Zs1H1kTjSDvRXcWZpZkSFeM530XFngcmuIsiB/9xbjxna16MxwsWo0h1ALNF//nlfngsUm6vEeblXRIS
jWV6Be+4FpH55VAHypeM59FNTfk+4vcPAHL3u4lBbi8O4zOSKjRyGaEKMn0gTF/xsiIqK9lbjZ7yZQc/
5jyNqVFQZmyQLiNGp8nIXVoCWJFI5a5TqigPybk0vE9uuWkuBWdamNK7cPkAmMZkjdMwgCevGA+n/JIk
8HelwHtWPOxuvPPMVSY1Ns4GLioVhnIUiauVSl71VkuYbEXecnOwvDHm5sxt02STgyt0XvEg+J0oPB12
AngyT5fwf4Ug+9MTf8m4vyRqA97dbTGj7NLY3X6jCLwdkK3y6Ip7SyG00pIkBagvEu2TrbK0DIhxpsG7
Bc9dH+DpfTksPoDnycybmrzHThur5iudjR++xFNZrwQP4el/PY5xQ1z+IuNL1CQsWNcCxiZPv9rn6WDK
WaFIoFarYv40O6nKCdVHTX2kyvzfra/I6Jgk1LCD7r/JeH4zfvNuvphcnz+9P9zGHhqRk7d/TN9Obkbv
Fn/eLP6aTc7dHfKr2IvRYnR+3zZXKhX4PuMh3nXd2i4T/m3fbwf37byD0Q7aT+9rjZCH9rN23hEoI/LG
gkHYBkV52nY4HtoPVZViEZoKodfrnfZ61VJEbDnKAKQQujKztmm4PuNvRIw+0oFn9PAz9ejqK1vyW814
/2l2ykyRa9ZsjiOzzsadXm/Y63Xqvkul4N2NSGW08yut9G915FJAwrQ2b+7p4H2G9tP7+oeAhzb8/DPg
HdPQq62kqYxMXGARcg3e6iiJ38DXcVJVpEYvvm3E1RxIqc1XadFNLEI47fW+c7XY8v3WBd9B43CGzn78
DFERx4SHla1foaYb7yCGla3BPeza4Kt+1SotzPKUd7hdefamcIz6152vbmSi4ddfJ1evjH8Y4dROuVjt
H4meV7PF9Ort/LztGWG8ULJblOdkq4xo4AZFoiEbyRLPeTnxNOCs1V2Uzy8gD9U9A5hcvarrB/+9uLq4
CkBiLG4RPma9w0R9BGE/TW0QVsLcAhlfwzJdA1OwYncYVk1piHmQR7g105t0aZtqJiUX2nFkjVz7TKkU
lX/y4mWNzF6E2kzWHsgREiNBwvK+5y3ZqnQmSmS73HYRg0Pn6X25+/vQqVrtG/yjMTj5qZIWkYuccvBC
+N8aEMDzTJA/b+catI+gJCpNpD4n0Zbs1BHQRihtm94f86ePR5C3IkpjPPdvifRlmgvaVYJ+Cmw0Kww0
klhHLEQZkaXy9/3wJlzN2PDzb9XAlBPompqyG4l1eXOzLvYj97bUXf83b+2jNs1uWBi6E/KjG5tA/2zQ
7Z91h4NuP3jRHzy3P/w0TJoXIHQWoz/m59m3kaBUNHaOrhnNpjevJ3+d1+zdvOIWmj2tYfAogUQK6ge+
MUD2LMVRMLXBMYernfJXKhs8tixzs0wWb78jdX/NP6sU3bXhU0ftDy9Knzug8Zpr50offybj+f6+CCoj
4TJu3v6xf82y7+Ls+zYlflO+ttflQpdjmsyk0IKKKABNi/e4V1LEMyF1AJ1+scxbiGz09Pnzk+fFmTEL
5TQJoN/r2n9+/7RmlTlGq2tcoUROS73x9hEbZTK7I3yBCfJQXfEA2iVU+0um3Bumbnw4YoAm5ZsVn7su
SEnkRn7/HwAA//9prxJE8CQAAA==
`,
	},

	"/templates/src/network-stack.yml": {
		local:   "templates/src/network-stack.yml",
		size:    3730,
		modtime: 1485155260,
		compressed: `
H4sIAAAAAAAA/7yWW0/jPBCG7/MrhuqTctOWHPhWwnfZtLDRLhC1VZFAe+GmBiJS24od2Ar1v6/ipDm0
TU8shRuw3xk/Mx6/bafT0Zz74YjMeIQluWLxDMsxiUXIKALdMkyjY1x2jEtd6xERxCGX2U7fHcLYd+GW
yHcWvyKQ7wx4MonCoJ39HYdvWBIQyYQSKdqAg5gJofacB9HVNe0Gcx7SZ4E0gKGSuYw+hc/p/5BmR/AB
rtcbIDCNrvo9N7/BQm1nAcaKxOoa59ZFTWKuSOx1ibUiuViX2CuS/5cS7S6RPJGqiDEPvGlOj6OEIDgb
kKe0ErXW/8NZLLN9gFs8SwXDZAKt/z6c+yFCQ4mD13R90VGpWlrRGcNXvc2C61eRC/Lud/W182sZtpPo
6yS1aL0kMncRmbuIzE8RmWtElp/NXDOStRzLRqZljiOhluEVKnsnlb2Tyv4clV1SDYhgSRyQbF59N0sz
mnOCQAX3XQuh5cT6MeMklmEmT3/ccBp/j1jwiuDsKqRTj95gDo+199uG1th3W21opQ+mBb/zWI8KiWlA
RoRiGswRTMkTTiKZb/cpnkSkR8Uw4apCkHFCVjd/MCEpnhFR2x7h5wIRoAM/yRyp1hRr9ba26o1SD+0a
S/KO55s74lFJYkpkLlrtDnwsKikcKXHwMiNUNrZ3TZlPByd0Ku4ogoaDilbWcLxpXlc1aulHpQlpAH3P
rzLpJVTf8/WG03pshkOK4I0HaY5bLDf0qpKqFDRldKKIBTh9BQown96+53erOws9l2fTVdRSN7TNHrna
80zSgHODeRbr8Tv6Cyc0eKmNl/OGwwhPwiiU8wdG1YMjEQkkPILRhrNrIp0HAbpezPqe7yQHX38rG65u
g+2erEjz80WaBxRZd/KDqvziq7IOqML+mir+wV3Ye1UxYIkkYtu8KcUodeaGOpqzbmvNMWl72WdJ7jsq
w5bkTUZHhAypcp9KP5ffPY1ctd10S/hCUG1k3UCP5Gw+o/I9Yd9ySpoiWblUcdfshG3DXGI5QrAgVCc3
lLDV0vfuZG4tJ0Azj0OzToBm1W9953CUBnUCOPtAuL8BAAD//01NaJ+SDgAA
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
