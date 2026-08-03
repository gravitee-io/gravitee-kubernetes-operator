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

# A V4 proxy API whose membership is driven by variables, so the journey can
# grant, re-role and revoke it with re-applies. Members are inline on apim_apiv4;
# there is no standalone membership resource.
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
  default = "e2e-sa-api-member"
}

# Whether the API grants the member at all — false revokes the membership.
variable "with_member" {
  type    = bool
  default = true
}

# Empty means "declare the member without a role" and let APIM default it.
variable "member_role" {
  type    = string
  default = ""
}

variable "notify_members" {
  type    = bool
  default = false
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "api-with-members-tf"
  name            = "api-with-members-tf"
  description     = "V4 proxy API whose membership the journey grants, changes and revokes"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"
  notify_members  = var.notify_members

  members = var.with_member ? [
    {
      source    = "gravitee"
      source_id = var.member_source_id
      role      = var.member_role != "" ? var.member_role : null
    }
  ] : []

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/api-with-members-tf/" }
        ]
        entrypoints = [
          { type = "http-proxy" }
        ]
      }
    }
  ]

  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          inherit_configuration = false
          configuration         = jsonencode({ target = "https://api.gravitee.io/echo" })
        }
      ]
    }
  ]

  flow_execution = {
    mode           = "DEFAULT"
    match_required = false
  }

  plans = [
    {
      hrid       = "keyless"
      name       = "Free plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

output "api_id" {
  value = apim_apiv4.api.id
}

output "api_context_path" {
  value = "/api-with-members-tf"
}
