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

# GKO-2564: apim_group resource lifecycle (create / update / idempotency /
# destroy). Parameterized via variables so a single workspace can be driven
# through multiple apply phases via *.auto.tfvars.json (see helpers/terraform.ts
# writeVars).
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  # Configured via environment variables:
  #   APIM_SERVER_URL  (e.g. http://localhost:30083/automation)
  #   APIM_USERNAME
  #   APIM_PASSWORD
}

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
  default = "lifecycle"
}

variable "group_name" {
  type    = string
  default = "e2e-tf-group"
}

variable "notify_members" {
  type    = bool
  default = true
}

resource "apim_group" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-group-${var.hrid_suffix}"
  name            = var.group_name
  notify_members  = var.notify_members
}

output "group_id" {
  value = apim_group.test.id
}

output "group_hrid" {
  value = apim_group.test.hrid
}

output "group_name" {
  value = apim_group.test.name
}
