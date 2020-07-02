package main

import (
	"log"
	"net/url"
	"os"

	"github.com/akamensky/argparse"
)

var (
	version        = "development"
	commit         = "not_git_repo"
	buildTimestamp = "no_build_time"
)

func main() {
	log.Printf("Initializing ado-jenkins-build-%s (%s)", version, commit)
	log.Printf("Build timestamp: %s", buildTimestamp)

	parser := argparse.NewParser("ado-jenkins-build", "Executes a jenkins job from Azure DevOps and waits for its completion. Returns 0 on SUCCESS, 1 on FAILURE, 2 on ABORTED, 3 on UNSTABLE and 4 on NOT_BUILT.")

	argEmulateGitLab := parser.Flag("g", "gitlab", &argparse.Options{Help: "Emulate gitlab variables (useful when migrating from GitLab)."})
	argTimeout := parser.Int("t", "timeout", &argparse.Options{
		Help:    "Maximum number of seconds to wait (poll).",
		Default: 3600, // One hour
	})
	argUrl := parser.String("u", "url", &argparse.Options{
		Required: true,
		Help:     "Alternative URL for the manifest.plist file.",
	})

	parseErr := parser.Parse(os.Args)
	if parseErr != nil {
		log.Fatal(parser.Usage(parseErr))
	}

	if *argEmulateGitLab {
		log.Printf("GitLab emulation enabled.")
	}

	parsedUrl, err := url.Parse(*argUrl)
	if err != nil {
		log.Fatalf("Unable to parse url: %s", argUrl)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		log.Fatalf("Unknown protocol: %s (%s)", parsedUrl.Scheme, *argUrl)
	}

	variables := ReadEnvironment(*argEmulateGitLab)

	response := TriggerBuild(parsedUrl, variables)

	returnCode := WaitForBuild(response, *argTimeout)
	log.Printf("Exiting with rc: %d", returnCode)
	os.Exit(returnCode)
}
