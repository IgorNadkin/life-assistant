package graph

func Evaluate(cond map[string]any, answers map[string]any) bool {
	typ := cond["type"].(string)

	switch typ {

	case "EQUALS":
		return answers[cond["key"].(string)] == cond["value"]

	case "GT":
		return toFloat(answers[cond["key"].(string)]) > toFloat(cond["value"])

	case "AND":
		for _, r := range cond["rules"].([]any) {
			if !Evaluate(r.(map[string]any), answers) {
				return false
			}
		}
		return true

	case "OR":
		for _, r := range cond["rules"].([]any) {
			if Evaluate(r.(map[string]any), answers) {
				return true
			}
		}
		return false
	}

	return false
}

func toFloat(v any) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case float64:
		return t
	}
	return 0
}
