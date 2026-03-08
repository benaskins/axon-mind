package mind

import "encoding/json"

// JSON returns the solution's variable bindings as a JSON object.
func (s Solution) JSON() ([]byte, error) {
	return json.Marshal(s.Bindings)
}

// SolutionsJSON returns a JSON array of solution objects.
// Returns "[]" for nil or empty slices.
func SolutionsJSON(solutions []Solution) ([]byte, error) {
	if len(solutions) == 0 {
		return []byte("[]"), nil
	}
	arr := make([]map[string]any, len(solutions))
	for i, s := range solutions {
		arr[i] = s.Bindings
	}
	return json.Marshal(arr)
}
