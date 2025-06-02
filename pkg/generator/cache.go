package generator

import (
	"sync"
	"text/template"
	"time"
)

// TemplateCache provides thread-safe caching for parsed templates
type TemplateCache struct {
	cache    map[string]*CachedTemplate
	mu       sync.RWMutex
	maxAge   time.Duration
	hitCount int64
	missCount int64
}

// CachedTemplate represents a cached template with metadata
type CachedTemplate struct {
	Template  *template.Template
	CreatedAt time.Time
	AccessCount int64
	LastAccess time.Time
}

// NewTemplateCache creates a new template cache
func NewTemplateCache(maxAge time.Duration) *TemplateCache {
	return &TemplateCache{
		cache:   make(map[string]*CachedTemplate),
		maxAge:  maxAge,
	}
}

// Get retrieves a template from cache
func (tc *TemplateCache) Get(key string) (*template.Template, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	cached, exists := tc.cache[key]
	if !exists {
		tc.missCount++
		return nil, false
	}

	// Check if template has expired
	if time.Since(cached.CreatedAt) > tc.maxAge {
		tc.mu.RUnlock()
		tc.mu.Lock()
		delete(tc.cache, key)
		tc.mu.Unlock()
		tc.mu.RLock()
		tc.missCount++
		return nil, false
	}

	// Update access statistics
	cached.AccessCount++
	cached.LastAccess = time.Now()
	tc.hitCount++

	return cached.Template, true
}

// Put stores a template in cache
func (tc *TemplateCache) Put(key string, tmpl *template.Template) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cache[key] = &CachedTemplate{
		Template:  tmpl,
		CreatedAt: time.Now(),
		AccessCount: 0,
		LastAccess: time.Now(),
	}
}

// Clear removes all templates from cache
func (tc *TemplateCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cache = make(map[string]*CachedTemplate)
	tc.hitCount = 0
	tc.missCount = 0
}

// Stats returns cache statistics
func (tc *TemplateCache) Stats() CacheStats {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	total := tc.hitCount + tc.missCount
	hitRatio := float64(0)
	if total > 0 {
		hitRatio = float64(tc.hitCount) / float64(total)
	}

	return CacheStats{
		Size:     len(tc.cache),
		Hits:     tc.hitCount,
		Misses:   tc.missCount,
		HitRatio: hitRatio,
	}
}

// CacheStats represents cache performance statistics
type CacheStats struct {
	Size     int
	Hits     int64
	Misses   int64
	HitRatio float64
}

// CleanupExpired removes expired templates from cache
func (tc *TemplateCache) CleanupExpired() int {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	now := time.Now()
	removedCount := 0

	for key, cached := range tc.cache {
		if now.Sub(cached.CreatedAt) > tc.maxAge {
			delete(tc.cache, key)
			removedCount++
		}
	}

	return removedCount
}

// StartCleanupWorker starts a background worker to clean up expired templates
func (tc *TemplateCache) StartCleanupWorker(interval time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tc.CleanupExpired()
		case <-done:
			return
		}
	}
}