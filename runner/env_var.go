package runner

import "bytes"

// EnvVar wraps a environment variable sent to the container engine
type EnvVar struct {
	Key   string
	Value string
}

// EnvVars models a collection of EnvVar
type EnvVars []EnvVar

func (envVar *EnvVar) String() string {
	return envVar.Key + "=" + envVar.Value
}

func (envVars EnvVars) String() string {
	var buffer bytes.Buffer
	for _, envVar := range envVars {
		buffer.WriteString(envVar.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
