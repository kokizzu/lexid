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
need to set ServerUniqueId if you have multiple server
*/

var Now time.Time

var AtomicCounter uint32
var Separator = `~`
var ServerUniqueId = `~0`
var MinCounterLength = len(S.EncodeCB63(math.MaxUint32, 1))
var MinTimeLength = 0
var MinNanoTimeLength = 0

func init() {
	Now = fastime.Now()
}

// generate unique ID (second, smaller)
func ID() string {
	counter := atomic.AddUint32(&AtomicCounter, 1)
	return S.EncodeCB63(Now.Unix(), MinTimeLength) + Separator + S.EncodeCB63(int64(counter), MinCounterLength) + ServerUniqueId
}

// generate unique ID (accurate)
func NanoID() string {
	counter := atomic.AddUint32(&AtomicCounter, 1)
	return S.EncodeCB63(Now.UnixNano(), MinNanoTimeLength) + Separator + S.EncodeCB63(int64(counter), MinCounterLength) + ServerUniqueId
}

type Segments struct {
	Time     int64
	Counter  uint32
	ServerID string
}

func Parse(id string) (*Segments, error) {
	segments := strings.Split(id, Separator)
	if len(segments) != 3 {
		return nil, fmt.Errorf(`invalid lexid or separator: %#v %s`, segments, Separator)
	}
	timePart, timeOk := S.DecodeCB63(segments[0])
	ctrPart, ctrOk := S.DecodeCB63(segments[1])
	var err error
	res := &Segments{
		Time:     timePart,
		Counter:  uint32(ctrPart),
		ServerID: segments[2],
	}
	if !timeOk {
		err = fmt.Errorf(`unable to parse time segment: %#v`, segments[0])
	} else if !ctrOk {
		err = fmt.Errorf(`unable to parse counter segment: %#v`, segments[1])
	}
	return res, err
}

// object-oriented version
type Generator struct {
	AtomicCounter     uint32
	Separator         string
	ServerUniqueId    string
	MinCounterLength  int
	MinTimeLength     int
	MinNanoTimeLength int
}

func NewGenerator(serverUniqueId string) *Generator {
	return &Generator{
		AtomicCounter:     0,
		Separator:         Separator,
		ServerUniqueId:    serverUniqueId,
		MinCounterLength:  len(S.EncodeCB63(math.MaxUint32, 1)),
		MinTimeLength:     0,
		MinNanoTimeLength: 0,
	}
}

func (gen *Generator) ID() string {
	counter := atomic.AddUint32(&gen.AtomicCounter, 1)
	return S.EncodeCB63(Now.Unix(), gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.ServerUniqueId
}

func (gen *Generator) NanoID() string {
	counter := atomic.AddUint32(&gen.AtomicCounter, 1)
	return S.EncodeCB63(Now.UnixNano(), gen.MinTimeLength) + gen.Separator + S.EncodeCB63(int64(counter), gen.MinCounterLength) + gen.ServerUniqueId
}

func (gen *Generator) Parse(id string) (*Segments, error) {
	segments := strings.Split(id, gen.Separator)
	if len(segments) != 3 {
		return nil, fmt.Errorf(`invalid lexid or separator: %#v %s`, segments, gen.Separator)
	}
	timePart, timeOk := S.DecodeCB63(segments[0])
	ctrPart, ctrOk := S.DecodeCB63(segments[1])
	var err error
	res := &Segments{
		Time:     timePart,
		Counter:  uint32(ctrPart),
		ServerID: segments[2],
	}
	if !timeOk {
		err = fmt.Errorf(`unable to parse time segment: %#v`, segments[0])
	} else if !ctrOk {
		err = fmt.Errorf(`unable to parse counter segment: %#v`, segments[1])
	}
	return res, err
}
