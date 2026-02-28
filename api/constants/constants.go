package constants

import "time"

// Database timeouts
const (
	// DBConnectionTimeout is the timeout for establishing a database connection
	DBConnectionTimeout = 5 * time.Second

	// DBPingTimeout is the timeout for pinging the database
	DBPingTimeout = 5 * time.Second

	// SSMTimeout is the timeout for SSM parameter retrieval
	SSMTimeout = 5 * time.Second
)

// HTTP request timeouts
const (
	// DefaultRequestTimeout is the default timeout for HTTP requests
	DefaultRequestTimeout = 10 * time.Second

	// ShortRequestTimeout is a shorter timeout for simple requests
	ShortRequestTimeout = 5 * time.Second

	// LambdaHandlerTimeout is the timeout for Lambda handler execution
	LambdaHandlerTimeout = 5 * time.Second
)

// Pagination defaults
const (
	// DefaultLimit is the default number of items per page
	DefaultLimit = 20

	// MaxLimit is the maximum number of items allowed per page
	MaxLimit = 100

	// SearchLimit is the limit for search results
	SearchLimit = 10
)

// Database pool configuration
const (
	// DBMaxConns is the maximum number of database connections in the pool
	DBMaxConns = 2

	// DBMinConns is the minimum number of database connections in the pool
	DBMinConns = 1

	// DBMaxConnIdleTime is the maximum time a connection can be idle before being closed
	DBMaxConnIdleTime = 5 * time.Minute

	// DBMaxConnLifetime is the maximum lifetime of a database connection
	DBMaxConnLifetime = 30 * time.Minute
)

// CORS configuration
const (
	// CORSMaxAge is the maximum age for CORS preflight requests
	CORSMaxAge = 300 // 5 minutes
)

// Environment values
const (
	// EnvLocal is the value for local environment
	EnvLocal = "local"

	// EnvProd is the value for production environment
	EnvProd = "prod"
)
