package registry

// Watcher is an interface that returns updates
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

// Result is returned by a call to Next on the watcher
type Result struct {
	Action  string
	Service *Service
}
