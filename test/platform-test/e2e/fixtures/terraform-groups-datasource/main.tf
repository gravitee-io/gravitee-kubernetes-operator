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

# GKO-2564: apim_group data source. A group is created by the resource, then
# read back through a `data "apim_group"` block keyed on its hrid — the data
# source only resolves groups created via the Automation API (hrid lookup;
# UUID/name lookup is pending API-side work).
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

variable "hrid_suffix" {
  type    = string
  default = "datasource"
}

resource "apim_group" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-group-${var.hrid_suffix}"
  name            = "e2e-tf-group-${var.hrid_suffix}"
  notify_members  = false
}

# Reading the data source by hrid pulls the group's attributes back through
# the provider's read path. depends_on defers the read until after the
# resource is applied, so the lookup never runs against a not-yet-created
# group.
data "apim_group" "lookup" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = apim_group.test.hrid

  depends_on = [apim_group.test]
}

output "resource_id" {
  value = apim_group.test.id
}

output "resource_name" {
  value = apim_group.test.name
}

output "ds_id" {
  value = data.apim_group.lookup.id
}

output "ds_name" {
  value = data.apim_group.lookup.name
}

output "ds_hrid" {
  value = data.apim_group.lookup.hrid
}

output "ds_notify_members" {
  value = data.apim_group.lookup.notify_members
}
