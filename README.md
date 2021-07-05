
# LexID: fast lexicographically orderable ID

Can generate ~10 millions id per second (single core only).

Consist of 3 segment:
- Unix or UnixNano (current time. 6/11 character)
- Atomic Counter (limit to single core, 6 character)
- Server Unique ID (or process or thread ID, min. 0 character)
- 2 separator character (default: `~`, can be removed, 2x 0-1 character)

```
Min length (ID without separator and server ID): 6+6+0+0 = 12 bytes
Max length (NanoID with separator and server ID): 11+6+N+2 = 19+N bytes
``` 

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
BenchmarkLexIdNoSep-8  	10020199	   118.1 ns/op
BenchmarkLexIdNoLex-8  	10475790	   116.3 ns/op
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
	
	// optional separator, 
	// you can set this to empty string if you set the Min*TimeLength >= 6 or 11
	lexid.Separator = `~`
	
	// optional minimum counter segment length, 
	// if set lower than 6 will not lexicographically orderable anymore
	lexid.MinCounterLength = 0
	
	// optional minimum time segment length, default: 6
	// if set lower than 6 might not lexicographically orderable anymore
	lexid.MinTimeLength = 0
	
	// optional minimum nano time segment length, default: 11
	// if set lower than 11 might not lexicographically orderable anymore
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

this shows minimum length and length after 10 million generated id with specific configuration (8-15 characters for `ID`, 15-20 characters for `NanoID`)

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


ID Separator=`` MinTimeLength=6
first 0Vsccp-----00
 len= 13
last 0Vsccp--a8P00
 len= 13

NanoID Separator=`` MinNanoTimeLength=11
first 0PDmclT1CmN-----00
 len= 18
last 0PDmclT1CmN--a8P00
 len= 18
 
 
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
- the `AtomicCounter` is overflowed on the exact same second/nanosecond, you might want to reset the counter every >1 second to overcome this (or you might want to ignore this if ordering doesn't matter if the event happened on the same second/nanodescond)
- you change `Separator` to other character that have lower ASCII/UTF-8 encoding value.
- you set `Min*Length` too low, it should be `>=6` for `MinTimeLength` and `>=11` for `MinNanoTimeLength`, and `6` for `MinCounterLength`
- the `time` segment already pass the `MinTimeLength`, earliest will happen at year 4147.

it might duplicate if:
- your processor so powerful, that it can call the function above faster than 4 billion (`MaxUint32`=`4,294,967,295`) per second, there's no workaround for this.
- you set the `AtomicCounter` multiple time on the same second/nanosecond (eg. to a number lower than current counter)
- using same/shared server/process/thread id on different server/process/thread 
- unsynchronized time on same server

it will impossible to parse (to get time, counter, and server id) if:
- you set `Separator` to empty string and all other `Min*Length` to lower than recommended value 


## Difference with XID

- have locks, so by default only utilize single core performance (unless using the object-oriented version to spawn multiple instance with different server/process/thread id)
- 256x more uniqueness generated id per sec guaranteed: 4 billion vs 16 million
- configurable (length, separator, server/process/thread id)
- same or less length for string representation (depends on your configuration)
- base63 vs base32
- defaults to string representation (12 to 19+N bytes) vs have 12-bytes binary representation
