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

	"/templates/build/ecs-app.json": {
		local:   "templates/build/ecs-app.json",
		size:    5748,
		modtime: 1447920981,
		compressed: `
H4sIAAAJbogA/6xYW2/bNhR+968g9NQWbeN4F3TCNsBxncVD0hmxlz5sfqBl2iYikQJJpUuH/PedQ8m2
SFGym6U3VNJ3Dj+ey8fD/NsjJBp+ns1ZlqfUsEupMmrumNJciigm0aB/3n/X/wn+RG8R+5HpRPHcVJ9H
ioGVJlSQ8WhGZkw98ITB44pcS7oiFzSlImGKcEEo2fAHVgJHaaENU3+L0uuUKpoxeKHBKXLClWRyz9Qk
oxu2fwmv5485Pkczo7jYWPMS7hKbbxnhaEuMJKoQSMDAu0QKQ7lgKrKGT6V9NNq9/gQ8nrmcAFMi13aZ
lSVfW61muKZFatCI5nkLi1FehEh8KrKl66tJwkKQRpIXpBDcaLJiK55AmlYYCycI5NV5f/A9+YWcwzvF
XodoIiLM8oZlUj0+k2hG/8l4kRGayUIYJJxZdx7BBCqLpqlE/iF2g34/TG4qlXkmtRxMifTqxdZ0ijW9
3NU0RFPmTIRofXBZQclXFf8CtXVwhgxSWohkaz9APaWYaLCDcndL6y5PQktD68fxeDSI47vpKI4nq04W
gIGFqLGrYSze7WPBNXaZgH00lp4VS8GMDi1/zbX5+cChRCKNXzt56NJjyWUtVYDPoetf7dBV+OhSPtit
vPZosqRQ3Dz+pmSRd8fKgXZGbUh0hSUbBO/Zjq8vSvpbCuqZJEzrXXfW9JFMhDa4IV1x7VV8oz8Kkxe1
oLaU2B1Ni7qewatbtkZmNXz17ckvWVTwnYCf7NQxcj3vyd8yLQuFu9rTn1N9D+3DQbDK0LWFH/7xsPvY
TxW0ozKcaZfcXhMORoj4a48gNbS18E4Bf5vuaeGAnt66nlwlb3WEsE4/Y62ZMJymaGZUcRDDCuAflP5S
9eO0c6WGqreSrpCd3lCGb0CYQBj8mOMvfw2wuJLaeOLdycJiG0CPh7VsOxue477X9bzohb4swn12ffFc
5Tle+RbmSdKsKUgoRqCUayXLkxi6rFZeeHaAxLnx2sUJz5VeIOiunE7ERoHCdTfdJIdtGJnIssaT3C/x
SyWzltSdkDa/Mufy5XyN+EpNMHVR/739fdaPvqUG2rXWLYGUasOTAxa6Ko4d06Ml0TyN63vefT2e0iPJ
3B0IfnX3WoIYNm+3PUS0RhDHCSbKW0QHt3rAXq4Gdkf1y3mst8PVfD49HoArRlOzHW1Zcu8neE7VhjWJ
RZcijn+XXATEOYp8EfURpCQWN4AhZf9f2n3mwxbtils3r0LyON+CAG1lijo2cAB/im0Tct53MBMBc9KD
PXy/c7/MecZkgXH9IdgzsFXBEtTejwq2DA07lTCk++drNBZ0mbJV8HSvLfJj/7BKm5RUl/DuGWoHOj48
NYZK+7p1inS3D2cPV2w1wjsehrXerLUuPNKwbbdzn823zGXeLPCh/7ZdJI4s2zL11tPj/7/etLcy9b3X
PFeJsqBgiFsn57ojD3RyEXnc3EKaDG/i2CKOVtFQ6yKz3sryh2kUnoWvR3AVhnt+9cEVm0bwx+s1NJYl
k6byS0OCgAuHEzK3XRuYMw9t0lQ1+MwS/Z5m9KsU9It+n8isqVMLfwhsUBgmVU4CwqmNjg9h6Za3ehEt
glUwpWaLsTirD2421mUuOrqrzEhV47jvd9oTBwfYkjqLaE9fc93K5FgaLagjjvY7Kwcj/BHAcj8YlXPv
kr0J+my3UmyDc4Ta371x7gwOWaf4u/W9feZme5K3ZHDCFgA0LMxWKv6VBUfugN0iGOPdpRxT8eb4jae9
QFun3R7+fer9FwAA//8uUjJJdBYAAA==
`,
	},

	"/templates/build/ecs-cluster-byo-vpc.json": {
		local:   "templates/build/ecs-cluster-byo-vpc.json",
		size:    9236,
		modtime: 1447920982,
		compressed: `
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
`,
	},

	"/templates/build/ecs-cluster.json": {
		local:   "templates/build/ecs-cluster.json",
		size:    11899,
		modtime: 1447920982,
		compressed: `
H4sIAAAJbogA/9RabXPbuBH+7l+Bqp370IksiX7JnWbajqI4qXrnRGMpuenVmRsIhCxMSIADgHGUjv/7
LUBSIkiQop04yfkusQTsLh7sG3aB/P8Iod7k18WSxkmENX0hZIz1WyoVE7w3Rr1gOBr2hz/B/70nhvY5
VUSyROfTU0mBSyHM0STGnwRHFxFWmhE0FVxjxqlECyo/MEIRiVKlqbzmmaQ5ljimMKBAkMEBYz/T7SsY
3A3UF1xuKFJqg97TbYKZRKmiIdICYUKoUkjDNCUKMa405jBk17KCltvECDa7HY8vpsF4DKvNQYb9YJe1
pHcZR2+S6o2Q7BMN3ygA+UZGLbBe2084Qn004SiVkcGUUMlEyAiOoi0KxS2PBA4teryT/TtsRKG1FHEd
6UJLxm/248/pGqeRNlMu1Fm+25yxTXcaSJBY7xRkcIIO0VpIq7wWxbXA0cFxzIgU+8lJFIlbGr7FUUqN
gf+XT8AUOT2OsLyhO+Js7KNvMPCOnnpHf6yP1nBlYyoGk1TpaMjSuDwYn/gHa6vAWH1pGCzA54PvHJtd
4o8LcIAD5orxRxanMeJpvIJQKhlOGctFOOVkUzfVK0vuM9W5AwJWZJKGU5xgwvT2AJgwo0YkJ0d4DeH7
ABAnriYY76IJBjr4kpoYuZoQ5D2V/05XJth5JQc1R4AD878itQnJ8iPIhSagMsEIJHcI5B2Kixiz6KEQ
qGFGOAylSYmfhWOOlboVMnwolCTnb0PxSlyQjbBpRKa0Dd1RjrD3OtVJqksHx8V0Mc2OlzJSm31KAzB0
RddGYIk+n7tzFLCgJJXg4i+lSJNuIl0Wr9S3CekmyxD6caUrTksb90t5wcfj/wjGnbwL409KCQo5U6jE
XtmWXXHYc6bvntyDd/QZvEGFt/Tt3VH1013FSy5xkoCHltwEzv4regMOuhSTy1lZi6nqUyhc+iNXk0A1
e27g4Jj1w5A8XZ2vnu5M86TEfUuBO2jhxic/hesROalz0zTjblt7PVqdnq9+DOrcOOlzISG4DsE/GT59
ehacUa8IJdJcRNsegpPVKV2dDnsNGr+iIEcS6kRmUJQncynWLPJmVluUzSaX8FeFeAcSBqCg0owqF+Ac
640RMSifvlciqpQdVU8r/AwWNcRlP7tr8qyi3MpZWvdhKQ6CnyiVxtTQzkXEyBbyI3zn2qEyca+hzM4n
DsTtxXpNic2atgbrPakSzCFtE5bgqLJIsVRWrtcWyqcpCY6xLfXxrTomULfWqN5VRu5qECakOCWUVuO9
ErpEe0ma3/aXmOMbGmYKnUhe9YMelnwM4McMx2P7IbGkA5XtvC8BySBrZ8B7d41MrhiolGHURdvoJxaE
a3WPr2RQO7i6pcs7pF3VX+He0X0NX9oZsu4rPRKJNLzFmmzG81RfUigXyHOscU1KQbu2zSfIG2elxIoC
VPK+SCo+PvDGHfESQ6qvkLyrId5Jg039/b7+9qhpxV/DVNpW+MspXErM0LCKBTSbcOTVyhZHSonQ+ZJx
HfTCt/Ppb4LTWQguxNbMwj2sEE8d4eSFNq7Rg7iCBuWX7PmL7RogwtfsJpXW9yr3D2WxHuqe103qHU1Z
TDHrZ621hQ5rPutlbW7myiKqVF5RNpRarfozNQv0nFMunyoKUghO1S+atCqRcS58A4loojO17ur/9ojZ
eae9eQIL5NmvUkVnMb5gN7x2zPWWLKZQ7ZgF58vR2aUDrTcVKd+3h2UU+7XfJCEg961ciiZQjfmV0VYx
gA8U6VvN+P7IHTlYwNzPTP7MPcKdnGNoNc1eso1U9/ErZvo1dzWh6kq+86Yin6d3yiY+xi5FkCDMaDRd
gUpnySRrXWstIdofepeCMy1sA+qhcnqxbgnb276VNeRmjlkMNcYsrNrVdF4vGA9nHJqP2rFY7T+6dUJW
yxlbayeVl+oHD7GG28Lykg6JXwiOm8v6sihP/e8VWL/6LUtxbmirrObKxhYWHmM8w4qen9Yq3aYW2c7V
Cg1PcfPXvwxWjA9WWG1Q/+MHes23aZxdSUUR6m8R1JV9sub9lRBaaYmTaz4QiR7AuOU0cwx8GPU/oH5f
mUIHeSqceoFecQxbIrm68ehoxwpryTwnIE+w2mnjaQ/C4nPSJiAedSibpFCfor/961vpxFNHfU2d1OrX
o6Zvdy3nE9TaOKwGRAZnagrtF7tCG9ptVmsQiPWIesyY4FWeprE3oJoM4LQ3f44buAvJut6R7KZbgjKn
8FjA/Pips7L49+kvbxbLi6t/NPD6bZjzN14XVn885sxF1I2a/VS7ZP+Yz01iEdrjdzgcng+HvrZI3HJb
lvckZB8fwU3eIGQERwfWzAxMpODHGwiTaDuoPGB9J9Z2kjJNr2F5BQn2E7puEtXF+J6HwPs7QQ/98AOi
HyHfD685MQ+EkOAgnjjgW38X+P6JBjpOqpa95vEH7wQabERMIeKDvnnyGBwrtakzkw04KwIv7UoOnosK
onE3pkcMr6ePEl6tRRwRcYx56M2zawjDTX+vAKuVxuCzcuzt2KHobQfYfug4fUPjs0n98d0lPdgjWLLK
C9Pb+XQAslAhCrmykH1Hqdfo9TcWt5x0kM34Td6KtLQQswRga0FEZJsRklS73RdSxHMhTSoMgsrcUjTN
TFkoZ9aNhsf2v8HwPvdJDfvoapOCo5tpGtXc0Fc5vUSL+sqqG5ZvKAqtnZ+dnZyVrZf1u85O7ovNq87K
411de+COh5VlbPosEuS94R0Vdh2d+3rcJeXwy96zhPlb6EPuaGwh3HBH09pzNpTR/oTg98GX0NXf4m27
4mZwTktOdUHcrkTvAhMNWDeV226/jeocB41WAdjoT8UGvF7ePRf5ozm/P20PYEvU4cblA2YRXrEInN/c
5fo65wWNsrt/t9waHri0MKwvqZ781nAuHe6Nmo8f/31GPaRGEFLB6cMOghblj7658kd/AuUHj6P84Jsr
P/gTKP/kSyv/SqS65fnS6t7SLPGqy4P3l4Cz/4dBh0AdxgMlpWbcXoY42txXXCVdflb+32upkT9XdreT
oItlMtL9ysVNe6dr+Yy5uXiqPug94mZH38dmR19ls8H3sdngsTZ7ZP7cHf0RAAD//6xJNR17LgAA
`,
	},

	"/templates/build/vpc.json": {
		local:   "templates/build/vpc.json",
		size:    3348,
		modtime: 1447911594,
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
