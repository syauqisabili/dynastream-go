package network

import "sync"

var (
	instance *NetCfg
	once     sync.Once
	mu       sync.RWMutex // Mutex to access thread-safe
)

// Automatically initialize the singleton instance
func init() {
	once.Do(func() {
		instance = &NetCfg{} // Initialize the instance
	})
}

// Set value for thread-safe
func Set(config NetCfg) {
	mu.Lock()
	defer mu.Unlock()
	*instance = config
}

// Get value
func Get() NetCfg {
	mu.RLock()
	defer mu.RUnlock()
	return *instance
}
