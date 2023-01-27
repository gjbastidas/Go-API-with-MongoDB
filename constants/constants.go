package constants

import "time"

const (
	ServerTimeout  = 1 * time.Second  // this is to timeout server shutdown
	RequestTimeout = 10 * time.Second // this is to timeout requests
	PColl          = "posts"          // Post collection name
	CColl          = "comments"       // Comments collection name
)
