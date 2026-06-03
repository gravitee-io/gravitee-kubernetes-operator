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

# GKO-1383 (TF variant): create an application via Terraform; the test will
# subsequently `terraform destroy` and verify the application is gone in APIM.
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

resource "apim_application" "e2e_tf_1383_app" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-1383-app"
  name            = "e2e-tf-1383-app"
  description     = "E2E test: TF-based application delete lifecycle"
  settings = {
    app = {
      type      = "SIMPLE"
      client_id = "e2e-tf-1383-client"
    }
  }
}

output "app_id" {
  value = apim_application.e2e_tf_1383_app.id
}
