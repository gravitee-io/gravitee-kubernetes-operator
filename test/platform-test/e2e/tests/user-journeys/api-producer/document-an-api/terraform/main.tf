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

# Document a V4 API with an inline markdown page, through the Terraform APIM
# provider. Pages are an inline attribute of apim_apiv4 (no standalone apim_page
# resource). with_page = false re-applies with an empty list to strip the page.
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

variable "page_revised" {
  type    = bool
  default = false
}

variable "with_page" {
  type    = bool
  default = true
}

# Heredocs cannot appear inside a conditional expression, so both revisions of
# the page body are locals the ternary below picks between.
locals {
  initial_content = <<-EOT
    # Getting started

    Call `GET /` to reach the upstream echo endpoint.
  EOT
  revised_content = <<-EOT
    # Quick start

    Send any request to the base path; the upstream echoes it back.
  EOT
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "documented-api-tf"
  name            = "documented-api-tf"
  description     = "V4 proxy API documented with an inline markdown page"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  # The page keeps its hrid across the revision, so a rename has to update the
  # existing page rather than replace it.
  pages = var.with_page ? [
    {
      hrid    = "getting-started"
      name    = var.page_revised ? "Quick start" : "Getting started"
      type    = "MARKDOWN"
      content = var.page_revised ? local.revised_content : local.initial_content
      published  = true
      visibility = var.page_revised ? "PRIVATE" : "PUBLIC"
    }
  ] : []

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/documented-api-tf/" }
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
  value = "/documented-api-tf"
}
