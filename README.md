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
  - [lmdb](https://github.com/Data-Corruption/lmdb-go)
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
| map          | 2   | 1858 | 1859 | 1113    | 4081720 | 431388  | 6086506  | 110670  | 2339050  | 2841354  |
| map/memory   | 2   | 1754 | 1755 | 4070366 | 2114611 | 7337892 | 144778   | 3585491 | 2316429  |          |
| rocksdb      | 4   | 0    | 1    | 1309    | 366805  | 194705  | 424482   | 93561   | 256414   | 205428   |
| badger       | 6   | 1870 | 1877 | 3335    | 373883  | 260851  | 778606   | 8221    | 704142   | 291006   |
| btree        | 6   | 1579 | 1583 | 1113    | 2490422 | 353988  | 5771065  | 71163   | 1110697  | 2194178  |
| btree/memory | 6   | 1363 | 1395 | 2497127 | 1478243 | 6240280 | 80139    | 1430160 | 1208952  |          |
| nutsdb       | 7   | 1788 | 1808 | 1280    | 2280523 | 270797  | 2499417  | 56107   | 862956   | 1293521  |
| buntdb       | 8   | 1484 | 1498 | 1260    | 989     | 107842  | 6517283  | 22772   | 350649   | 1175236  |
| sniper       | 9   | 233  | 241  | 1953    | -1      | 544316  | 51672884 | 192979  | 37113968 | 32968963 |
| leveldb      | 18  | 19   | 20   | 1064    | 139123  | 165550  | 801371   | 39425   | 526989   | 516589   |
| bitcask      | 24  | 2254 | 2314 | 1071    | 638231  | 73913   | 1036781  | 29959   | 500502   | 316459   |
| pogreb       | 29  | 2    | 3    | 1154    | -1      | 122380  | 5858526  | 58014   | 836904   | 2361001  |
| rosedb       | 29  | 712  | 718  | 1091    | 5       | 299429  | 5617256  | 63427   | 967331   | 301889   |
| pebble       | 46  | 3    | 5    | 1060    | 94758   | 143077  | 126942   | 104595  | 36713    | 497503   |
| bolt         | 104 | 37   | 38   | 1581    | 704869  | 31513   | 903939   | 10138   | 632362   | 33263    |
| bbolt        | 108 | 39   | 41   | 1572    | 642533  | 31814   | 778625   | 12299   | 563267   | 108132   |


**Index ranking**

The higher the ranking, the better

| Rank | BatchWrite | MemUsage | DiskUsage | Prefix | Set | Get | Setmixed | Getmixed | Delete |
|:----:|:-----------|:---------|:----------|:-------|:---|:---|:--------|:--------|:------|
| 1    | rocksdb    | rocksdb  | pebble    | nutsdb  | sniper | sniper | sniper | sniper | sniper |
| 2    | badger     | pogreb   | leveldb   | bolt    | nutsdb | buntdb | pebble | pogreb | pogreb |
| 3    | nutsdb     | pebble   | bitcask   | bbolt   | badger | rosedb | rosedb | rosedb | nutsdb |
| 4    | buntdb     | leveldb  | rosedb    | bitcask | rosedb | pogreb | pogreb | nutsdb | buntdb |
| 5    | sniper     | bolt     | nutsdb    | badger  | rosedb | nutsdb | nutsdb | rosedb | badger |
| 6    | leveldb    | bbolt    | rocksdb   | rocksdb | rocksdb | bitcask | rocksdb | buntdb | rosedb |
| 7    | bitcask    | sniper   | sniper    | leveldb | buntdb | badger | buntdb | badger | bitcask |
| 8    | rosedb     | buntdb   | bolt      | pebble  | badger | rocksdb | badger | leveldb | rocksdb |

* pogreb and sniper does not support prefix queries

### fsync

**throughputs**

Coming soon...
