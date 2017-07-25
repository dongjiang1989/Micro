package systeminfo

type STATE string

// TODO: add PROCESS STATE CODES
const (
	RUNNING STATE = "Running" // Running or runnable (on run queue)
	WAITING STATE = "Waiting" // Interruptible sleep (waiting for an event to complete
	STOPPED STATE = "Stopped" // Stopped, either by a job control signal or because it is being traced.
)
