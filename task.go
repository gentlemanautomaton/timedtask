package timedtask

import (
	"fmt"
	"strings"
	"time"
)

// Task is a timed task.
type Task struct {
	parent      *Task
	description string
	depth       int // The nesting depth of the task, 0 for root tasks
	quiet       bool
	checks      []Checkable
	functions   []Func

	start   time.Time
	end     time.Time
	flushed bool
	err     error
}

// Run runs the given simple task function as a timed task.
//
// Information about the task is printed to stdout.
func Run(description string, options ...Option) Summary {
	t := Task{
		description: description,
	}
	for _, opt := range options {
		opt.apply(&t)
	}
	t.run()
	return Summary{
		Start: t.start,
		End:   t.end,
		Err:   t.err,
	}
}

// Run runs a subtask with the given options.
func (task *Task) Run(description string, options ...Option) Summary {
	child := Task{
		parent:      task,
		description: description,
		depth:       task.depth + 1,
	}
	for _, opt := range options {
		opt.apply(&child)
	}
	child.run()
	return Summary{
		Start: child.start,
		End:   child.end,
		Err:   child.err,
	}
}

// Duration returns the duration of the task.
//
// If the task is still running, it returns the task duration so far.
func (task *Task) Duration() time.Duration {
	if task.end.IsZero() {
		return time.Since(task.start)
	}
	return task.end.Sub(task.start)
}

// Logf prints the given format and values to stdout with an indendation
// level matching the task's nesting depth.
func (task *Task) Logf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}

	task.flush()
	task.log(indent(task.depth + 1))
	task.log(s)
}

// logf prints the given string format and values to the log.
func (task *Task) logf(depth int, format string, a ...any) {
	if depth > 0 {
		task.log(indent(depth))
	}
	task.log(fmt.Sprintf(format, a...))
}

// log writes the given string to the log.
func (task *Task) log(s string) {
	fmt.Print(s)
}

// run executes the task's function and logs the result.
func (task *Task) run() {
	if !task.quiet {
		if task.parent != nil {
			task.parent.flush()
		}
		task.logf(task.depth, "%s...", task.description)
	}

	task.start = time.Now()
	for _, check := range task.checks {
		if task.err != nil {
			break
		}
		task.err = check.Err()
	}
	for _, function := range task.functions {
		if task.err != nil {
			break
		}
		task.err = function(task)
	}
	task.end = time.Now()

	duration := task.Duration().Round(time.Millisecond)

	if task.err != nil {
		if task.flushed {
			task.logf(task.depth, "%s... failed. %s\n", task.description, duration)
			return
		}

		if task.quiet {
			if task.parent != nil {
				task.parent.flush()
			}
			task.logf(task.depth, "%s... failed. %s\n", task.description, duration)
			return
		}

		task.logf(0, " failed. %s\n", duration)
	} else {
		if task.flushed {
			task.logf(task.depth, "%s... done. %s\n", task.description, duration)
			return
		}

		if !task.quiet {
			task.logf(0, " done. %s\n", duration)
		}
	}
}

// flush causes the task to flush any pending text to the log ahead of an
// unrelated write.
func (task *Task) flush() {
	if task.flushed {
		return
	}
	if task.parent != nil {
		task.parent.flush()
	}
	if task.quiet {
		task.logf(task.depth, "%s...\n", task.description)
	} else {
		task.log("\n")
	}
	task.flushed = true
}

// indent returns indentation for the given task depth.
func indent(depth int) string {
	return strings.Repeat(" ", depth*2)
}
