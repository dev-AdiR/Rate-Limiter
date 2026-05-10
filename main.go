// package main

// import (
// 	"rate-limiter/proxy"
// 	ratelimit "rate-limiter/rate_limit"
// )

// func main() {
// 	proxy.RegisterProxy()
// 	ratelimit.RateLimit()
// }

package main

import (
	"rate-limiter/proxy"
)

func main() {
	proxy.RegisterProxy()
}
