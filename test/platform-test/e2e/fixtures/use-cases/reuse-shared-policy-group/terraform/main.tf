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

# Use case: reuse a Shared Policy Group across a V4 API, through the Terraform
# APIM provider. apim_shared_policy_group defines the reusable step; the API's
# request flow invokes it via the shared-policy-group-policy step. Setting
# attach_spg = false drops the flow so a re-apply detaches the SPG.
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

variable "attach_spg" {
  type    = bool
  default = true
}

resource "apim_shared_policy_group" "spg" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-uc-spg"
  name            = "e2e-tf-uc-spg"
  description     = "E2E use-case: reusable shared policy group"
  api_type        = "PROXY"
  phase           = "REQUEST"
  steps = [
    {
      enabled = true
      name    = "Transform Headers"
      policy  = "transform-headers"
      configuration = jsonencode({
        addHeaders = [
          { name = "X-SPG-Test", value = "spg-header" }
        ]
      })
    }
  ]
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-uc-spg-api"
  name            = "e2e-tf-uc-spg-api"
  description     = "E2E use-case: V4 API reusing a shared policy group"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/e2e-tf-uc-spg-api/" }
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
          name                  = "default-endpoint"
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

  # The API reuses the SPG via a shared-policy-group-policy step, referencing it
  # by HRID. attach_spg = false drops the flow to detach the SPG.
  flows = var.attach_spg ? [
    {
      enabled = true
      name    = "Flow with SPG"
      selectors = [
        {
          http = {
            type          = "HTTP"
            path          = "/"
            path_operator = "STARTS_WITH"
          }
        }
      ]
      request = [
        {
          enabled = true
          name    = "SPG step"
          policy  = "shared-policy-group-policy"
          configuration = jsonencode({
            sharedPolicyGroupId = "{#sharedPolicyGroup['${apim_shared_policy_group.spg.hrid}']}"
          })
        }
      ]
    }
  ] : []

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
  value = "/e2e-tf-uc-spg-api"
}
