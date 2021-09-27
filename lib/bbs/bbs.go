package bbs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type BBS int

const (
	ANSI BBS = iota
	Celerity
	PCBoard
	Renegade
	Telegard
	Wildcat
	WWIVHash
	WWIVHeart
)

const (
	PCBClear = `@CLS@`
)

// String returns the BBS color format name and toggle characters.
func (b BBS) String() string {
	if b < ANSI || b > WWIVHeart {
		return ""
	}
	return [...]string{
		"ANSI ←[",
		"Celerity |",
		"PCBoard @",
		"Renegade |",
		"Telegard `",
		"Wildcat! @X",
		"WWIV |#",
		"WWIV ♥",
	}[b]
}

// Bytes returns the BBS color code toggle characters.
func (b BBS) Bytes() []byte {
	const (
		etx               byte = 3  // CP437 ♥
		esc               byte = 27 // CP437 ←
		hash                   = byte('#')
		atSign                 = byte('@')
		grave                  = byte('`')
		leftSquareBracket      = byte('[')
		verticalBar            = byte('|')
		upperX                 = byte('X')
	)
	switch b {
	case ANSI:
		return []byte{esc, leftSquareBracket}
	case Celerity, Renegade:
		return []byte{verticalBar}
	case PCBoard:
		return []byte{atSign, upperX}
	case Telegard:
		return []byte{grave}
	case Wildcat:
		return []byte{atSign}
	case WWIVHash:
		return []byte{verticalBar, hash}
	case WWIVHeart:
		return []byte{etx}
	default:
		return nil
	}
}

func Find(r io.Reader) BBS {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		ts := bytes.TrimSpace(b)
		if ts == nil {
			return -1
		}
		const l = len(PCBClear)
		if len(ts) > l {
			if bytes.Equal(ts[0:l], []byte(PCBClear)) {
				b = ts[l:]
			}
		}
		switch {
		case bytes.Contains(b, ANSI.Bytes()):
			return ANSI
		case bytes.Contains(b, Celerity.Bytes()):
			if f := findRenegade(b); f == Renegade {
				return Renegade
			}
			if f := findCelerity(b); f == Celerity {
				return Celerity
			}
			return -1
		case bytes.Contains(b, PCBoard.Bytes()):
			return findPCBoard(b)
		case bytes.Contains(b, Telegard.Bytes()):
			return findTelegard(b)
		case bytes.Contains(b, Wildcat.Bytes()):
			return findWildcat(b)
		case bytes.Contains(b, WWIVHash.Bytes()):
			return findWWIVHash(b)
		case bytes.Contains(b, WWIVHeart.Bytes()):
			return findWWIVHeart(b)
		}
	}
	return -1
}

func findCelerity(b []byte) BBS {
	const (
		bb = byte('B')
		c  = byte('C')
		d  = byte('D')
		g  = byte('G')
		k  = byte('K')
		m  = byte('M')
		r  = byte('R')
		s  = byte('S')
		y  = byte('Y')
		w  = byte('W')
	)
	codes := []byte{bb, c, d, g, k, m, r, s, y, w}
	for _, code := range codes {
		if bytes.Contains(b, []byte{Celerity.Bytes()[0], code}) {
			return Celerity
		}
	}
	return -1
}

func findPCBoard(b []byte) BBS {
	/*
	   const inRange = (a = -1, b = -1) => {
	     if (a >= 48 && b >= 48 && a <= 70 && b <= 70) return true
	     return false
	   }
	   48 = 0; 0030
	   70 = F; 0046
	   "@X<Background><Foreground>"
	*/
	const first, last = 0, 15
	const hexxed = "%X%X"
	for bg := first; bg <= last; bg++ {
		for fg := first; fg <= last; fg++ {
			subslice := []byte(fmt.Sprintf(hexxed, bg, fg))
			subslice = append(PCBoard.Bytes(), subslice...)
			if bytes.Contains(b, subslice) {
				return PCBoard
			}
		}
	}
	return -1
}

func findRenegade(b []byte) BBS {
	// |00 -> |23 (strconv)
	const first, last = 0, 23
	const leadingZero = "%01d"
	for i := first; i <= last; i++ {
		subslice := []byte(fmt.Sprintf(leadingZero, i))
		subslice = append(Renegade.Bytes(), subslice...)
		if bytes.Contains(b, subslice) {
			return Renegade
		}
	}
	return -1
}

func findTelegard(b []byte) BBS {
	// |00 -> |23 (strconv)
	const first, last = 0, 23
	const leadingZero = "%01d"
	for i := first; i <= last; i++ {
		subslice := []byte(fmt.Sprintf(leadingZero, i))
		subslice = append(Telegard.Bytes(), subslice...)
		if bytes.Contains(b, subslice) {
			return Telegard
		}
	}
	return -1
}

func findWildcat(b []byte) BBS {
	const first, last = 0, 15
	for bg := first; bg <= last; bg++ {
		for fg := first; fg <= last; fg++ {
			subslice := []byte(fmt.Sprintf("%s%X%X%s",
				Wildcat.Bytes(), bg, fg, Wildcat.Bytes()))
			if bytes.Contains(b, subslice) {
				return Wildcat
			}
		}
	}
	return -1
}

func findWWIVHash(b []byte) BBS {
	const first, last = 0, 9
	for i := first; i <= last; i++ {
		subslice := append(WWIVHash.Bytes(), []byte(strconv.Itoa(i))...)
		if bytes.Contains(b, subslice) {
			return WWIVHash
		}
	}
	return -1
}

func findWWIVHeart(b []byte) BBS {
	const first, last = 0, 9
	for i := first; i <= last; i++ {
		subslice := append(WWIVHeart.Bytes(), []byte(strconv.Itoa(i))...)
		if bytes.Contains(b, subslice) {
			return WWIVHeart
		}
	}
	return -1
}
