package timedtask_test

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gentlemanautomaton/timedtask"
)

func Example() {
	downloadTask := timedtask.Spec{
		Description: "Downloading files",
	}
	downloadTask.Run(func(downloadTask *timedtask.Task) error {
		const fileCount = 10
		downloadTask.AddNote(fmt.Sprintf("%d files", fileCount))

		for i := 1; i <= fileCount; i++ {
			fileTask := timedtask.Spec{
				Description: fmt.Sprintf("Downloading file %d", i),
				Parent:      downloadTask,
			}
			if err := fileTask.Run(func(fileTask *timedtask.Task) error {
				// Download code goes here.
				if i == 5 {
					const simulatedRetries = 3
					fileTask.Logf("Download took %d retries.", simulatedRetries)
				}
				return nil
			}); err != nil {
				return err
			}
		}

		return nil
	})

	output, err := timedtask.SpecFor[int]{Description: "Crunching data"}.Run(func(numbersTask *timedtask.Task) (int, error) {
		validationTask := timedtask.Spec{
			Description: "Validating",
			Parent:      numbersTask,
		}
		err := validationTask.RunSimple(func() error {
			// Data validation code goes here.
			return nil
		})
		if err != nil {
			return 0, err
		}

		solverTask := timedtask.SpecFor[int]{Description: "Running solver", Parent: numbersTask}
		result, err := solverTask.Run(func(solverTask *timedtask.Task) (int, error) {
			const rounds = 5
			solverTask.AddNoteWithLabel("Rounds", strconv.Itoa(rounds))
			return timedtask.SpecFor[int]{Description: "Solving all the things", Parent: solverTask}.RunSimple(func() (int, error) {
				// Data processing code goes here.
				value := 7
				for i := 1; i < rounds; i++ {
					value *= i
				}
				return value, errors.New("a solution did not present itself")
			})
		})

		return result, err
	})

	fmt.Printf("Result: %d, Error: %s\n", output, err)

	// Output:
	// Downloading files...
	//   Downloading file 1... done. (0s)
	//   Downloading file 2... done. (0s)
	//   Downloading file 3... done. (0s)
	//   Downloading file 4... done. (0s)
	//   Downloading file 5...
	//     Download took 3 retries.
	//   Downloading file 5... done. (0s)
	//   Downloading file 6... done. (0s)
	//   Downloading file 7... done. (0s)
	//   Downloading file 8... done. (0s)
	//   Downloading file 9... done. (0s)
	//   Downloading file 10... done. (0s)
	// Downloading files... done. (0s, 10 files)
	// Crunching data...
	//   Validating... done. (0s)
	//   Running solver...
	//     Solving all the things... failed. (0s)
	//   Running solver... failed. (0s, Rounds: 5)
	// Crunching data... failed. (0s)
	// Result: 168, Error: a solution did not present itself
}
