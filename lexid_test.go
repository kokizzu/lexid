package lexid_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/google/uuid"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/lexid"
	"github.com/lithammer/shortuuid/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func ExampleRecommendedMinLength() {
	L.Print(`recommended MinDateOffset`, lexid.Offset2020.Unix())
	L.Print(`recommended MinNanoDateOffset`, lexid.Offset2020.UnixNano())
	L.Print(`recommended/default MinCounterLength`, lexid.Config.MinCounterLength)
	print(`recommended/default MinTimeLength`, S.EncodeCB63(lexid.Now.UnixNow(), 0))
	print(`recommended/default MinNanoTimeLength`, S.EncodeCB63(lexid.Now.UnixNanoNow(), 0))

	timeOverflow := `zzzzzz~0~0`
	seg, err := lexid.Parse(timeOverflow, false)
	L.PanicIf(err, `failed parse timeOverflow`)
	L.Describe(seg)
	print(`TimeLength=6 will overflow at`, time.Unix(seg.Time, 0).String())
	// UnixNano will never overflow

	L.Print(`MaxTime`, time.Unix(math.MaxInt64, 0))
	L.Print(`MaxNanoTime`, time.Unix(0, math.MaxInt64))

	lexid.Config.MinTimeLength = 0
	lexid.Config.MinCounterLength = 0
	lexid.Config.MinDateOffset = lexid.Offset2020.Unix()
	print(`offset example`, lexid.ID())
	lexid.Config.MinNanoDateOffset = lexid.Offset2020.UnixNano()
	print(`offset nano example`, lexid.NanoID())
}

func BenchmarkShortuuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := shortuuid.New()
		_ = res
	}
}

func BenchmarkKsuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := ksuid.New().String()
		_ = res
	}
}

func BenchmarkNanoid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res, err := gonanoid.New()
		assert.Nil(b, err)
		_ = res
	}
}

func BenchmarkUuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := uuid.New().String()
		_ = res
	}
}

func BenchmarkTime(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := time.Now().String()
		_ = res
	}
}

func BenchmarkSnowflake(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := I.UToS(snowflake.ID())
		_ = res
	}
}

func BenchmarkLexIdNano(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := lexid.NanoID()
		_ = res
	}
}

// without orderable/sortable property
func BenchmarkLexIdNoLex(b *testing.B) {
	lexid.Config.Separator = ``
	lexid.Config.MinTimeLength = 0
	lexid.Config.MinCounterLength = 0
	for z := 0; z < b.N; z++ {
		res := lexid.ID()
		_ = res
	}
}

func BenchmarkLexId(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := lexid.ID()
		_ = res
	}
}

// without separator
func BenchmarkLexIdNoSep(b *testing.B) {
	lexid.Config.Separator = ``
	for z := 0; z < b.N; z++ {
		res := lexid.ID()
		_ = res
	}
}

func BenchmarkXid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := xid.New().String()
		_ = res
	}
}

func print(prefix, id string) {
	fmt.Println(prefix, id)
	fmt.Println(` len=`, len(id))
}

const N = 1_000_000 // not 10 mil because there's map to check duplicate

