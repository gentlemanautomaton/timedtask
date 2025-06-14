package timedtask

import (
	"context"
	"io"
)

// SpecFor is a specification for a timed task that returns type T.
type SpecFor[T any] struct {
	// Parent is a parent task to run the task under.
	Parent *Task

	// Description is a description of the task.
	Description string

	// Quiet tasks will not be printed under normal circumstances.
	Quiet bool

	// Writer specifies a writer for the task to write its output to.
	// If nil, it will use the writer of its parent task, if available,
	// othwerwise it will write to os.Stdout.
	Writer io.Writer
}

// Run executes the given function as a timed task with the parameters from s.
func (s SpecFor[T]) Run(f FuncFor[T]) (T, error) {
	return s.RunCtx(context.Background(), f)
}

// Run executes the given function as a timed task with the parameters from s.
//
// If the given context returns an error when RunCtx is called, it will
// be returned without running the function.
func (s SpecFor[T]) RunCtx(ctx context.Context, f FuncFor[T]) (T, error) {
	if err := ctx.Err(); err != nil {
		var empty T
		return empty, err
	}

	writer := s.Writer
	if writer == nil && s.Parent != nil {
		writer = s.Parent.writer
	}

	task := Task{
		parent:      s.Parent,
		description: s.Description,
		depth:       childTaskDepth(s.Parent),
		quiet:       s.Quiet,
		writer:      writer,
	}

	task.start()
	result, err := f(&task)
	task.end(err)

	return result, err
}

// RunSimple executes the given simple function as a timed task with the
// parameters from s.
func (s SpecFor[T]) RunSimple(f SimpleFuncFor[T]) (T, error) {
	return s.RunCtx(context.Background(), func(t *Task) (T, error) {
		return f()
	})
}

// RunSimpleCtx executes the given simple function as a timed task with the
// parameters from s.
//
// If the given context returns an error when RunSimpleCtx is called, it will
// be returned without running the function.
func (s SpecFor[T]) RunSimpleCtx(ctx context.Context, f SimpleFuncFor[T]) (T, error) {
	return s.RunCtx(ctx, func(t *Task) (T, error) {
		return f()
	})
}
