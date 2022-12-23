package lexid

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kokizzu/gotro/S"
	"github.com/kpango/fastime"
)

/* generate based on 3 parts:
 1. current second
 2. atomic increment
 3. server unique id
need to set Identity if you have multiple server/instance
*/

var lastSec int64
var incNano int64

var Config *Generator

var DefaultSeparator = `~`
var DefaultMinCounterLength = 6
var DefaultMinTimeLength = 6
var DefaultMinNanoTimeLength = 11
var DefaultMinDateOffset = int64(0)
var DefaultMinNanoDateOffset = int64(0)
var Offset2020 time.Time // MinDateOffset = Offset2020.Unix() or MinNanoDateOffset = Offset2020.UnixNano()
var SeparatorIdentity = `~0`

func init() {
	DefaultMinCounterLength = len(S.EncodeCB63(int64(math.MaxUint32), 0))
	DefaultMinTimeLength = len(S.EncodeCB63(time.Now().Unix(), 0))
	DefaultMinNanoTimeLength = len(S.EncodeCB63(time.Now().UnixNano(), 0))

	Config = NewGenerator(SeparatorIdentity)

	Offset2020 = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	Config.MinTimeLength = DefaultMinTimeLength
	Config.MinNanoTimeLength = DefaultMinNanoTimeLength
	Config.MinDateOffset = DefaultMinDateOffset
	Config.MinNanoDateOffset = DefaultMinNanoDateOffset
}

func Parse(id string, isNano bool) (*Segments, error) {
	return Config.Parse(id, isNano)
}

func FromUnixCounterIdent(time int64, counter uint32, ident string) string {
	return Config.FromUnixCounterIdent(time, counter, ident)
}

func FromUnixCounter(time int64, counter uint32) string {
	return Config.FromUnixCounter(time, counter)
}

func FromUnix(time int64) string {
	return Config.FromUnix(time)
}

func FromNanoCounterIdent(time int64, counter uint32, ident string) string {
	return Config.FromNanoCounterIdent(time, counter, ident)
}

func FromNanoCounter(time int64, counter uint32) string {
	return Config.FromNanoCounter(time, counter)
}

func FromNano(time int64) string {
	return Config.FromNano(time)
}

// generate unique ID (second, smaller)
func ID() string {
	return Config.ID()
}

// generate unique ID (accurate)
func NanoID() string {
	return Config.NanoID()
}

type Segments struct {
	Time      int64
	Counter   uint32
	Identity  string
	IsNano    bool
	Generator *Generator
}

func (s *Segments) ToTime() time.Time {
	if s.IsNano {
		return time.Unix(0, s.Time)
	}
	return time.Unix(s.Time, 0)
}

func (s *Segments) ToID() string {
	if s.Generator == nil {
		s.Generator = Config
	}
	if s.IsNano {
		return s.Generator.FromNanoCounterIdent(s.Time, s.Counter, s.Identity)
	}
	return s.Generator.FromUnixCounterIdent(s.Time, s.Counter, s.Identity)
}

// object-oriented version
type Generator struct {
	AtomicCounter     uint32
	Separator         string
	Identity          string
	MinCounterLength  int
	MinTimeLength     int
	MinNanoTimeLength int
	MinDateOffset     int64
	MinNanoDateOffset int64
}

func NewGenerator(uniqStr string) *Generator {
	return &Generator{
		AtomicCounter:     0,
		Separator:         DefaultSeparator,
		Identity:          uniqStr,
		MinCounterLength:  DefaultMinCounterLength,
		MinTimeLength:     DefaultMinTimeLength,
		MinNanoTimeLength: DefaultMinNanoTimeLength,
		MinDateOffset:     DefaultMinDateOffset,
		MinNanoDateOffset: DefaultMinNanoDateOffset,
	}
}

