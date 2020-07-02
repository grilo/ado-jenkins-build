package main

import (
	"log"
	"os"
	"strings"
)

func ReadEnvironment() map[string]string {

	/*
	   See: https://docs.microsoft.com/en-us/azure/devops/pipelines/build/variables?view=azure-devops&tabs=yaml

	   A variable such as Build.SourceBranchName will be exposed as BUILD_SOURCEBRANCHNAME.
	*/

	vars := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "AGENT_") || strings.HasPrefix(key, "BUILD_") || strings.HasPrefix(key, "PIPELINE_") || strings.HasPrefix(key, "SYSTEM_") {
			log.Printf("Found %s: %s", key, value)
            vars[key] = value
		}
	}

	return vars
}
