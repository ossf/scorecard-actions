# Copyright 2021 Security Scorecard Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Testing: docker run -e GITHUB_REF=refs/heads/main \
#           -e GITHUB_EVENT_NAME=branch_protection_rule \
#           -e INPUT_RESULTS_FORMAT=sarif \
#           -e INPUT_RESULTS_FILE=results.sarif \
#           -e GITHUB_WORKSPACE=/ \
#           -e INPUT_POLICY_FILE="/policy.yml" \
#           -e INPUT_REPO_TOKEN=$GITHUB_AUTH_TOKEN \
#           -e GITHUB_REPOSITORY="ossf/scorecard" \
#           laurentsimon/scorecard-action:latest
FROM gcr.io/openssf/scorecard:v4.3.1@sha256:gcr.io/openssf/scorecard@sha256:06e3ddde7f63619813c5749389010b596e753fa070c524a42fd0de756f96970f as base

# Build our image and update the root certs.
# TODO: use distroless.
FROM debian:11.3-slim@sha256:06a93cbdd49a265795ef7b24fe374fee670148a7973190fb798e43b3cf7c5d0f
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    jq ca-certificates curl

# Copy the scorecard binary from the official scorecard image.
COPY --from=base /scorecard /scorecard

# Copy a test policy for local testing.
COPY policies/template.yml  /policy.yml

# Our entry point.
# Note: the file is executable in the repo
# and permission carry over to the image.
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
