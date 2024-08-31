package timedtask

import "time"

// Summary returns a summary of the task's execution.
type Summary struct {
	Start time.Time
	End   time.Time
	Err   error
}

// Duration returns the duration of the task.
func (summary Summary) Duration() time.Duration {
	return summary.End.Sub(summary.Start)
}
