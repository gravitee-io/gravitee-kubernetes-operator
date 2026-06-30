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

# Use case: stand up a V4 MESSAGE (event) API through the Terraform APIM provider.
# HTTP-GET + webhook subscription entrypoints over a mock message endpoint.
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

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-uc-message"
  name            = "e2e-tf-uc-message"
  description     = "E2E use-case: consume a message API"
  version         = "1"
  type            = "MESSAGE"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/e2e-tf-uc-message/" }
        ]
        entrypoints = [
          { type = "http-get" }
        ]
      }
    },
    {
      subscription = {
        type = "SUBSCRIPTION"
        entrypoints = [
          { type = "webhook" }
        ]
      }
    }
  ]

  endpoint_groups = [
    {
      name = "Default Mock group"
      type = "mock"
      endpoints = [
        {
          name                  = "Default Mock endpoint"
          type                  = "mock"
          inherit_configuration = false
          configuration = jsonencode({
            messageInterval = 1000
            messageContent  = "{\"message\": \"hello from mock\"}"
            messagesCount   = 1
          })
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
      hrid     = "keyless"
      name     = "Free plan"
      type     = "API"
      mode     = "STANDARD"
      status   = "PUBLISHED"
      security = { type = "KEY_LESS" }
    }
  ]
}

output "api_id" {
  value = apim_apiv4.api.id
}

output "api_context_path" {
  value = "/e2e-tf-uc-message"
}
