// Copyright OpenSSF Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ossf/scorecard-action/dependencydiff"
	"github.com/ossf/scorecard-action/entrypoint"
	"github.com/ossf/scorecard-action/options"
	"github.com/ossf/scorecard-action/signing"
)

const (
	EVENT_PULL_REQUEST = "pull_request"
)

func main() {
	fmt.Println(os.Getenv(options.EnvGithubEventName))
	fmt.Println(os.Getenv(options.EnvGithubRef))
	fmt.Println("****&&&&****^^^^")
	switch os.Getenv(options.EnvGithubEventName) {
	case EVENT_PULL_REQUEST:
		// Run the dependency-diff on pull requests.
		ctx := context.Background()
		err := dependencydiff.RunDependencyDiff(ctx)
		if err != nil {
			log.Fatalf("error running dependency-diff: %v", err)
		}
	default:
		// Run the root Scorecard-action.
		action, err := entrypoint.New()
		if err != nil {
			log.Fatalf("creating scorecard entrypoint: %v", err)
		}

		if err := action.Execute(); err != nil {
			log.Fatalf("error during command execution: %v", err)
		}

		if os.Getenv(options.EnvInputPublishResults) == "true" {
			// Get json results by re-running scorecard.
			jsonPayload, err := signing.GetJSONScorecardResults()
			if err != nil {
				log.Fatalf("error generating json scorecard results: %v", err)
			}

			// Sign json results.
			if err = signing.SignScorecardResult("results.json"); err != nil {
				log.Fatalf("error signing scorecard json results: %v", err)
			}

			// Processes json results.
			repoName := os.Getenv(options.EnvGithubRepository)
			repoRef := os.Getenv(options.EnvGithubRef)
			accessToken := os.Getenv(options.EnvInputRepoToken)
			if err := signing.ProcessSignature(jsonPayload, repoName, repoRef, accessToken); err != nil {
				log.Fatalf("error processing signature: %v", err)
			}
		}
	}

}
