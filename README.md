
# LexID: Fast lexicographically orderable/sortable ID

Highly customizable ID generator that can generate ~10 millions IDs per second (single core only).

Consist of 3 segment:
- Unix or UnixNano (current time, `6`/`11` character)
- Atomic Counter (limit to single core, `6` character)
- Server/Process/Thread Unique `Identity` (optional, min. `0` character, default: `~0`)
- 1 or 2 separator character (default: `~`, can be removed, 2x `0`-`1` character)

Default formats:
- `ID()`: `tttttt~cccccc~i`
- `NanoID()`: `ttttttttttt~cccccc~i`
- `t` = timestamp
- `c` = counter
- `i` = identity, can be more than 1 character

Based on [lexicographically sortable encoding](//github.com/kokizzu/gotro/tree/master/S), URL-safe encoding.

## Configuration Comparison

| Type   | Min[Nano]<br/>Date<br/>Offset | Min[Nano]<br/>Time<br/>Length | Min<br/>Counter<br/>Length | Separator | Byte use<br/>without<br/>Identity | Ordered | Unique |
|:------:|------------------------------:|------------------------------:|---------------------------:|:---------:|----------------------------------:|:-------:|:------:|
| ID     | 0                             | 6                             | 6                          | ~         | 13                                | Y       | Y      |
| ID     | 0                             | 0                             | 0                          | ~         | 8-13                              | N       | Y      |
| ID     | 0                             | 6                             | 6                          |           | 12                                | Y       | Y      |
| ID     | 0                             | 0                             | 0                          |           | 6-12                              | N       | N      |
| ID     | 1577836800                    | 0                             | 0                          | ~         | 7                                 | N       | Y      |
| ID     | 1577836800                    | 0                             | 0                          |           | 6                                 | N       | N      |
| NanoID | 0                             | 11                            | 6                          | ~         | 19                                | Y       | Y      |
| NanoID | 0                             | 0                             | 0                          | ~         | 14-19                             | N       | Y      |
| NanoID | 0                             | 11                            | 6                          |           | 17                                | Y       | Y      |
| NanoID | 0                             | 0                             | 0                          |           | 12-17                             | N       | N      |
| NanoID | 1577836800<br>000000000       | 0                             | 0                          | ~         | 12                                | N       | Y      |
| NanoID | 1577836800<br>000000000       | 0                             | 0                          |           | 11                                | N       | N      |

Note: `1577836800` = unix timestamp of `2021-01-01 00:00:00`

Uniqueness configuration (when `Separator` or `Min*TimeLength` set, this is the default)
```
Min length (ID with separator and server identity): 
  6+6+0+1 = 13 bytes (format: `tttttt~cccccc`)
  5+1+0+1 = 7 bytes (with 2020 offset, format: `ttttt`~`c`)
Max length (NanoID with separator and server identity): 
  11+6+N+2 = 19+N bytes (format: `ttttttttttt~cccccc~i`)
  10+1+N+1 = 12+N bytes (with 2020 offset, format: `tttttttttt~c~i`)
``` 

Ordered/sortable configuration (when `Min*TimeLength` set, may unset the `Separator`)
```
Min length (ID without separator and server identity): 
  6+6+0+0 = 12 bytes (format: `ttttttcccccc`)
Max length (NanoID without separator and with server identity): 
  11+6+N+0 = 17+N bytes (format: `tttttttttttcccccci`)
```

## Benchmark

```
cpu: AMD Ryzen 3 3100 4-Core Processor    
BenchmarkShortuuid-8      118908      8572 ns/op
BenchmarkKsuid-8          760924      1493 ns/op
BenchmarkNanoid-8         759548      1485 ns/op
BenchmarkUuid-8           935152      1304 ns/op
BenchmarkTime-8          1690483       720.0 ns/op
BenchmarkSnowflake-8     4911249       244.7 ns/op
BenchmarkLexIdNano-8     8483720       138.8 ns/op <--
BenchmarkLexIdNoSep-8   10396551       116.3 ns/op <--
BenchmarkLexIdNoLex-8   10590300       115.1 ns/op <--
BenchmarkLexId-8         9991906       114.9 ns/op <--
BenchmarkXid-8          13754178        86.02 ns/op
BenchmarkId64-8        276799974         4.362 ns/op
```

## Usage

```
import "github.com/kokizzu/lexid"

func main() {
	// set if multiserver, can be empty if not multi-server/process/thread/instance
	lexid.Config.Identity = `` // default: ~0
	
	// optional starting counter
	lexid.Config.AtomicCounter = 0
	
	// optional segment separator
	// you can set this to empty string if you keep the Min*Length as default 
	lexid.Config.Separator = `` // default: ~
	
	// optional minimum counter segment length, default: 6
	// if set lower than 6 will not lexicographically orderable anymore
	lexid.Config.MinCounterLength = 0
	
	// optional minimum time segment length, default: 6
	// if set lower than 6 might not lexicographically orderable anymore
	lexid.Config.MinTimeLength = 0
	
	// optional minimum nano time segment length, default: 11
	// if set lower than 11 might not lexicographically orderable anymore
	lexid.Config.MinNanoTimeLength = 0
	
	// optional date offset, can reduce length of the time segment
	lexid.Config.MinDateOffset = lexid.OffsetY2020.Unix()
	lexid.Config.MinNanoDateOffset = lexid.OffsetY2020.UnixNano()
	
	// smaller id, second resolution
	id := lexid.ID()
	
	// larger id, nanosecond resolution (`~5 ms` to be exact)
	nanoid := lexid.NanoID()
	
	// parse to get time, counter, and server id
	seg, err := lexid.Parse(id)
	seg, err = lexid.Parse(nanoid)  
	
	// generate id from segment/component
	id = lexid.FromUnix(time)
	id = lexid.FromUnixCounter(time,counter)
	id = lexid.FromUnixCounterIdent(time,counter,identity)
	id = lexid.FromNano(time)
	id = lexid.FromNanoCounter(time,counter)
	id = lexid.FromNanoCounterIdent(time,counter,identity)
	
	// object-oriented version, 
	// eg. if you need to generate uniquely one per core/thread
	//     or when each database table need different ID format
	gen := lexid.NewGenerator(`~1`) // ~2 for 2nd core, ~3 for 3rd core, and so on
	// gen.Identity = `~1`
	// gen.Separator = `~`
	// gen.AtomicCounter = 0
	// gen.MinCounterLength = 6
	// gen.MinTimeLength = 6
	// gen.MinNanoTimeLength = 11
	// gen.MinDateOffset = 0 
	// gen.MinNanoDateOffset = 0 
	// gen.From*() also exists
	
	id = gen.ID()
	nanoid = gen.NanoID()
}
```

## Example generated ID

shows minimum length and length after 1-10 million generated ID with specific configuration (6-15 characters for `ID`, 11-20 characters for `NanoID`). 

Default config (fixed length):
```
ID 
first: 0Vsccp~-----0~0
 len= 15
last: 0Vsccp~--a8P0~0
 len= 15

NanoID
first: 0PDmclT1CmN~-----0~0
 len= 20
last: 0PDmclT1CmN~--a8P0~0
 len= 20
```
 
Separatorless config and without Identity:
```
ID Separator=`` Identity=`` MinTimeLength=6 (default)
first: 0Vsccp-----0
 len= 12
last: 0Vsccp--a8P0
 len= 12

NanoID Separator=`` Identity=`` MinNanoTimeLength=11 (default)
first: 0PDmclT1CmN-----0
 len= 17
last: 0PDmclT1CmN--a8P0
 len= 17
```

Variable length config (not lexicographically sortable):
```
ID MinCounterLength=0
first: 0Vsc0a~0~0 
 len= 10
last:  0Vsc0a~2o80~0 
 len= 13

NanoID MinCounterLength=0
first 0PDm7hn0KSs~0~0
 len= 15
last 0PDm7hn0KSs~2o80~0
 len= 18 
```
 
Allows duplicate config:
```
ID Separator=`` Identity=`` MinCounterLength=0
first: 0Vsccp0
 len= 7
last: 0Vsccpa8P0
 len= 10

NanoID Separator=`` Identity=`` MinCounterLength=0
first: 0PDmclT1CmN0
 len= 12
last: 0PDmclT1CmNa8P0
 len= 15
``` 

Offsetted config (reduce time segment by 2020-01-01):
```
ID MinTimeLength=0 MinCounterLength=0
example: 1pkHb~0~0
 len= 9
 
NanoID MinNanoTimeLength=0 MinCounterLength=0
example: 1dGr84ixhV~1~0
 len= 14
``` 
 
Offsetted with minimum length and allows duplicate:
```
ID MinTimeLength=0 MinCounterLength=0 Separator=`` Identity=``
example: 1pkHb0
 len= 6
 
NanoID MinNanoTimeLength=0 MinCounterLength=0 Separator=`` Identity=``
example: 1dGr84ixhV1
 len= 11
```

## Gotchas

it might not lexicographically ordered/sorted if:
- the `AtomicCounter` is overflowed on the exact same second, can be happening when your processor able to call `ID()` function more than 4 billion (`MaxUint32`=`4,294,967,295`) times per second.
- you change `Separator` to other character that have lower value than `z` (`122`) of ASCII/UTF-8, the default is `~` (`126`).
- you set `Min*Length` less than recommended value, it should be `>=6` for `MinTimeLength` and `>=11` for `MinNanoTimeLength`, and `6` for `MinCounterLength`.
- you unset `Separator` and set `MinCounterLength` lower than `6`
- the `time` segment already pass the `MinTimeLength=6`, earliest will happen at year `4147`.
- you change system time to earlier time.

it might duplicate if:
- your processor too powerful, that it can call the function `ID()` more than 4 billion times per second, workaround: use `NanoID()`.
- you set the `AtomicCounter` multiple time on the same second/nanosecond (eg. to a number lower than current counter).
- using same/shared `Identity` on different server/process/thread.
- unsynchronized time on same server.
- you set `Min*DateOffset` too low or differently on each run.
- you change `Separator` to empty string or characters that are in `EncodeCB63` with `Min*Length` less than recommended value.
- you change system time to earlier time.

it will impossible to parse (to get time, counter, and server id) if:
- you unset `Separator` and `Min*Length` to lower than recommended value.

hidden logic/implementation detail?
- counter resets to 0 when second changed but only when `ID()` called, so it would never overflow (unless your processor too powerful).
- nanosecond incremented by 1 every 4 billion `NanoID()` calls, so it would never collide.

## Difference with XID

- have locks, so by default can only utilize single core (unless using the object-oriented version to spawn multiple instance with different server/process/thread `Identity`)
- 256x more max unique IDs per sec guaranteed: 4 billion vs 16 million
- configurable (length, separator, server/process/thread identity, date offset)
- EncodeCB63 (base64-variant) vs base32 (20% bigger)
- defaults to string representation (12 to 17+N bytes) vs have 12-bytes binary representation (20 bytes for string representation)
- server/process/thread `Identity` are optional
- can be offsetted (subtracted with certain value, eg. `2020-01-01 00:00:00`)


## See also

[id64](//github.com/kokizzu/id64) - quick non-distributed uint64 id generator (can generate ~276 million ID/second)
