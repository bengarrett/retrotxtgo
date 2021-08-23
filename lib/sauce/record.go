package sauce

type (
	record   []byte
	id       [5]byte
	version  [2]byte
	title    [35]byte
	author   [20]byte
	group    [20]byte
	date     [8]byte
	fileSize [4]byte
	dataType [1]byte
	fileType [1]byte
	tInfo1   [2]byte
	tInfo2   [2]byte
	tInfo3   [2]byte
	tInfo4   [2]byte
	comments [1]byte
	tFlags   [1]byte
	tInfoS   [22]byte
)

func (b id) String() string {
	return string(b[:])
}
func (b version) String() string {
	return string(b[:])
}
func (b title) String() string {
	return string(b[:])
}
func (b author) String() string {
	return string(b[:])
}
func (b group) String() string {
	return string(b[:])
}
func (b date) String() string {
	return string(b[:])
}

func (t tInfoS) String() string {
	const nul = 0
	s := ""
	for _, b := range t {
		if b == nul {
			continue
		}
		s += string(b)
	}
	return s
}

func (r record) extract() data {
	i := Scan(r...)
	if i == -1 {
		return data{}
	}
	d := data{
		id:       r.id(i),
		version:  r.version(i),
		title:    r.title(i),
		author:   r.author(i),
		group:    r.group(i),
		date:     r.date(i),
		filesize: r.fileSize(i),
		datatype: r.dataType(i),
		filetype: r.fileType(i),
		tinfo1:   r.tInfo1(i),
		tinfo2:   r.tInfo2(i),
		tinfo3:   r.tInfo3(i),
		tinfo4:   r.tInfo4(i),
		tFlags:   r.tFlags(i),
		tInfoS:   r.tInfoS(i),
	}
	d.comments = r.comments(i)
	d.comnt = r.comnt(d.comments, i)
	return d
}

func (r record) author(i int) author {
	var a author
	const (
		start = 42
		end   = start + len(a)
	)
	for j, c := range r[start+i : end+i] {
		a[j] = c
	}
	return a
}

func (r record) comments(i int) comments {
	return comments{r[i+104]}
}

func (r record) comnt(count comments, sauceIndex int) (block comnt) {
	block = comnt{
		count: count,
	}
	if int(unsignedBinary1(count)) == 0 {
		return block
	}
	id, l := []byte(comntID), len(r)
	var backwardsLoop = func(i int) int {
		return l - 1 - i
	}
	// search for the id sequence in b
	for i := range r {
		if i > comntLineSize*comntMaxLines {
			break
		}
		i = backwardsLoop(i)
		if i < comntLineSize {
			break
		}
		if i >= sauceIndex {
			continue
		}
		// do matching in reverse
		if r[i-1] != id[4] {
			continue // T
		}
		if r[i-2] != id[3] {
			continue // N
		}
		if r[i-3] != id[2] {
			continue // M
		}
		if r[i-4] != id[1] {
			continue // O
		}
		if r[i-5] != id[0] {
			continue // C
		}
		block.index = i
		block.length = sauceIndex - block.index
		block.lines = r[i : i+block.length]
		return block
	}
	return block
}

func (r record) dataType(i int) dataType {
	return dataType{r[i+94]}
}

func (r record) date(i int) date {
	var d date
	const (
		start = 82
		end   = start + len(d)
	)
	for j, c := range r[start+i : end+i] {
		d[j] = c
	}
	return d
}

func (r record) fileSize(i int) fileSize {
	return fileSize{r[i+90], r[i+91], r[i+92], r[i+93]}
}

func (r record) fileType(i int) fileType {
	return fileType{r[i+95]}
}

func (r record) group(i int) group {
	var g group
	const (
		start = 62
		end   = start + len(g)
	)
	for j, c := range r[start+i : end+i] {
		g[j] = c
	}
	return g
}

func (r record) id(i int) id {
	return id{r[i+0], r[i+1], r[i+2], r[i+3], r[i+4]}
}

func (r record) tFlags(i int) tFlags {
	return tFlags{r[i+105]}
}

func (r record) title(i int) title {
	var t title
	const (
		start = 7
		end   = start + len(t)
	)
	for j, c := range r[start+i : end+i] {
		t[j] = c
	}
	return t
}

func (r record) tInfo1(i int) tInfo1 {
	return tInfo1{r[i+96], r[i+97]}
}

func (r record) tInfo2(i int) tInfo2 {
	return tInfo2{r[i+98], r[i+99]}
}

func (r record) tInfo3(i int) tInfo3 {
	return tInfo3{r[i+100], r[i+101]}
}

func (r record) tInfo4(i int) tInfo4 {
	return tInfo4{r[i+102], r[i+103]}
}

func (r record) tInfoS(i int) tInfoS {
	var s tInfoS
	const (
		start = 106
		end   = start + len(s)
	)
	for j, c := range r[start+i : end+i] {
		if c == 0 {
			continue
		}
		s[j] = c
	}
	return s
}

func (r record) version(i int) version {
	return version{r[i+5], r[i+6]}
}
