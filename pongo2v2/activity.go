package pongo2v2prompt

import (
	"github.com/flosch/pongo2/v6"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

// activityMd holds the metadata for this activity
var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	activity.Register(&ImprovedActivity{}, New)
}

// ImprovedActivity represents the enhanced pongo2-prompt activity with better type handling
type ImprovedActivity struct {
	template string
}

// New creates a new instance of the ImprovedActivity (standard Flogo interface)
func New(ctx activity.InitContext) (activity.Activity, error) {
	return NewImproved(ctx)
}

// NewImproved creates a new instance of the ImprovedActivity
func NewImproved(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Info("Improved Pongo2-prompt activity initialized with enhanced type handling")
	return &ImprovedActivity{template: s.Template}, nil
}

// Eval executes the improved activity with better data type handling
func (a *ImprovedActivity) Eval(ctx activity.Context) (done bool, err error) {
	template := ctx.GetInput("template").(string)
	variables := ctx.GetInput("variables")

	// Convert variables to a map
	var varMap map[string]interface{}
	if variables != nil {
		varMap = variables.(map[string]interface{})

		// Fix float to int conversion for numeric comparisons in templates
		for key, value := range varMap {
			if f, ok := value.(float64); ok {
				// If it's a whole number, convert to int for proper template comparison
				if f == float64(int64(f)) {
					varMap[key] = int64(f)
				}
			}
		}
	}

	// Create template
	tpl, err := pongo2.FromString(template)
	if err != nil {
		return false, err
	}

	// Execute template
	out, err := tpl.Execute(pongo2.Context(varMap))
	if err != nil {
		return false, err
	}

	ctx.SetOutput("result", out)
	return true, nil
}

// Metadata returns the activity's metadata (same as original)
func (a *ImprovedActivity) Metadata() *activity.Metadata {
	return activityMd
}
