package primenumbers

import (
	"fmt"

	"github.com/project-flogo/core/activity"
)

// Settings for the activity (empty for now)
type Settings struct{}

// Input defines the input structure
type Input struct {
	Start int `md:"start,required"`
	End   int `md:"end,required"`
}

// Output defines the output structure
type Output struct {
	Primes []int `md:"primes"`
}

// activityMd holds the metadata for this activity
var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

type Activity struct{}

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Metadata returns the activity's metadata.
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// New creates a new instance of the Activity.
func New(ctx activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

func (a *Activity) Eval(ctx activity.Context) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("activity context is nil")
	}

	// Get input values safely
	startInput := ctx.GetInput("start")
	if startInput == nil {
		return false, fmt.Errorf("'start' input is nil")
	}

	start, ok := startInput.(int)
	if !ok {
		return false, fmt.Errorf("invalid 'start' input type, expected int, got %T", startInput)
	}

	endInput := ctx.GetInput("end")
	if endInput == nil {
		return false, fmt.Errorf("'end' input is nil")
	}

	end, ok := endInput.(int)
	if !ok {
		return false, fmt.Errorf("invalid 'end' input type, expected int, got %T", endInput)
	}

	// Generate primes in the range from Start to End
	primes := generatePrimesInRange(start, end)

	// Set output
	ctx.SetOutput("primes", primes)

	return true, nil
}

func generatePrimesInRange(start, end int) []int {
	primes := []int{}
	for num := start; num <= end; num++ {
		if isPrime(num) {
			primes = append(primes, num)
		}
	}
	return primes
}

func isPrime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}
