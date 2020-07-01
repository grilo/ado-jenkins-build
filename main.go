package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
)

var (
	version        = "development"
	commit         = "not_git_repo"
	buildTimestamp = "no_build_time"
)

func main() {
	log.Printf("Initializing jenkins-integration-%s (%s)", version, commit)
	log.Printf("Build timestamp: %s", buildTimestamp)

	parser := argparse.NewParser("webapp", "Launches deployment service.")

	gitlab := parser.Flag("g", "gitlab", &argparse.Options{Help: "Emulate gitlab variables."})
	url := parser.String("u", "url", &argparse.Options{
		Default: "https://jenkins-spain.ic.ing.net",
		Help:    "Alternative URL for the manifest.plist file.",
	})

	parseErr := parser.Parse(os.Args)
	if parseErr != nil {
		log.Fatal(parser.Usage(parseErr))
	}

	if *gitlab {
		log.Printf("GitLab emulation enabled.")
	}

	log.Printf("Using URL: %s", *url)

	Get(*url)
}
