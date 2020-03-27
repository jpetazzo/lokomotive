// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build aws packet
// +build poste2e

package monitoring

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	testutil "github.com/kinvolk/lokomotive/test/components/util"
)

//nolint:funlen
func testComponentsPrometheusMetrics(t *testing.T, v1api v1.API) {
	testCases := []struct {
		componentName string
		query         string
		platforms     []testutil.Platform
	}{
		{
			componentName: "kube-apiserver",
			query:         "apiserver_request_total",
		},
		{
			componentName: "coredns",
			query:         "coredns_build_info",
		},
		{
			componentName: "kube-scheduler",
			query:         "scheduler_schedule_attempts_total",
		},
		{
			componentName: "kube-controller-manager",
			query:         "workqueue_work_duration_seconds_bucket",
		},
		{
			componentName: "kube-proxy",
			query:         "kubeproxy_sync_proxy_rules_duration_seconds_bucket",
		},
		{
			componentName: "kubelet",
			query:         "kubelet_running_pod_count",
		},
		{
			componentName: "metallb",
			query:         "metallb_bgp_session_up",
			platforms:     []testutil.Platform{testutil.PlatformPacket},
		},
		{
			componentName: "contour",
			query:         "contour_dagrebuild_timestamp",
			platforms:     []testutil.Platform{testutil.PlatformPacket, testutil.PlatformAWS},
		},
		{
			componentName: "cert-manager",
			query:         "certmanager_controller_sync_call_count",
			platforms:     []testutil.Platform{testutil.PlatformPacket, testutil.PlatformAWS},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("prometheus-%s", tc.componentName), func(t *testing.T) {
			if !testutil.IsPlatformSupported(t, tc.platforms) {
				t.Skip()
			}

			t.Parallel()

			t.Logf("querying %q", tc.query)

			const contextTimeout = 10

			var err error

			// This loop ensures that we try to query prometheus for the state of targets multiple times.
			// This is the retry logic.
			for i := 0; i < 20; i++ {
				t.Logf("Running scrape test iteration #%d", i)

				// Use function to be able to use defer.
				err = func() error {
					ctx, cancel := context.WithTimeout(context.Background(), contextTimeout*time.Second)
					defer cancel()

					results, warnings, err := v1api.Query(ctx, tc.query, time.Now())
					if err != nil {
						return fmt.Errorf("error querying Prometheus: %w", err)
					}

					if len(warnings) > 0 {
						t.Logf("warnings: %v", warnings)
					}

					if len(results.String()) == 0 {
						return fmt.Errorf("no metrics found")
					}

					t.Logf("found %d results for %s", len(strings.Split(results.String(), "\n")), tc.query)

					return nil
				}()

				// If there is no errors, break the retry loop.
				if err == nil {
					break
				}

				// Wait a bit before next attempt.
				time.Sleep(30 * time.Second) //nolint:gomnd
			}

			// If there are still some errors after all retries, fail the test.
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
