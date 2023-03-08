package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"

	hashi_aes_gcm "github.com/valli0x/hashi-vault-barrier-aes-gcm"
	hashi_phicical "github.com/valli0x/hashi-vault-physical"
)

func main() {
	/*
		create physical backend
	*/

	typePhisical := "inmem"       // backend type
	config := map[string]string{} // backend config
	logger := log.New(nil)        // backend logger

	factory := hashi_phicical.PhysicalBackends[typePhisical]

	// create backend
	phicicalBackend, err := factory(config, logger.Named("physical"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// put entry
	phicicalBackend.Put(context.Background(), &physical.Entry{
		Key:   "key-example",
		Value: []byte("phisical entry"),
	})

	// get entry
	entry, err := phicicalBackend.Get(context.Background(), "key-example")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(entry.Value))

	/*
		create cache backend
	*/

	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	cacheBackend := physical.NewCache(phicicalBackend, 0, logger.Named("cache"), inm)

	// put entry
	cacheBackend.Put(context.Background(), &physical.Entry{
		Key:   "key-example",
		Value: []byte("cache entry"),
	})

	// get entry
	cacheEntry, err := cacheBackend.Get(context.Background(), "key-example")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(cacheEntry.Value))

	// clear cache
	cacheBackend.Purge(context.Background())

	/*
		create aes-gcm backend
	*/

	aesBackend, err := hashi_aes_gcm.NewAESGCMBarrier(cacheBackend)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialize and unseal
	key, err := aesBackend.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	// init aes-gcm barrier
	aesBackend.Initialize(context.Background(), key, nil, rand.Reader)
	// unseal barrier
	aesBackend.Unseal(context.Background(), key)

	// put entry
	aesBackend.Put(context.Background(), &logical.StorageEntry{
		Key:   "key-example",
		Value: []byte("storage entry"),
	})

	// get entry
	storageEntry, err := aesBackend.Get(context.Background(), "key-example")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(storageEntry.Value))

	/*
		create barrierView with prefix
	*/

	barrierView := logical.NewStorageView(aesBackend, "prefix-example")

	// put entry
	barrierView.Put(context.Background(), &logical.StorageEntry{
		Key:   "key-example",
		Value: []byte("storageView entry"),
	})

	// get entry
	barrierViewEntry, err := barrierView.Get(context.Background(), "key-example")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(barrierViewEntry.Value))
}