func (gen *Generator) Reinit() {
	gen.AtomicCounter = 0
	gen.Separator = DefaultSeparator
	gen.MinCounterLength = DefaultMinCounterLength
	gen.MinTimeLength = DefaultMinTimeLength         // >=6
	gen.MinNanoTimeLength = DefaultMinNanoTimeLength // >=11
	gen.MinDateOffset = DefaultMinDateOffset
	gen.MinNanoDateOffset = DefaultMinNanoDateOffset
}

func (gen *Generator) ID() string {
	now := fastime.UnixNow()
	counter := atomic.AddUint32(&gen.AtomicCounter, 1)
	if now != lastSec { // reset to 0 if not the same second
		atomic.SwapInt64(&lastSec, now)
		atomic.SwapUint32(&gen.AtomicCounter, 0) // ignore old value
		counter = atomic.AddUint32(&gen.AtomicCounter, 1)
	}
	return S.EncodeCB63(now-gen.MinDateOffset, gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.Identity
}

func (gen *Generator) NanoID() string {
	counter := atomic.AddUint32(&gen.AtomicCounter, 1)
	if counter == 0 { // add 1 nanosecond everytime generating 4 million IDs
		atomic.AddInt64(&incNano, 1)
	}
	return S.EncodeCB63(fastime.UnixNanoNow()+incNano-gen.MinNanoDateOffset, gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.Identity
}

func (gen *Generator) FromUnixCounterIdent(time int64, counter uint32, ident string) string {
	return S.EncodeCB63(time-gen.MinDateOffset, gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + ident
}

func (gen *Generator) FromUnixCounter(time int64, counter uint32) string {
	return S.EncodeCB63(time-gen.MinDateOffset, gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.Identity
}

func (gen *Generator) FromUnix(time int64) string {
	counter := atomic.AddUint32(&gen.AtomicCounter, 1)
	return S.EncodeCB63(time-gen.MinDateOffset, gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.Identity
}

func (gen *Generator) FromNanoCounterIdent(time int64, counter uint32, ident string) string {
	return S.EncodeCB63(time-gen.MinNanoDateOffset, gen.MinNanoTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + ident
}

func (gen *Generator) FromNanoCounter(time int64, counter uint32) string {
	return S.EncodeCB63(time-gen.MinNanoDateOffset, gen.MinNanoTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.Identity
}

func (gen *Generator) FromNano(time int64) string {
	counter := atomic.AddUint32(&gen.AtomicCounter, 1)
	return S.EncodeCB63(time-gen.MinNanoDateOffset, gen.MinNanoTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.Identity
}

func (gen *Generator) Parse(id string, isNano bool) (*Segments, error) {
	var segments []string
	if gen.Separator == `` {
		// try parse as unixnano
		start := gen.MinNanoTimeLength
		end := start + gen.MinCounterLength
		if len(id) < end {
			// try parse as unix
			start = gen.MinTimeLength
			end = start + gen.MinCounterLength
			if len(id) < end {
				return nil, fmt.Errorf(`invalid lexid length: %s %d < %d+%d`, id, len(id), gen.MinTimeLength, gen.MinCounterLength)
			}
		}
		segments = []string{
			id[:start],
			id[start:end],
			id[end:],
		}
	} else {
		segments = strings.Split(id, gen.Separator)
		if len(segments) != 3 {
			return nil, fmt.Errorf(`invalid lexid or separator: %#v %s`, segments, gen.Separator)
		}
	}
	timePart, timeOk := S.DecodeCB63[int64](segments[0])
	ctrPart, ctrOk := S.DecodeCB63[int64](segments[1])
	var err error
	if isNano {
		timePart += gen.MinNanoDateOffset
	} else {
		timePart += gen.MinDateOffset
	}
	res := &Segments{
		Time:     timePart,
		Counter:  uint32(ctrPart),
		Identity: segments[2],
		IsNano:   isNano,
	}
	if !timeOk {
		err = fmt.Errorf(`unable to parse time segment: %#v`, segments[0])
	} else if !ctrOk {
		err = fmt.Errorf(`unable to parse counter segment: %#v`, segments[1])
	}
	return res, err
}
