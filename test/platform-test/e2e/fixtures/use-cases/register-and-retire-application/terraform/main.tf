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

# Use case: register an application through the Terraform APIM provider. The
# Terraform arm of the register-and-retire-application journey. `description` is
# re-applied to exercise the update; `create_application = false` drops the
# resource so a re-apply retires (archives) it the way a user would.
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

variable "description" {
  type    = string
  default = "E2E use-case: application registered"
}

# Toggle the application on/off so a test can retire it the way a user would:
# drop the resource from the desired state and `terraform apply`.
variable "create_application" {
  type    = bool
  default = true
}

resource "apim_application" "app" {
  count           = var.create_application ? 1 : 0
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-uc-app-lifecycle"
  name            = "e2e-tf-uc-app-lifecycle"
  description     = var.description
  settings = {
    app = {
      type      = "SIMPLE"
      client_id = "e2e-tf-uc-app-lifecycle-client"
    }
  }
}

output "app_id" {
  value = one(apim_application.app[*].id)
}
