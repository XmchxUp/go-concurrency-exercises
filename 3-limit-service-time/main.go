//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int32 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	done := make(chan struct{})
	start := time.Now()
	go func() {
		process()
		done <- struct{}{}
	}()

	for {
		elapsed := int32(time.Since(start).Seconds())
		currentTimeUsed := atomic.LoadInt32(&u.TimeUsed)
		if currentTimeUsed+elapsed >= 10 {
			atomic.StoreInt32(&u.TimeUsed, 10)
			return false
		}

		select {
		case <-done:
			atomic.AddInt32(&u.TimeUsed, int32(time.Since(start).Seconds()))
			return true
		case <-time.After(100 * time.Millisecond):
		}

	}
}

func main() {
	RunMockServer()
}
