package utils

import "net/http"

func GetPathParameterValues(r *http.Request, parameters []string) map[string]string {
	parameterValues := map[string]string{}
	for _, parameter := range parameters {
		parameterValues[parameter] = r.PathValue(parameter)
	}

	return parameterValues
}
