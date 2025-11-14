# KVBench

---

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

---

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

---

## benchmark

* The following benchmarks show the throughput of inserting/reading keys (of size 9 bytes) and values (of size 256 bytes). 
* Batch write cost is the time it takes to write 4,000,000 keys and values.

### SSD

Computer configuration: Apple M4 Pro, 14（10性能和4能效）, 48 GB RAM, 500GB SSD

#### nofsync

**throughputs**

| name | batch write cost(s) | MemUsage(MiB) | HeapInuse(MiB) | DiskUsage(MiB) | Prefix op/s | Set op/s | Get op/s | Setmixed op/s | Getmixed op/s | Del op/s |
|------|---------------------|---------------|----------------|----------------|-----------|----------|----------|---------------|---------------|----------|
| badger  | 6   | 1303 | 1309 | 3654 | 373046  | 275498 | 731472  | 9791  | 653621 | 301646  |
| rocksdb | 6   | 3    | 4    | 1310 | 390185  | 197135 | 294849  | 91438 | 234092 | 207075  |
| buntdb  | 8   | 1743 | 1748 | 1269 | 1053    | 115869 | 6036609 | 19891 | 305514 | 1202369 |
| nutsdb  | 8   | 2009 | 2025 | 1280 | 2273910 | 270784 | 2724855 | 51480 | 784030 | 1384351 |
| leveldb | 18  | 21   | 23   | 1061 | 52545   | 156204 | 488964  | 94729 | 158719 | 513809  |
| rosedb  | 25  | 689  | 695  | 1091 | 5       | 260050 | 5659095 | 58831 | 906305 | 306553  |
| pebble  | 42  | 4    | 6    | 1062 | 103175  | 138720 | 105370  | 99093 | 37688  | 479873  |
| bbolt   | 105 | 55   | 56   | 1571 | 646685  | 31412  | 790603  | 12358 | 572332 | 116347  |


**Index ranking**

The higher the ranking, the better

| Rank | BatchWrite | MemUsage | DiskUsage | Prefix | Set | Get | Setmixed | Getmixed | Delete |
|:----:|:-----------|:---------|:----------|:-------|:---|:---|:--------|:--------|:------|
| 1    | badger     | rocksdb  | leveldb   | nutsdb | badger  | buntdb  | pebble  | rosedb  | nutsdb  |
| 2    | rocksdb    | pebble   | pebble    | bbolt  | nutsdb  | rosedb  | leveldb | nutsdb  | buntdb  |
| 3    | buntdb     | leveldb  | rosedb    | rocksdb| rosedb  | nutsdb  | rocksdb | badger  | leveldb |
| 4    | nutsdb     | bbolt    | buntdb    | badger | rocksdb | bbolt   | rosedb  | bbolt   | pebble  |
| 5    | leveldb    | rosedb   | nutsdb    | pebble | leveldb | badger  | nutsdb  | buntdb  | rosedb  |
| 6    | rosedb     | badger   | rocksdb   | leveldb| pebble  | leveldb | buntdb  | rocksdb | badger  |
| 7    | pebble     | buntdb   | bbolt     | buntdb | buntdb  | rocksdb | bbolt   | leveldb | rocksdb |
| 8    | bbolt      | nutsdb   | badger    | rosedb | bbolt   | pebble  | badger  | pebble  | bbolt   |

* pogreb and sniper does not support prefix queries

#### fsync

Coming soon...


---

### HDD

Computer configuration: 

Coming soon...
