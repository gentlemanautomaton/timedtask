package timedtask

// Func is a timed task function that can have subtasks.
//
// The provided task can be used to start subtasks. It must not be retained
// after the function returns.
type Func func(t *Task) error

// SimpleFunc is a simple timed task function without subtasks.
type SimpleFunc func() error

// Func is a timed task function that can have subtasks and returns a value
// of type T.
//
// The provided task can be used to start subtasks. It must not be retained
// after the function returns.
type FuncFor[T any] func(t *Task) (T, error)

// SimpleFunc is a simple timed task function without subtasks that returns
// a value of type T.
type SimpleFuncFor[T any] func() (T, error)
