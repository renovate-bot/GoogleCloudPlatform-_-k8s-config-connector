#!/bin/bash
# Copyright 2025 Google LLC
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


set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd ${REPO_ROOT}/dev/tools/controllerbuilder

go run . generate-types \
  --service google.cloud.gkehub.v1beta --api-version gkehub.cnrm.cloud.google.com/v1beta1 \
  --resource GKEHubFeatureMembership:MembershipFeatureSpec

# NOTYET - not yet using proto
# go run . generate-mapper \
#   --service google.cloud.gkehub.v1beta --api-version gkehub.cnrm.cloud.google.com/v1beta1

# NOTYET - not yet following full pattern
rm ${REPO_ROOT}/apis/gkehub/v1beta1/membershipfeaturespec_identity.go
rm ${REPO_ROOT}/apis/gkehub/v1beta1/membershipfeaturespec_reference.go
rm ${REPO_ROOT}/apis/gkehub/v1beta1/membershipfeaturespec_types.go

cd ${REPO_ROOT}
dev/tasks/generate-crds

go run -mod=readonly golang.org/x/tools/cmd/goimports@latest -w  pkg/controller/direct/gkehub/

