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

	notes     []string
	startTime time.Time
	endTime   time.Time
	flushed   bool
}

// Duration returns the duration of the task.
//
// If the task is still running, it returns the task duration so far.
func (task *Task) Duration() time.Duration {
	if task.endTime.IsZero() {
		return time.Since(task.startTime)
	}
	return task.endTime.Sub(task.startTime)
}

// AddNote adds a note to the task that will be reported at its completion.
//
// If the note is empty it will not be added.
func (task *Task) AddNote(note string) {
	if note == "" {
		return
	}
	task.notes = append(task.notes, note)
}

// AddNoteWithLabel adds a note to the task that will be reported at its
// completion. The note will be prefixed with the given label.
//
// If the note is empty it will not be added.
func (task *Task) AddNoteWithLabel(label, note string) {
	if note == "" {
		return
	}
	if label != "" {
		note = label + ": " + note
	}
	task.notes = append(task.notes, note)
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

// start marks the task start time and prints its description if the
// task is not quiet.
func (task *Task) start() {
	if !task.quiet {
		if task.parent != nil {
			task.parent.flush()
		}
		task.logf(task.depth, "%s...", task.description)
	}
	task.startTime = time.Now()
}

// start marks the task end time and prints its result if the
// task is not quiet or if an error was encountered.
func (task *Task) end(taskErr error) {
	task.endTime = time.Now()

	notes := []string{task.Duration().Round(time.Millisecond).String()}
	notes = append(notes, task.notes...)
	suffix := "(" + strings.Join(notes, ", ") + ")"

	if taskErr != nil {
		if task.flushed {
			task.logf(task.depth, "%s... failed. %s\n", task.description, suffix)
			return
		}

		if task.quiet {
			if task.parent != nil {
				task.parent.flush()
			}
			task.logf(task.depth, "%s... failed. %s\n", task.description, suffix)
			return
		}

		task.logf(0, " failed. %s\n", suffix)
	} else {
		if task.flushed {
			task.logf(task.depth, "%s... done. %s\n", task.description, suffix)
			return
		}

		if !task.quiet {
			task.logf(0, " done. %s\n", suffix)
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
