
# LexID: fast lexicographically orderable ID

Can generate ~10 millions id per second (single core only).

Consist of 3 segment:
- Unix or UnixNano (current time)
- Atomic Counter (limit to single core)
- Server Unique ID (or process or thread ID)
- 2 separator character (default: `~`)

Based on [lexicographically sortable encoding](//github.com/kokizzu/gotro/tree/master/S), URL-safe encoding.

```
cpu: AMD Ryzen 3 3100 4-Core Processor    
BenchmarkShortuuid-8   	  139238	  8461 ns/op
BenchmarkKsuid-8       	  749536	  1497 ns/op
BenchmarkNanoid-8      	  775210	  1490 ns/op
BenchmarkUuid-8        	  875721	  1340 ns/op
BenchmarkTime-8        	 1674458	   712.9 ns/op
BenchmarkSnowflake-8   	 4909974	   244.8 ns/op
BenchmarkLexIdNano-8   	 7718455	   142.0 ns/op
BenchmarkLexId-8       	 9906074	   118.8 ns/op
BenchmarkXid-8         	13322355	    86.91 ns/op
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
	
	// parse to get time, counter, and server id
	seg, err := lexid.Parse(id)
	seg, err = lexid.Parse(nanoid)  
	
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
- the `AtomicCounter` is overflowed on the exact same second/nanosecond, you might want to reset the counter every >1 second to overcome this (or you might want to ignore this if order doesn't matter if it's the event happened on the same second/nanodescond)
- the `time` segment already pass current length, you might want to set `MinTimeLength` to `>6` and `MinNanoTimeLength` to `>11` 
- you change `Separator` to other character that have lower ASCII/UTF-8 encoding value

it might duplicate if:
- your processor can call the function above, faster than 4 billion (`MaxUint32`=`4,294,967,295`) per second, there's no workaround for this.
- you set the `AtomicCounter` multiple time on the same second/nanosecond (eg. to a number lower than current counter)

## Difference with XID

- have locks, so by default only utilize single core performance (unless using OO version to spawn multiple instance with different server/process/thread id)
- 256x more uniqueness generated id per sec guaranteed: 4 billion vs 16 million
- configurable (length, separator, server/process/thread id)
- same or less length for string representation (depends on your configuration)
- base63 vs base32
- defaults to string representation vs have 12-bytes binary representation
