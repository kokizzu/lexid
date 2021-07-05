
# LexID: fast lexicographically orderable ID

Can generate ~10 millions id per second (single core only).

Consist of 3 segment:
- Unix or UnixNano
- Atomic Counter (limit to single core)
- Server Unique ID (or thread ID)

Based on [lexicographically sortable encoding](//github.com/kokizzu/gotro/S)

```
cpu: AMD Ryzen 3 3100 4-Core Processor            
BenchmarkUuid
BenchmarkUuid-8        	  992215	      1213 ns/op
BenchmarkNanoid
BenchmarkNanoid-8      	  754716	      1487 ns/op
BenchmarkLexId
BenchmarkLexId-8       	10611122	       114.2 ns/op
BenchmarkLexIdNano
BenchmarkLexIdNano-8   	 8493916	       138.2 ns/op
PASS
```

## Usage

```
import "github.com/kokizzu/lexid"

func main() {
	// set if multiserver
	lexid.UniqueServerId = `~1`
	
	// optional starting counter
	lexid.AtomicCounter = 0
	
	// optional separator
	lexid.Separator = `~`
	
	// optional minimum counter segment length, 
	// if set too low will not lexicographically orderable anymore
	lexid.MinCounterLength = 0
	
	// optional minimum time segment length
	lexid.MinTimeLength = 0
	
	// optional minimum nano time segment length
	lexid.MinNanoTimeLength = 0
	
	// smaller id, second resolution
	id := lexid.ID()
	
	// larger id, nanosecond resolution
	nanoid := lexid.NanoID()
	
	// object-oriented version, eg. if you need to generate uniquely one per core
	gen := lexid.NewGenerator(`~1`) // ~2 for 2nd core, ~3 for 3rd core, and so on
	
	id = gen.ID()
	nanoid = gen.NanoID()
}
```

## Example generated id

this shows minimum length and length after 10 million generated id with specific configuration (10-15 characters for `ID`, 15-20 characters for `NanoID`)

```
ID 
first 0Vsccp~-----0~0
 len= 15
last 0Vsccp~--a8P0~0
 len= 15

NanoID
first 0PDmclT1CmN~-----0~0
 len= 20
last 0PDmclT1CmN~--a8P0~0
 len= 20


ID MinCounterLength = 0
first 0Vsc0a~0~0 
 len=10
last  0Vsc0a~2o80~0 
 len=13

NanoID MinCounterLength = 0
first 0PDm7hn0KSs~0~0
 len= 15
last 0PDm7hn0KSs~2o80~0
 len= 18 
```

## Gotchas

it might not lexicographically ordered if:
- the `AtomicCounter` is overflow, you might want to reset the counter every >1 second for example
- the `time` already pass current length, you might want to set `MinTimeLength` to `>6` and `MinNanoTimeLength` to `>11` 
- you change `Separator` to other character that have lower ASCII/UTF-8 encoding value

it might duplicate if:
- your processor can call the function above, faster than 4 billion (`MaxUint32`=`4,294,967,295`) per second, there's no workaround for this.
- you set the `AtomicCounter` multiple time on the same second/nanosecond (eg. to a number lower than current counter)
