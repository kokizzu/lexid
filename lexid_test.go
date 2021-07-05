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
	"github.com/matoous/go-nanoid/v2"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func BenchmarkKsuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := ksuid.New().String()
		_ = res
	}
}

func BenchmarkTime(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := time.Now().String()
		_ = res
	}
}

func BenchmarkShortuuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := shortuuid.New()
		_ = res
	}
}

func BenchmarkXid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := xid.New().String()
		_ = res
	}
}

func BenchmarkSnowflake(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := I.UToS(snowflake.ID())
		_ = res
	}
}

func BenchmarkUuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := uuid.New().String()
		_ = res
	}
}

func BenchmarkNanoid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res, err := gonanoid.New()
		L.PanicIf(err, `error generating nanoid`)
		_ = res
	}
}

func BenchmarkLexId(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := lexid.ID()
		_ = res
	}
}

func BenchmarkLexIdNano(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := lexid.NanoID()
		_ = res
	}
}

func print(prefix, id string) {
	fmt.Println(prefix, id)
	fmt.Println(` len=`, len(id))
}

const N = 1_000_000

func TestLexiId(t *testing.T) {
	m := map[string]bool{}
	id := lexid.ID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.ID()
		if past >= id {
			t.Errorf(`past should be lower or equal: %s >= %s`, past, id)
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
	id := lexid.NanoID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.NanoID()
		if past >= id {
			t.Errorf(`past should be lower or equal: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate nano lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestOverflow(t *testing.T) {
	lexid.AtomicCounter = math.MaxUint32 - N/2
	m := map[string]bool{}
	id := lexid.ID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.ID()
		if past >= id {
			t.Logf(`past should be lower or equal: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate nano lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestOverflowNano(t *testing.T) {
	lexid.AtomicCounter = math.MaxUint32 - N/2
	m := map[string]bool{}
	id := lexid.NanoID()
	print(`first`, id)
	for z := 0; z < N; z++ {
		past := id
		id = lexid.NanoID()
		if past >= id {
			t.Logf(`past should be lower or equal: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate nano lexid`)
		}
		m[id] = true
	}
	print(`last`, id)
}

func TestRecommendedMinTime(t *testing.T) {
	print(`recommended MinTimeLength`, S.EncodeCB63(lexid.Now.Unix(), 0))
	print(`recommended MinNanoTimeLength`, S.EncodeCB63(lexid.Now.UnixNano(), 0))
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
			t.Errorf(`past should be lower or equal: %s >= %s`, past, id)
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
			t.Errorf(`past should be lower or equal: %s >= %s`, past, id)
		}
		if m[id] {
			panic(`duplicate lexid`)
		}
		m[id] = true
	}
	print(`last`, id)

}

func TestParse(t *testing.T) {
	id := lexid.ID()
	print(`id`, id)
	seg, err := lexid.Parse(id)
	L.IsError(err, `failed parse id`)
	L.Print(`Time`, time.Unix(seg.Time, 0))
	L.Print(`Counter`, seg.Counter)
	assert.Equal(t, lexid.AtomicCounter, seg.Counter)
	L.Print(`ServerID`, seg.ServerID)
	assert.Equal(t, lexid.ServerUniqueId[1:], seg.ServerID)
}

func TestParseOO(t *testing.T) {
	gen := lexid.NewGenerator(`~1`)
	gen.AtomicCounter = 123
	id := gen.ID()
	print(`id`, id)
	seg, err := gen.Parse(id)
	L.IsError(err, `failed parse id`)
	L.Print(`Time`, time.Unix(seg.Time, 0))
	L.Print(`Counter`, seg.Counter)
	assert.Equal(t, gen.AtomicCounter, seg.Counter)
	L.Print(`ServerID`, seg.ServerID)
	assert.Equal(t, gen.ServerUniqueId[1:], seg.ServerID)
}

func TestParseFixedNanoOO(t *testing.T) {
	gen := lexid.NewGenerator(`1`)
	gen.Separator = ``
	gen.AtomicCounter = 123
	gen.ServerUniqueId = `2`
	L.Print(`MaxTime`, time.Unix(0, math.MaxInt64))
	id := gen.NanoID()
	print(`id`, id)
	L.Describe(gen)
	seg, err := gen.Parse(id)
	L.IsError(err, `failed parse id`)
	L.Print(`Time`, time.Unix(0, seg.Time))
	L.Print(`Counter`, seg.Counter)
	assert.Equal(t, gen.AtomicCounter, seg.Counter)
	L.Print(`ServerID`, seg.ServerID)
	assert.Equal(t, gen.ServerUniqueId, seg.ServerID)
}

func TestParseFixedNano(t *testing.T) {
	defer func() {
		// restore configuration
		lexid.Separator = `~`
		lexid.ServerUniqueId = `~0`
	}()
	lexid.Separator = ``
	lexid.AtomicCounter = 123
	lexid.ServerUniqueId = `2`
	L.Print(`MaxTime`, time.Unix(0, math.MaxInt64))
	id := lexid.NanoID()
	print(`id`, id)
	seg, err := lexid.Parse(id)
	L.IsError(err, `failed parse id`)
	L.Print(`Time`, time.Unix(0, seg.Time))
	L.Print(`Counter`, seg.Counter)
	assert.Equal(t, lexid.AtomicCounter, seg.Counter)
	L.Print(`ServerID`, seg.ServerID)
	assert.Equal(t, lexid.ServerUniqueId, seg.ServerID)
}

func TestParseFixedOO(t *testing.T) {
	gen := lexid.NewGenerator(`1`)
	gen.Separator = ``
	gen.AtomicCounter = 123
	gen.ServerUniqueId = `2`
	L.Print(`MaxTime`, time.Unix(math.MaxInt64, 0))
	id := gen.ID()
	print(`id`, id)
	seg, err := gen.Parse(id)
	L.IsError(err, `failed parse id`)
	L.Print(`Time`, time.Unix(seg.Time, 0))
	assert.Equal(t, gen.AtomicCounter, seg.Counter)
	L.Print(`ServerID`, seg.ServerID)
	assert.Equal(t, gen.ServerUniqueId, seg.ServerID)
}

func TestParseFixed(t *testing.T) {
	defer func() {
		// restore configuration
		lexid.Separator = `~`
		lexid.ServerUniqueId = `~0`
	}()
	lexid.Separator = ``
	lexid.AtomicCounter = 123
	lexid.ServerUniqueId = `2`
	L.Print(`MaxTime`, time.Unix(math.MaxInt64, 0))
	id := lexid.ID()
	print(`id`, id)
	seg, err := lexid.Parse(id)
	L.IsError(err, `failed parse id`)
	L.Print(`Time`, time.Unix(seg.Time, 0))
	assert.Equal(t, lexid.AtomicCounter, seg.Counter)
	L.Print(`ServerID`, seg.ServerID)
	assert.Equal(t, lexid.ServerUniqueId, seg.ServerID)
}
