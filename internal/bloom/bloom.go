package bloom

import (
	"sync"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	Store    []bool            // to Store Bitset
	Size     int32             // Size of the filters
	mHashers []murmur3.Hash128 // HashFunctions to improve the probabilistic accuracy

	mu sync.Mutex
}

// NewBloomFilter helps to create a new bloom.Filter with the given size.
//
// It initializes the bloom.Store with the size and generates the
// murmur3 hash functions with different seeds.
//
// The number of hash functions is fixed to 8.
//
// Returns a new bloom.Filter.
func NewBloomFilter(size int32) *BloomFilter {
	// Generating the Seed and MHashers (Murmur 128 style)
	//
	return &BloomFilter{
		Store: make([]bool, size),
		Size:  size,
		mHashers: []murmur3.Hash128{
			murmur3.New128WithSeed(uint32(11)),
			murmur3.New128WithSeed(uint32(31)),
			murmur3.New128WithSeed(uint32(131)),
			murmur3.New128WithSeed(uint32(989)),
			murmur3.New128WithSeed(uint32(1919)),
			murmur3.New128WithSeed(uint32(2007)),
			murmur3.New128WithSeed(uint32(31313)),
			murmur3.New128WithSeed(uint32(9281917)),
		},
	}
}

// Info returns the required information for the bloom.Filter configuration
func (bf *BloomFilter) Info() map[string]any {
	return map[string]any{
		"size":           bf.Size,
		"totalHashFuncs": len(bf.mHashers),
	}
}

// ComputeMurmurHash computes and returns the murmur has of a query string `key`
// you have to module it with the bloom.Size to set the index True in bloom.Store
//
// Non-cryptographic hash, fast and efficient and implementation specific.
func (bf *BloomFilter) ComputeMurmurHash(key string, hashFn int) uint64 {
	bf.mHashers[hashFn].Write([]byte(key))
	val, _ := bf.mHashers[hashFn].Sum128()
	bf.mHashers[hashFn].Reset()
	return val
}

// Add helps to add add given key into the bloom.Store. Remember it does not store the
// actual keys. rather it a probabilistic representtal of their presence.
func (bf *BloomFilter) Add(key string) {
	// index := bf.ComputeMurmurHash(key) % uint64(bf.Size)
	// bf.Store[index] = true
	// Utilizing all has functions
	bf.mu.Lock()
	for i := 0; i < len(bf.mHashers); i++ {
		index := bf.ComputeMurmurHash(key, i) % uint64(bf.Size)
		bf.Store[index] = true
	}
	bf.mu.Unlock()
}

// Exists helps to lookup if the key present in the bloom.Store.
// In real, the key might not be present even if the return is true. as it
// works as a probabilistic estimation of finding the presence.
func (bf *BloomFilter) Exists(key string) (uint64, bool) {
	// index := bf.ComputeMurmurHash(key) % uint64(bf.Size)
	// return index, bf.Store[index]

	for i := 0; i < len(bf.mHashers); i++ {
		index := bf.ComputeMurmurHash(key, i) % uint64(bf.Size)
		if !bf.Store[index] {
			return index, false
		}
	}
	return 0, true
}
