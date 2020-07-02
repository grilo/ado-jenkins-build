package main

import (
	"log"
	"os"
	"strings"
)

func glabTranslate(variable string) string {

	switch variable {

	case "BUILD_REPOSITORY_URI":
		newName := "gitlabSourceRepoURL"
		log.Printf("Translating %s into: %s", variable, newName)
		return newName
	case "BUILD_SOURCEBRANCH":
		newName := "gitlabSourceBranch"
		log.Printf("Translating %s into: %s", variable, newName)
		return newName
	}

	return variable
}

func ReadEnvironment(emulateGitLab bool) map[string]string {

	/*
	   See: https://docs.microsoft.com/en-us/azure/devops/pipelines/build/variables?view=azure-devops&tabs=yaml

	   A variable such as Build.SourceBranchName will be exposed as BUILD_SOURCEBRANCHNAME.
	*/

	vars := make(map[string]string)

	//BUILD_REPOSITORY_URI
	//BUILD_SOURCEBRANCHNAME

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "AGENT_") || strings.HasPrefix(key, "BUILD_") || strings.HasPrefix(key, "PIPELINE_") || strings.HasPrefix(key, "SYSTEM_") {
			log.Printf("Found %s: %s", key, value)
			if emulateGitLab {
				vars[glabTranslate(key)] = value
			} else {
				vars[key] = value
			}
		}
	}

	return vars
}
