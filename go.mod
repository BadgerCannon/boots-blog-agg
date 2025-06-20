module github.com/BadgerCannon/boot-blog-agg

go 1.23.4

// replace internal/config v0.0.0 => ./internal/config/

// require internal/config v0.0.0

// replace internal/database v0.0.0 => ./internal/database/

// require internal/database v0.0.0

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)
