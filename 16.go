package main

import (
	"bytes"
	"io"
	"log"
	"math"
	"strconv"
)

func init() {
	addSolutions(16, problem16)
}

func problem16(ctx *problemContext) {
	b, err := io.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	s := &bitStream{hex: bytes.TrimSpace(b)}
	ctx.reportLoad()

	pkt := s.parseBitsPacket()
	if !s.restZeros() {
		panic("trailing non-zero")
	}
	var versionSum int64
	pkts := QueueOf(pkt)
	for {
		pkt, ok := pkts.Pop()
		if !ok {
			break
		}
		versionSum += int64(pkt.version())
		for _, child := range pkt.children() {
			pkts.Push(child)
		}
	}
	ctx.reportPart1(versionSum)

	ctx.reportPart2(pkt.eval())
}

type bitStream struct {
	cur     uint64
	numBits int
	numRead int64
	hex     []byte
}

func (s *bitStream) pop(n int) uint64 {
	for s.numBits < n {
		h := s.hex[0]
		s.hex = s.hex[1:]
		v, err := strconv.ParseUint(string(h), 16, 4)
		if err != nil {
			panic(err)
		}
		s.cur = (s.cur << 4) | v
		s.numBits += 4
	}
	r := s.cur >> (s.numBits - n)
	s.cur &= (1 << (s.numBits - n)) - 1
	s.numBits -= n
	s.numRead += int64(n)
	return r
}

func (s *bitStream) restZeros() bool {
	if s.cur != 0 {
		return false
	}
	for _, h := range s.hex {
		if h != '0' {
			return false
		}
	}
	return true
}

type bitsPacket interface {
	version() uint8
	children() []bitsPacket
	eval() uint64
}

type literalPacket struct {
	ver uint8
	v   uint64
}

func (p *literalPacket) version() uint8 {
	return p.ver
}

func (p *literalPacket) children() []bitsPacket {
	return nil
}

func (p *literalPacket) eval() uint64 {
	return p.v
}

type opPacket struct {
	ver uint8
	op  uint8
	ch  []bitsPacket
}

func (p *opPacket) version() uint8 {
	return p.ver
}

func (p *opPacket) children() []bitsPacket {
	return p.ch
}

func (p *opPacket) eval() uint64 {
	switch p.op {
	case 0:
		var sum uint64
		for _, c := range p.ch {
			sum += c.eval()
		}
		return sum
	case 1:
		prod := uint64(1)
		for _, c := range p.ch {
			prod *= c.eval()
		}
		return prod
	case 2:
		min := uint64(math.MaxUint64)
		for _, c := range p.ch {
			if v := c.eval(); v < min {
				min = v
			}
		}
		return min
	case 3:
		var max uint64
		for _, c := range p.ch {
			if v := c.eval(); v > max {
				max = v
			}
		}
		return max
	case 5:
		if len(p.ch) != 2 {
			panic("need 2 children")
		}
		if p.ch[0].eval() > p.ch[1].eval() {
			return 1
		}
		return 0
	case 6:
		if len(p.ch) != 2 {
			panic("need 2 children")
		}
		if p.ch[0].eval() < p.ch[1].eval() {
			return 1
		}
		return 0
	case 7:
		if len(p.ch) != 2 {
			panic("need 2 children")
		}
		if p.ch[0].eval() == p.ch[1].eval() {
			return 1
		}
		return 0
	default:
		panic("bad op")
	}
}

func (s *bitStream) parseBitsPacket() bitsPacket {
	ver := uint8(s.pop(3))
	typeID := uint8(s.pop(3))
	if typeID == 4 {
		// Literal
		pkt := &literalPacket{ver: ver}
		for {
			v := s.pop(5)
			pkt.v <<= 4
			pkt.v |= v & ((1 << 4) - 1)
			if v&(1<<4) == 0 {
				break
			}
		}
		return pkt
	}
	lengthTypeID := uint8(s.pop(1))
	if lengthTypeID == 0 {
		// 15 bit total length
		subLen := s.pop(15)
		target := s.numRead + int64(subLen)
		pkt := &opPacket{ver: ver, op: typeID}
		for {
			pkt.ch = append(pkt.ch, s.parseBitsPacket())
			if s.numRead > target {
				panic("read too many bits for child")
			}
			if s.numRead == target {
				return pkt
			}
		}
	}
	// 11-bit number of sub-packets
	numSub := int(s.pop(11))
	pkt := &opPacket{ver: ver, op: typeID}
	for i := 0; i < numSub; i++ {
		pkt.ch = append(pkt.ch, s.parseBitsPacket())
	}
	return pkt
}
