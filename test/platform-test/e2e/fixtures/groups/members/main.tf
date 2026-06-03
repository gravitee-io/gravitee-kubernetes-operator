# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# GKO-2564: apim_group resource with a `members` list. The member set is
# supplied via *.auto.tfvars.json (see helpers/terraform.ts writeVars) so the
# same workspace can be re-applied with a resolvable member and then with a
# non-resolvable one.
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {}

variable "environment_id" {
  type    = string
  default = "DEFAULT"
}

variable "organization_id" {
  type    = string
  default = "DEFAULT"
}

variable "hrid_suffix" {
  type    = string
  default = "members"
}

variable "members" {
  type = list(object({
    source    = string
    source_id = string
    roles     = map(string)
  }))
  default = []
}

resource "apim_group" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-group-${var.hrid_suffix}"
  name            = "e2e-tf-group-${var.hrid_suffix}"
  notify_members  = false
  members         = var.members
}

output "group_id" {
  value = apim_group.test.id
}

output "group_hrid" {
  value = apim_group.test.hrid
}
