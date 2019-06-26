package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Report all errors to stdout.
func handle(err error) {
	if err != nil && err != redis.ErrNil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Scan and queue source keys.
func get(conn redis.Conn, found map[string]bool, queue chan<- map[string]string) {
	var (
		// cursor int64
		keys    []string
		allKeys []string
	)

	start := time.Now()
	fmt.Printf("Fetching keys: %s\n", start.String())
	allKeys, err := redis.Strings(conn.Do("KEYS", "*"))
	handle(err)
	fmt.Printf("Took: %s\n", time.Now().Sub(start))
	fmt.Printf("Total Keys: %d\n", len(keys))

	// Only copy what is no longer in the new cluster.
	for _, key := range allKeys {
		if !found[key] {
			keys = append(keys, key)
		}
	}

	batchSize := 10
	for i := 0; i < len(keys); i += batchSize {
		last := i + batchSize
		if last > len(keys) {
			last = len(keys)
		}

		// Get pipelined dumps.
		for _, key := range keys[i:last] {
			conn.Send("DUMP", key)
		}
		dumps, err := redis.Strings(conn.Do(""))
		handle(err)

		// Build batch map.
		batch := make(map[string]string)
		for j, key := range keys[i:last] {
			batch[key] = dumps[j]
		}

		fmt.Printf(">")
		queue <- batch
	}

	close(queue)
	fmt.Printf("Total Time: %s\n", time.Now().Sub(start))

	// for {
	// 	// Scan a batch of keys.
	// 	values, err := redis.Values(conn.Do("SCAN", cursor))
	// 	handle(err)
	// 	values, err = redis.Scan(values, &cursor, &keys)
	// 	handle(err)

	// 	// Get pipelined dumps.
	// 	for _, key := range keys {
	// 		conn.Send("DUMP", key)
	// 	}
	// 	dumps, err := redis.Strings(conn.Do(""))
	// 	handle(err)

	// 	// Build batch map.
	// 	batch := make(map[string]string)
	// 	for i := range keys {
	// 		batch[keys[i]] = dumps[i]
	// 	}

	// 	// Last iteration of scan.
	// 	if cursor == 0 {
	// 		// queue last batch.
	// 		select {
	// 		case queue <- batch:
	// 		}
	// 		close(queue)
	// 		break
	// 	}

	// 	fmt.Printf(">")
	// 	// queue current batch.
	// 	queue <- batch
	// }
}

// Restore a batch of keys on destination.
// func put(conn redis.Conn, queue <-chan map[string]string) {
func put(conn redis.Conn, queue <-chan map[string]string) {
	for batch := range queue {
		for key, value := range batch {
			conn.Send("RESTORE", key, "0", value)
		}
		_, err := conn.Do("")
		handle(err)

		fmt.Printf(".")
	}
}

func main() {
	from := flag.String("from", "", "example: redis://127.0.0.1:6379/0")
	to := flag.String("to", "", "example: redis://127.0.0.1:6379/1")
	flag.Parse()

	source, err := redis.DialURL(*from)
	handle(err)
	destination, err := redis.DialURL(*to)
	handle(err)
	defer source.Close()
	defer destination.Close()

	// Channel where batches of keys will pass.
	queue := make(chan map[string]string, 100)

	// Get the keys from the `destination` cluster and make them
	// a "set" of strings.
	keys, err := redis.Strings(destination.Do("KEYS", "*"))
	handle(err)
	found := map[string]bool{}
	for _, key := range keys {
		found[key] = true
	}

	// Scan and send to queue.
	// go get(source, queue)
	go get(source, found, queue)

	// Restore keys as they come into queue.
	put(destination, queue)

	fmt.Println("Sync done.")
}
