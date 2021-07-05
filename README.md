
# LexID: fast lexicographically orderable ID

Can generate ~10 millions id per second (single core only).

Consist of 3 segment:
- Unix or UnixNano (current time. 6/11 character)
- Atomic Counter (limit to single core, 6 character)
- Server Unique ID (or process or thread ID, min. 0 character)
- 2 separator character (default: `~`, can be removed, 2x 0-1 character)

```
Min length (ID without separator and server ID): 
  6+6+0+0 = 12 bytes
Max length (NanoID with separator and server ID): 
  11+6+N+2 = 19+N bytes
``` 
|  Type  | Min[Nano]<br/>Time<br/>Length | Min<br/>Counter<br/>Length | Separator | Byte use<br/>without<br/>ServerID | Ordered | Unique |
|:------:|:---------------------:|:------------------:|:---------:|:--------------------------:|:-------:|:------:|
| ID     |                     6 |                  6 |     ~     |             14             |    Y    |    Y   |
| ID     |                     0 |                  0 |     ~     |            9-14            |    N    |    Y   |
| ID     |                     6 |                  6 |           |             12             |    Y    |    Y   |
| ID     |                     0 |                  0 |           |            7-12            |    N    |    N   |
| NanoID |                    11 |                  6 |     ~     |             19             |    Y    |    Y   |
| NanoID |                     0 |                  0 |     ~     |            14-19           |    N    |    Y   |
| NanoID |                    11 |                  6 |           |             17             |    Y    |    Y   |
| NanoID |                     0 |                  0 |           |            12-17           |    N    |    N   |

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
	// set if multiserver, can be empty if not multi-server
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

## Example generated ID

this shows minimum length and length after 1-10 million generated ID with specific configuration (7-15 characters for `ID`, 10-20 characters for `NanoID`)

```
Default config (fixed length):

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
 
 
Separatorless config and without ServerUniqueID:

ID Separator=`` ServerUniqueId=`` MinTimeLength=6 (default)
first 0Vsccp-----0
 len= 12
last 0Vsccp--a8P0
 len= 12

NanoID Separator=`` ServerUniqueId=`` MinNanoTimeLength=11 (default)
first 0PDmclT1CmN-----0
 len= 17
last 0PDmclT1CmN--a8P0
 len= 17
 

Config with variable length (not lexicographically sortable):
 
ID MinCounterLength=0
first 0Vsc0a~0~0 
 len= 10
last  0Vsc0a~2o80~0 
 len= 13

NanoID MinCounterLength=0
first 0PDm7hn0KSs~0~0
 len= 15
last 0PDm7hn0KSs~2o80~0
 len= 18 
 
 
Config that allows duplicate:

ID Separator=`` ServerUniqueId=`` MinCounterLength=0
first 0Vsccp0
 len= 7
last 0Vsccpa8P0
 len= 10

NanoID Separator=`` ServerUniqueId=`` MinCounterLength=0
first 0PDmclT1CmN0
 len= 12
last 0PDmclT1CmN~a8P0
 len= 16
```

## Gotchas

it might not lexicographically ordered if:
- the `AtomicCounter` is overflowed on the exact same second/nanosecond, you might want to reset the counter every >1 second to overcome this (or you might want to ignore this if ordering doesn't matter if the event happened on the same second/nanodescond).
- you change `Separator` to other character that have lower ASCII/UTF-8 encoding value.
- you set `Min*Length` less than recommended value, it should be `>=6` for `MinTimeLength` and `>=11` for `MinNanoTimeLength`, and `6` for `MinCounterLength`.
- the `time` segment already pass the `MinTimeLength`, earliest will happen at year 4147.

it might duplicate if:
- your processor so powerful, that it can call the function above faster than 4 billion (`MaxUint32`=`4,294,967,295`) per second, there's no workaround for this.
- you set the `AtomicCounter` multiple time on the same second/nanosecond (eg. to a number lower than current counter).
- using same/shared `ServerUniqueID` on different server/process/thread.
- unsynchronized time on same server.
- you change `Separator` to empty string or characters that are in `EncodeCB63` with `Min*Length` less than recommended value.

it will impossible to parse (to get time, counter, and server id) if:
- you set `Separator` to empty string and all other `Min*Length` to lower than recommended value.

## Difference with XID

- have locks, so by default can only utilize single core (unless using the object-oriented version to spawn multiple instance with different server/process/thread id)
- 256x more uniqueness generated id per sec guaranteed: 4 billion vs 16 million
- configurable (length, separator, server/process/thread id)
- same or less length for string representation (depends on your configuration)
- EncodeCB63 (base64-variant) vs base32 (20% space usage)
- defaults to string representation (12 to 19+N bytes) vs have 12-bytes binary representation

## Further Optimization
- add `MinDate` (eg. `20xx-01-01 00:00:00`) to subtract the generated `Time` segment so it would use less space.
