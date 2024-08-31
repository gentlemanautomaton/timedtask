package timedtask

// Option is an option for a timed task.
type Option interface {
	apply(*Task)
}

// QuietOption defines whether logging should be reduced or not.
type QuietOption bool

func (option QuietOption) apply(t *Task) {
	t.quiet = bool(option)
}

// Quiet is a timed task option that reduces logging.
const Quiet = QuietOption(true)

// Check returns a CheckableOption that will add the given checkable to the
// list of error checking functions consulted before starting a task.
func Check(c Checkable) CheckableOption {
	return CheckableOption{checkable: c}
}

// CheckableOption defines an error checking option.
type CheckableOption struct {
	checkable Checkable
}

func (option CheckableOption) apply(t *Task) {
	t.checks = append(t.checks, option.checkable)
}

// Checkable is an interface that can be checked for errors.
type Checkable interface {
	Err() error
}
