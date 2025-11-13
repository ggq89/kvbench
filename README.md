# KVBench

Cloned from [liulhdarks/go-kvbench](https://github.com/liulhdarks/go-kvbench). Compared to the [smallnest/kvbench](https://github.com/smallnest/kvbench) codebase:
1. Fixed some incorrect logic of KV database prefix query.
2. Fixed an issue where some KV database configurations were incorrect and persistence was not enabled
3. Batch writing is changed to write a fixed amount of data instead of a fixed time, which can make the subsequent query evaluation fairer.
4. Add the evaluation of prefix query.
5. Add the evaluation of memory usage and disk usage.

KVBench is a Redis server clone backed by a few different Go databases. 

It's intended to be used with the `redis-benchmark` command to test the
performance of various Go databases.  It has support for redis pipelining. The
`redis-benchmark` can run as explained here https://github.com/tidwall/kvbench#examples.

This cloned version adds more kv databases and automatic scripts.

Features:

- Databases
  - [badger](https://github.com/dgraph-io/badger)
  - [BboltDB](https://github.com/etcd-io/bbolt)
  - [BoltDB](https://github.com/boltdb/bolt)
  - [buntdb](https://github.com/tidwall/buntdb)
  - [LevelDB](https://github.com/syndtr/goleveldb)
  - [modernc.org/kv](https://gitlab.com/cznic/kv)
  - [rocksdb](https://github.com/linxGnu/grocksdb)
  - [pebble](https://github.com/cockroachdb/pebble)
  - [pogreb](https://github.com/akrylysov/pogreb)
  - [nutsdb](https://github.com/nutsdb/nutsdb)
  - [sniper](https://github.com/recoilme/sniper)
  - btree (in-memory) with [AOF persistence](https://redis.io/topics/persistence)
  - map (in-memory) with [AOF persistence](https://redis.io/topics/persistence)
  - [bitcask](https://git.mills.io/prologic/bitcask)
  - [rosedb](https://github.com/rosedblabs/rosedb)
- Option to disable fsync
- Compatible with Redis clients

## Quickstart

Run the following commands:
```shell
cd cmd/cli
go build -o cli main.go 
./test.sh
```

Or manual test cli command:
```shell
Usage of ./cli:
  -c int
        concurrent goroutines (default runtime.NumCPU())
  -d duration
        test duration for each case (default 10s)
  -fsync
        fsync (default false)
  -s string
        store type (default "map")
  -save string
        save path, ouput csv file path (default "", not output)
  -set int
        batch set count (default 4000000)
  -size int
        data size for each value (default 256)
```

Example:
```shell
./cli -d 10s -size 256 -s "bbolt" -save "benchmarks/nofsync.csv" >> benchmarks/test.log 2>&1
```

## SSD benchmark
The following benchmarks show the throughput of inserting/reading keys (of size
9 bytes) and values (of size 256 bytes). Batch write cost is the time it takes to write 4,000,000 keys and values.

Computer configuration: Apple M4 Pro, 14（10性能和4能效）, 48 GB RAM, 500GB SSD

### nofsync

**throughputs**

| name | batch write cost(s) | MemUsage(MiB) | HeapInuse(MiB) | DiskUsage(MiB) | Prefix op/s | Set op/s | Get op/s | Setmixed op/s | Getmixed op/s | Del op/s |
|------|---------------------|---------------|----------------|----------------|-----------|----------|----------|---------------|---------------|----------|
| map          | 2   | 1672 | 1673 | 1113 | 71      | 402469  | 5996264  | 84967   | 1573401  | 2884469  |
| map/memory   | 2   | 1928 | 1929 | 75   | 2057561 | 7177591 | 102080   | 2002329 | 2270046  |          |
| badger       | 6   | 909  | 912  | 3335 | 401682  | 284798  | 740151   | 9442    | 654084   | 310713   |
| rocksdb      | 6   | 3    | 4    | 1291 | 361015  | 199367  | 299135   | 90945   | 241338   | 203154   |
| btree        | 6   | 1532 | 1536 | 1113 | 1751    | 382492  | 6233683  | 65373   | 1021511  | 2003842  |
| btree/memory | 6   | 1451 | 1463 | 1576 | 1413167 | 6210927 | 63676    | 1079349 | 1193747  |          |
| buntdb       | 7   | 1680 | 1685 | 1279 | 804     | 116168  | 6017012  | 22275   | 344268   | 1199046  |
| nutsdb       | 7   | 1848 | 1865 | 1280 | 2281420 | 265854  | 2456997  | 53381   | 812805   | 1234043  |
| sniper       | 8   | 149  | 159  | 1953 | -1      | 550118  | 35392977 | 267492  | 29127345 | 31385085 |
| leveldb      | 17  | 15   | 17   | 1068 | 148753  | 153399  | 665573   | 67210   | 171993   | 500855   |
| bitcask      | 23  | 2981 | 3014 | 1071 | 624214  | 75022   | 1020881  | 30134   | 519674   | 307039   |
| rosedb       | 28  | 602  | 608  | 1091 | 6       | 294411  | 5630722  | 60743   | 922339   | 299996   |
| pogreb       | 30  | 2    | 3    | 1154 | -1      | 127821  | 6019969  | 60578   | 875630   | 2456936  |
| pebble       | 41  | 3    | 5    | 1055 | 103501  | 141950  | 94970    | 99085   | 34565    | 463641   |
| bbolt        | 104 | 38   | 40   | 1574 | 653242  | 30107   | 783146   | 12522   | 570597   | 113127   |
| bolt         | 107 | 35   | 36   | 1582 | 774438  | 30352   | 1010986  | 9791    | 644485   | 30554    |



**Index ranking**

The higher the ranking, the better

| Rank | BatchWrite | MemUsage | DiskUsage | Prefix | Set | Get | Setmixed | Getmixed | Delete |
|:----:|:-----------|:---------|:----------|:-------|:---|:---|:--------|:--------|:------|
| 1    | rocksdb    | rocksdb  | pebble    | nutsdb  | sniper  | sniper  | sniper  | sniper  | sniper |
| 2    | badger     | pogreb   | leveldb   | bolt    | rosedb  | pogreb  | pebble  | rosedb  | pogreb |
| 3    | nutsdb     | pebble   | bitcask   | bbolt   | badger  | buntdb  | rocksdb | pogreb  | nutsdb |
| 4    | buntdb     | leveldb  | rosedb    | bitcask | nutsdb  | rosedb  | leveldb | nutsdb  | buntdb |
| 5    | sniper     | bolt     | pogreb    | badger  | rocksdb | nutsdb  | rosedb  | badger  | leveldb |
| 6    | leveldb    | bbolt    | buntdb    | rocksdb | leveldb | bitcask | pogreb  | bolt    | pebble |
| 7    | bitcask    | sniper   | nutsdb    | leveldb | pebble  | bolt    | nutsdb  | bbolt   | badger  |
| 8    | rosedb     | rosedb   | rocksdb   | pebble  | pogreb  | bbolt   | bitcask | bitcask | bitcask  |

* pogreb and sniper does not support prefix queries

### fsync

**throughputs**

Coming soon...
