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

# Create and update a Shared Policy Group through the Terraform APIM provider.
# `updated` swaps the step and the description, mirroring the two GKO fixture
# files. Only `hrid` forces replacement on this resource, so re-applying with
# updated = true changes the same shared policy group in place.
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

variable "updated" {
  type    = bool
  default = false
}

resource "apim_shared_policy_group" "spg" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "update-spg-tf"
  name            = "update-spg-tf"
  description     = var.updated ? "Shared policy group after its update" : "Shared policy group as first authored"
  api_type        = "PROXY"
  phase           = "REQUEST"
  steps = [
    {
      enabled = true
      name    = var.updated ? "Inject the updated tracking header" : "Inject the tracking header"
      policy  = "transform-headers"
      configuration = jsonencode({
        addHeaders = [
          { name = "X-SPG-Test", value = var.updated ? "spg-header-updated" : "spg-header" }
        ]
      })
    }
  ]
}

output "spg_id" {
  value = apim_shared_policy_group.spg.id
}
