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

# A native Kafka API whose connection-metrics reporter is driven by a variable,
# so the journey can turn it back on with a re-apply.
#
# The listener host and the plan's broker port range are unique to this arm:
# APIM admission enforces global uniqueness on both, so this fixture and its GKO
# counterpart must not share either.
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

variable "reporter_metrics_enabled" {
  type    = bool
  default = false
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "observability-native-api-tf"
  name            = "observability-native-api-tf"
  description     = "Native Kafka API whose connection-metrics reporter is toggled"
  version         = "1.0"
  type            = "NATIVE"
  state           = "STOPPED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  analytics = {
    enabled                  = true
    reporter_metrics_enabled = var.reporter_metrics_enabled
  }

  listeners = [
    {
      kafka = {
        type = "KAFKA"
        host = "kafka-observability-tf.example.local"
        port = 9092
        entrypoints = [
          { type = "native-kafka" }
        ]
      }
    }
  ]

  endpoint_groups = [
    {
      name = "Default Kafka group"
      type = "native-kafka"
      endpoints = [
        {
          name                  = "Default Kafka endpoint"
          type                  = "native-kafka"
          inherit_configuration = false
          configuration = jsonencode({
            bootstrapServers = "kafka-observability-tf.example.local:9092"
          })
          shared_configuration_override = jsonencode({
            security = { protocol = "PLAINTEXT" }
          })
        }
      ]
    }
  ]

  plans = [
    {
      hrid               = "keyless"
      name               = "Free plan"
      type               = "API"
      mode               = "STANDARD"
      validation         = "AUTO"
      status             = "PUBLISHED"
      bootstrap_port     = 9492
      broker_range_start = 9500
      broker_range_end   = 9502
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

output "api_id" {
  value = apim_apiv4.api.id
}
