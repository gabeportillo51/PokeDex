package pokecache
import ("sync"; "time") 

type Cache struct {
	// This is the struct that represents the entire cache
	cache map[string]cacheEntry
	mu *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	// This struct represents an individual entry in the cache
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache {
	new_cache := Cache{
		cache: make(map[string]cacheEntry),
		mu: &sync.Mutex{},
		interval: interval,
	}
	go new_cache.reapLoop()
	return &new_cache
}

func (c *Cache) Add(key string, val []byte){
	// method to add entries to the cache
	c.mu.Lock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// method to retrieve the data in a cache entry
	c.mu.Lock()
	entry, ok := c.cache[key]
	c.mu.Unlock()
	if ok {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	// method that refreshes the cache so that it doesn't get too big
	ticker := time.NewTicker(c.interval)    // create a ticker that ticks every interval
	defer ticker.Stop()                    // ensure the ticker closes after all of the below code executes
	for range ticker.C {                   // execute this loop every time the ticker ticks
		c.mu.Lock()                        // lock access to the map
		now := time.Now()                    // get the current time
		for key, entry := range c.cache {           // iterate through the cache entries in the map
			if now.Sub(entry.createdAt) > c.interval{ // check the entries age. If older than interval Delete it.
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
