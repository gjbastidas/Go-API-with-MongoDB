package constants

import "time"

const (
	ServerTimeout  = 1 * time.Second           // This is to timeout server shutdown
	RequestTimeout = 10 * time.Second          // This is to timeout requests
	DbName         = "simple-api-with-mongodb" // Database name
	PColl          = "posts"                   // Post collection name
	CColl          = "comments"                // Comments collection name
)
