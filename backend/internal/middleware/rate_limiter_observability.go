package middleware

// Snapshot returns a copy of per-IP counters for observability.
func (r *RateLimiter) Snapshot() map[string]int {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]int, len(r.limits))
	for ip, b := range r.limits {
		out[ip] = b.count
	}
	return out
}

// Reset clears the counter for one client IP.
func (r *RateLimiter) Reset(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.limits, ip)
}

// Count returns the number of tracked client IPs.
func (r *RateLimiter) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.limits)
}

// Peek returns the current counter for one client IP.
func (r *RateLimiter) Peek(ip string) (int, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.limits[ip]
	if !ok {
		return 0, false
	}
	return b.count, true
}

// ClearAll drops every tracked client.
func (r *RateLimiter) ClearAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for ip := range r.limits {
		delete(r.limits, ip)
	}
}
