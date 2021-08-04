// Copyright 2021 Google LLC
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

package confgenerator

type MetricsExporterGoogleCloudMonitoring struct {
	ConfigComponent `yaml:",inline"`
}

func (r MetricsExporterGoogleCloudMonitoring) Type() string {
	return "google_cloud_monitoring"
}

func (r *MetricsExporterGoogleCloudMonitoring) ValidateParameters(subagent string, kind string, id string) error {
	panic("Should never be called")
}

func init() {
	metricsExporterTypes.registerType(func() component { return &MetricsExporterGoogleCloudMonitoring{} })
}

type LoggingExporterGoogleCloudLogging struct {
	ConfigComponent `yaml:",inline"`
}

func (r LoggingExporterGoogleCloudLogging) Type() string {
	return "google_cloud_logging"
}

func (r *LoggingExporterGoogleCloudLogging) ValidateParameters(subagent string, kind string, id string) error {
	panic("Should never be called")
}

func init() {
	loggingExporterTypes.registerType(func() component { return &LoggingExporterGoogleCloudLogging{} })
}