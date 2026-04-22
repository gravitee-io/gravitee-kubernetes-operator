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

# GKO-1380: invalid subscription configuration in Terraform. The subscription
# references a non-existent api_hrid/plan_hrid pair so `terraform apply` fails
# with a clear error message from the APIM provider.
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

resource "apim_application" "e2e_tf_1380_app" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-1380-app"
  name            = "e2e-tf-1380-app"
  description     = "E2E test: TF subscription error handling — application"
  settings = {
    app = {
      type      = "SIMPLE"
      client_id = "e2e-tf-1380-client"
    }
  }
}

# Intentionally invalid: the api_hrid + plan_hrid pair does not exist, so the
# APIM provider rejects the apply with a 404/validation error.
resource "apim_subscription" "e2e_tf_1380_sub" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = "e2e-tf-1380-sub"
  api_hrid         = "does-not-exist-api"
  plan_hrid        = "does-not-exist-plan"
  application_hrid = apim_application.e2e_tf_1380_app.hrid
}
