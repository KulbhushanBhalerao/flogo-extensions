package primenumbers

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/support/log"
)

func TestPrimeNumbersActivity(t *testing.T) {
	// Create activity with mock init context
	initCtx := &MockInitContext{}
	act, err := New(initCtx)
	if err != nil {
		t.Fatalf("Failed to create activity: %v", err)
	}

	// Create evaluation context with start=2, end=10
	evalCtx := &MockActivityContext{
		inputs: map[string]interface{}{
			"start": 2,
			"end":   10,
		},
		outputs: make(map[string]interface{}),
	}

	success, err := act.Eval(evalCtx)
	if err != nil || !success {
		t.Fatalf("Failed to evaluate activity: %v", err)
	}

	primes := evalCtx.outputs["primes"].([]int)
	expected := []int{2, 3, 5, 7}

	if len(primes) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, primes)
	}

	for i, prime := range primes {
		if prime != expected[i] {
			t.Errorf("Expected %v, got %v", expected[i], prime)
		}
	}
}

func TestPrimeNumbersActivityRange(t *testing.T) {
	// Test with different range
	initCtx := &MockInitContext{}
	act, err := New(initCtx)
	if err != nil {
		t.Fatalf("Failed to create activity: %v", err)
	}

	evalCtx := &MockActivityContext{
		inputs: map[string]interface{}{
			"start": 10,
			"end":   20,
		},
		outputs: make(map[string]interface{}),
	}

	success, err := act.Eval(evalCtx)
	if err != nil || !success {
		t.Fatalf("Failed to evaluate activity: %v", err)
	}

	primes := evalCtx.outputs["primes"].([]int)
	expected := []int{11, 13, 17, 19}

	if len(primes) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, primes)
	}

	for i, prime := range primes {
		if prime != expected[i] {
			t.Errorf("Expected %v, got %v", expected[i], prime)
		}
	}
}

// MockInitContext is a mock implementation of activity.InitContext for testing purposes.
type MockInitContext struct{}

func (m *MockInitContext) Settings() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *MockInitContext) MapperFactory() mapper.Factory {
	return nil
}

func (m *MockInitContext) Logger() log.Logger {
	return &MockLogger{}
}

func (m *MockInitContext) GetSetting(setting string) (value interface{}, exists bool) {
	return nil, false
}

// MockActivityContext is a mock implementation of activity.Context for testing purposes.
type MockActivityContext struct {
	inputs  map[string]interface{}
	outputs map[string]interface{}
}

func (m *MockActivityContext) GetInput(name string) interface{} {
	return m.inputs[name]
}

func (m *MockActivityContext) SetOutput(name string, value interface{}) error {
	m.outputs[name] = value
	return nil
}

func (m *MockActivityContext) ActivityHost() activity.Host {
	return nil
}

func (m *MockActivityContext) GetSharedTempData() map[string]interface{} {
	return nil
}

func (m *MockActivityContext) Logger() log.Logger {
	return &MockLogger{}
}

func (m *MockActivityContext) Name() string {
	return "MockActivity"
}

func (m *MockActivityContext) SetOutputObject(obj data.StructValue) error {
	return obj.FromMap(m.outputs)
}

func (m *MockActivityContext) GetInputObject(obj data.StructValue) error {
	return obj.FromMap(m.inputs)
}

// MockLogger is a mock implementation of log.Logger for testing purposes.
type MockLogger struct{}

func (l *MockLogger) Info(args ...interface{})                  {}
func (l *MockLogger) Infof(format string, args ...interface{})  {}
func (l *MockLogger) Debug(args ...interface{})                 {}
func (l *MockLogger) Debugf(format string, args ...interface{}) {}
func (l *MockLogger) Warn(args ...interface{})                  {}
func (l *MockLogger) Warnf(format string, args ...interface{})  {}
func (l *MockLogger) Error(args ...interface{})                 {}
func (l *MockLogger) Errorf(format string, args ...interface{}) {}
func (l *MockLogger) DebugEnabled() bool                        { return false }
func (l *MockLogger) Structured() log.StructuredLogger {
	return nil
}

func (l *MockLogger) Trace(args ...interface{})                 {}
func (l *MockLogger) Tracef(format string, args ...interface{}) {}
func (l *MockLogger) TraceEnabled() bool                        { return false }
