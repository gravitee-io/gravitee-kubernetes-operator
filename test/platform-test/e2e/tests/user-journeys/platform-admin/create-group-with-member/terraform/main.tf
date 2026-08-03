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

# A group created through the Terraform APIM provider. Lands in APIM via the
# Automation API (origin KUBERNETES).
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

variable "group_name" {
  type    = string
  default = "simple-group-tf"
}

# Gates the resource so the journey's remove("group") can drop it from the
# desired state and re-apply, the way a user deletes a resource block.
variable "create_group" {
  type    = bool
  default = true
}

resource "apim_group" "group" {
  count           = var.create_group ? 1 : 0
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "simple-group-tf"
  name            = var.group_name
  notify_members  = false
}

output "group_id" {
  value = var.create_group ? apim_group.group[0].id : ""
}
