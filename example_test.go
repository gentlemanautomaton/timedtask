package timedtask_test

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/timedtask"
)

func Example() {
	timedtask.Run("Downloading files", timedtask.Func(func(task *timedtask.Task) error {
		for i := 1; i <= 10; i++ {
			if download := task.Run(fmt.Sprintf("Downloading file %d", i), timedtask.Func(func(task *timedtask.Task) error {
				// Download code goes here.
				if i == 5 {
					const simulatedRetries = 3
					task.Logf("Download took %d retries.", simulatedRetries)
				}
				return nil
			})); download.Err != nil {
				return download.Err
			}
		}
		return nil
	}))

	timedtask.Run("Crunching data", timedtask.Func(func(task *timedtask.Task) error {
		if computations := task.Run("Validating", timedtask.SimpleFunc(func() error {
			// Data processing code goes here.
			return nil
		})); computations.Err != nil {
			return computations.Err
		}

		if computations := task.Run("Running solver", timedtask.Func(func(task *timedtask.Task) error {
			return task.Run("Solving all the things", timedtask.SimpleFunc(func() error {
				return errors.New("a solution did not present itself")
			})).Err
		})); computations.Err != nil {
			return computations.Err
		}

		return nil
	}))

	// Output:
	// Downloading files...
	//   Downloading file 1... done. 0s
	//   Downloading file 2... done. 0s
	//   Downloading file 3... done. 0s
	//   Downloading file 4... done. 0s
	//   Downloading file 5...
	//     Download took 3 retries.
	//   Downloading file 5... done. 0s
	//   Downloading file 6... done. 0s
	//   Downloading file 7... done. 0s
	//   Downloading file 8... done. 0s
	//   Downloading file 9... done. 0s
	//   Downloading file 10... done. 0s
	// Downloading files... done. 0s
	// Crunching data...
	//   Validating... done. 0s
	//   Running solver...
	//     Solving all the things... failed. 0s
	//   Running solver... failed. 0s
	// Crunching data... failed. 0s
}
