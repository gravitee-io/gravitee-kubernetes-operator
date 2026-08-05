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

# An mTLS-secured API plus a subscribed application whose client certificates are
# driven by `cert_mode`, so the journey can issue, rotate and revoke them with
# re-applies. Certificates are read from the shared PKI the gateway calls are made
# with, so both arms present the SAME certificate to the gateway.
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

# Absolute path, passed by the scenario: the workspace is copied to a temp dir, so
# a path relative to this file would not resolve.
variable "pki_dir" {
  type = string
}

# Distinguishes the two scenarios sharing this fixture; APIM names are a global
# namespace across the serial suite.
variable "resource_prefix" {
  type    = string
  default = "mtls-client-cert"
}

# client1 | both | client2 | none | deprecated
variable "cert_mode" {
  type    = string
  default = "client1"
}

locals {
  client1 = file("${var.pki_dir}/client1.crt")
  client2 = file("${var.pki_dir}/client2.crt")

  # Mirrors the GKO arm's application-*.yaml stages.
  cert_sets = {
    client1 = [
      {
        name      = "client1"
        content   = local.client1
        starts_at = "2026-01-01T00:00:00Z"
        ends_at   = "2030-01-01T00:00:00Z"
      }
    ]
    both = [
      {
        name      = "client1"
        content   = local.client1
        starts_at = "2026-01-01T00:00:00Z"
        ends_at   = "2030-01-01T00:00:00Z"
      },
      {
        name      = "client2"
        content   = local.client2
        starts_at = null
        ends_at   = null
      }
    ]
    client2 = [
      {
        name      = "client2"
        content   = local.client2
        starts_at = null
        ends_at   = null
      }
    ]
    none       = []
    deprecated = []
  }

  is_deprecated = var.cert_mode == "deprecated"
  name          = var.resource_prefix
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "${local.name}-tf"
  name            = "${local.name}-tf"
  description     = "V4 proxy API whose only plan is secured with mTLS"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  listeners = [
    {
      http = {
        type        = "HTTP"
        paths       = [{ path = "/${local.name}-tf/" }]
        entrypoints = [{ type = "http-proxy" }]
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

  plans = [
    {
      hrid       = "mtls-plan"
      name       = "mtls"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security   = { type = "MTLS" }
    }
  ]
}

resource "apim_application" "app" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "${local.name}-app-tf"
  name            = "${local.name}-app-tf"
  description     = "Application authenticating to an mTLS plan with a client certificate"

  settings = {
    app = {
      type      = "SIMPLE"
      client_id = "${local.name}-app-tf-client"
    }
    # The two forms are mutually exclusive in the provider, so the unused one is
    # null rather than absent: an object built by a conditional must keep a
    # consistent attribute set.
    tls = {
      client_certificate  = local.is_deprecated ? local.client1 : null
      client_certificates = local.is_deprecated ? null : local.cert_sets[var.cert_mode]
    }
  }
}

resource "apim_subscription" "sub" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = "${local.name}-sub-tf"
  api_hrid         = apim_apiv4.api.hrid
  application_hrid = apim_application.app.hrid
  plan_hrid        = "mtls-plan"
}

output "api_id" {
  value = apim_apiv4.api.id
}

output "app_id" {
  value = apim_application.app.id
}

output "sub_id" {
  value = apim_subscription.sub.id
}

output "api_context_path" {
  value = "/${local.name}-tf"
}
