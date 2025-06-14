package timedtask

import (
	"context"
	"io"
)

// Spec is a specification for a timed task.
type Spec struct {
	// Parent is a parent task within which the task will run.
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
func (s Spec) Run(f Func) error {
	return s.RunCtx(context.Background(), f)
}

// Run executes the given function as a timed task with the parameters from s.
//
// If the given context returns an error when RunCtx is called, it will
// be returned without running the function.
func (s Spec) RunCtx(ctx context.Context, f Func) error {
	if err := ctx.Err(); err != nil {
		return err
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
	err := f(&task)
	task.end(err)

	return err
}

// RunSimple executes the given simple function as a timed task with the
// parameters from s.
func (s Spec) RunSimple(f SimpleFunc) error {
	return s.RunCtx(context.Background(), func(t *Task) error {
		return f()
	})
}

// RunSimpleCtx executes the given simple function as a timed task with the
// parameters from s.
//
// If the given context returns an error when RunSimpleCtx is called, it will
// be returned without running the function.
func (s Spec) RunSimpleCtx(ctx context.Context, f SimpleFunc) error {
	return s.RunCtx(ctx, func(t *Task) error {
		return f()
	})
}

func childTaskDepth(parent *Task) int {
	if parent == nil {
		return 0
	}
	return parent.depth + 1
}