func TestLexiId(t *testing.T) {
	m := map[string]bool{}
	lexid.Config.MinTimeLength = 6
	lexid.Config.MinCounterLength = 6
	id := lexid.ID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.ID()
		if past >= id {
			t.Fatalf(`past should be lower: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestLexiIdNano(t *testing.T) {
	m := map[string]bool{}
	lexid.Config.MinNanoTimeLength = 11
	lexid.Config.MinCounterLength = 6
	id := lexid.NanoID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.NanoID()
		if past >= id {
			t.Fatalf(`past should be lower: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate nano lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestOverflow(t *testing.T) {
	const NFast = 2 * 10_000_000 // for processor that could calls 10m times
	lexid.Config.AtomicCounter = math.MaxUint32 - NFast/2
	m := map[string]bool{}
	id := lexid.ID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.ID()
		if past >= id {
			t.Fatalf(`past should be lower: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate nano lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestOverflowNano(t *testing.T) {
	lexid.Config.AtomicCounter = math.MaxUint32 - N/2
	m := map[string]bool{}
	id := lexid.NanoID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.NanoID()
		if past >= id {
			t.Fatalf(`past should be lower: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate nano lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestObject(t *testing.T) {
	m := map[string]bool{}
	gen := lexid.NewGenerator(`~1`)
	id := gen.ID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = gen.ID()
		if past >= id {
			t.Fatalf(`past should be lower: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate lexid`)
		}
		m[id] = true
	}
	print(`last`, id)

}

func TestObjectNano(t *testing.T) {
	m := map[string]bool{}
	gen := lexid.NewGenerator(`~1`)
	id := gen.NanoID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = gen.NanoID()
		if past >= id {
			t.Fatalf(`past should be lower: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate lexid`)
		}
		m[id] = true
	}
	print(`last`, id)

}

func printSeg(seg *lexid.Segments) {
	L.Print(`Time`, seg.ToTime())
	L.Print(`Counter`, seg.Counter)
	L.Print(`Identity`, seg.Identity)
}

func TestParseObject(t *testing.T) {
	gen := lexid.NewGenerator(`~1`)
	gen.AtomicCounter = 123
	id := gen.ID()
	print(`id`, id)
	seg, err := gen.Parse(id, false)
	assert.Nil(t, err)
	assert.Equal(t, gen.AtomicCounter, seg.Counter)
	assert.Equal(t, gen.Identity[1:], seg.Identity)
	printSeg(seg)
}

func TestParseFixedNanoObject(t *testing.T) {
	gen := lexid.NewGenerator(`1`)
	gen.Separator = ``
	gen.AtomicCounter = 123
	gen.Identity = `2`
	id := gen.NanoID()
	print(`id`, id)
	L.Describe(gen)
	seg, err := gen.Parse(id, true)
	assert.Nil(t, err)
	assert.Equal(t, gen.AtomicCounter, seg.Counter)
	assert.Equal(t, gen.Identity, seg.Identity)
	printSeg(seg)
}

func TestParseFixedObject(t *testing.T) {
	gen := lexid.NewGenerator(`1`)
	gen.Separator = ``
	gen.AtomicCounter = 123
	gen.Identity = `2`
	id := gen.ID()
	print(`id`, id)
	seg, err := gen.Parse(id, false)
	assert.Nil(t, err)
	assert.Equal(t, gen.AtomicCounter, seg.Counter)
	assert.Equal(t, gen.Identity, seg.Identity)
	printSeg(seg)
}

func TestFrom(t *testing.T) {
	lexid.Config.Reinit()
	id := lexid.FromUnixCounterIdent(0, 0, `~A`)
	assert.Equal(t, `------~------~A`, id)
	id = lexid.FromUnixCounter(0, 0)
	assert.Equal(t, `------~------~0`, id)
	id = lexid.FromUnix(0)
	assert.Equal(t, `------~-----0~0`, id)
	
	id = lexid.FromNanoCounterIdent(0, 0, `~A`)
	assert.Equal(t, `-----------~------~A`, id)
	id = lexid.FromNanoCounter(0, 0)
	assert.Equal(t, `-----------~------~0`, id)
	id = lexid.FromNano(0)
	assert.Equal(t, `-----------~-----1~0`, id)
}

func TestFromObject(t *testing.T) {
	gen := lexid.NewGenerator(`~X`)
	gen.MinTimeLength = 0
	gen.MinCounterLength = 0
	id := gen.FromUnixCounterIdent(0, 0, `~A`)
	assert.Equal(t, `-~-~A`, id)
	id = gen.FromUnixCounter(0, 0)
	assert.Equal(t, `-~-~X`, id)
	id = gen.FromUnix(0)
	assert.Equal(t, `-~0~X`, id)

	gen.Separator = ``
	gen.Identity = ``
	id = gen.FromUnixCounterIdent(0, 0, `Z`)
	assert.Equal(t, `--Z`, id)
	id = gen.FromUnixCounter(0, 0)
	assert.Equal(t, `--`, id)
	id = gen.FromUnix(0)
	assert.Equal(t, `-1`, id)
}
