package dolaterio

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

// StringArray returns an array of strings to send as docker Env
func (envVars EnvVars) StringArray() []string {
	res := make([]string, len(envVars))
	for idx, envVar := range envVars {
		res[idx] = envVar.String()
	}
	return res
}
