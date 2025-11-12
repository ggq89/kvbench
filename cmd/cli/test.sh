#!/bin/sh

SIZE=256

# STORES=("badger" "bbolt" "bolt" "leveldb" "kv" "buntdb" "pebble" "pogreb" "nutsdb" "rocksdb" "btree" "btree/memory" "map" "map/memory")
STORES=("badger" "bbolt" "bolt" "buntdb" "leveldb" "rocksdb" "pebble" "pogreb" "nutsdb" "sniper" "btree" "btree/memory" "map" "map/memory" "bitcask" "rosedb")

export LD_LIBRARY_PATH=/usr/local/lib

# CGO_CFLAGS="-I/usr/local/include/rocksdb" CGO_LDFLAGS="-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4"   go build main.go

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`
`rm -f benchmarks/test.log`

echo "=========== test nofsync ==========="
for i in "${STORES[@]}"
do
#	./main -d 1m -size ${SIZE} -s "$i" >> benchmarks/test.log 2>&1
  echo "$i"
	./cli -d 10s -size ${SIZE} -s "$i" -save "benchmarks/nofsync.csv" >> benchmarks/test.log 2>&1
done

`rm  -fr .*db`
`rm  -fr *.db`
`rm  -fr pogreb.*`

echo ""
echo "=========== test fsync ==========="

for i in "${STORES[@]}"
do
#   ./main -d 1m -size ${SIZE} -s "$i" -fsync >> benchmarks/test.log 2>&1
  echo "$i"
	./cli -d 10s -size ${SIZE} -s "$i" -save "benchmarks/fsync.csv" -fsync >> benchmarks/test.log 2>&1
done

`rm  -fr .*db` 
`rm  -fr *.db`
`rm  -fr pogreb.*`