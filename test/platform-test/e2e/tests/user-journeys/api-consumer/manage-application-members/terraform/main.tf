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

# An application whose membership is driven by variables, so the journey can
# grant, re-role and revoke it with re-applies. Members are inline on
# apim_application; there is no standalone membership resource.
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

variable "member_source_id" {
  type    = string
  default = "e2e-sa-app-member"
}

variable "with_member" {
  type    = bool
  default = true
}

variable "member_role" {
  type    = string
  default = "USER"
}

variable "notify_members" {
  type    = bool
  default = false
}

resource "apim_application" "app" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "shared-application-tf"
  name            = "shared-application-tf"
  description     = "Application whose membership the journey grants, changes and revokes"
  notify_members  = var.notify_members

  members = var.with_member ? [
    {
      source    = "gravitee"
      source_id = var.member_source_id
      role      = var.member_role
    }
  ] : []

  settings = {
    app = {
      type      = "SIMPLE"
      client_id = "shared-application-tf-client"
    }
  }
}

output "app_id" {
  value = apim_application.app.id
}
