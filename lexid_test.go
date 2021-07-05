package lexid_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/lexid"

	"github.com/google/uuid"
	"github.com/kokizzu/gotro/L"
	"github.com/matoous/go-nanoid/v2"
)

func BenchmarkUuid(b *testing.B) {
	for z := 0; z < b.N; z++ {
		res := uuid.New()
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

func TestLexiId(t *testing.T) {
	const N = 10_000_000
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
	const N = 10_000_000
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
	const N = 10_000_000
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
	const N = 10_000_000
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
	print(`recommended minimum length`, S.EncodeCB63(lexid.Now.Unix(), 0))
	print(`recommended minimum length`, S.EncodeCB63(lexid.Now.UnixNano(), 0))
}

func TestObject(t *testing.T) {
	const N = 10_000_000
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
	const N = 10_000_000
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
