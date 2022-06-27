// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fluentbit

import (
	"fmt"

	"github.com/GoogleCloudPlatform/ops-agent/internal/port"
)

const DefaultMetricsPort = 20202

const MetricsPortKey = "fluent-bit-metrics"

func MetricsInputComponent() Component {
	return Component{
		Kind: "INPUT",
		Config: map[string]string{
			"Name":            "fluentbit_metrics",
			"Scrape_On_Start": "True",
			"Scrape_Interval": "60",
		},
	}
}

func MetricsOutputComponent() Component {
	return Component{
		Kind: "OUTPUT",
		Config: map[string]string{
			// https://docs.fluentbit.io/manual/pipeline/outputs/prometheus-exporter
			"Name":  "prometheus_exporter",
			"Match": "*",
			"host":  "0.0.0.0",
			"port":  fmt.Sprintf("%d", MetricsPort()),
		},
	}
}

func MetricsPort() uint16 {
	var metricsPort uint16 = DefaultMetricsPort
	config, err := port.ReadConfig(port.DefaultConfigPath)
	if err == nil && config != nil {
		if reservedPort, ok := config.ReservedPorts[MetricsPortKey]; ok {
			metricsPort = reservedPort
		} else {
			chooser, err := port.NewRandomPortChooser()
			if err == nil {
				randomPort, err := chooser.Choose()
				if err == nil {
					metricsPort = randomPort
				}
			}
		}
	}
	return metricsPort
}
