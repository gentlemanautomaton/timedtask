package timedtask

// Func is a timed task function that can have subtasks.
//
// The provided task can be used to start subtasks. It must not be retained
// after the function returns.
type Func func(t *Task) error

func (f Func) apply(t *Task) {
	t.functions = append(t.functions, f)
}

// SimpleFunc is a simple timed task function without subtasks.
type SimpleFunc func() error

func (f SimpleFunc) apply(t *Task) {
	t.functions = append(t.functions, func(*Task) error {
		return f()
	})
}
