/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package features

import (
	"testing"

	backendconfigv1 "k8s.io/ingress-gce/pkg/apis/backendconfig/v1"
	"k8s.io/ingress-gce/pkg/composite"
	"k8s.io/ingress-gce/pkg/utils"
)

func TestEnsureCDN(t *testing.T) {
	testCases := []struct {
		desc           string
		sp             utils.ServicePort
		be             *composite.BackendService
		updateExpected bool
	}{
		{
			desc: "cdn setting are missing from spec, no update needed",
			sp: utils.ServicePort{
				BackendConfig: &backendconfigv1.BackendConfig{
					Spec: backendconfigv1.BackendConfigSpec{
						Cdn: nil,
					},
				},
			},
			be: &composite.BackendService{
				EnableCDN: true,
				CdnPolicy: &composite.BackendServiceCdnPolicy{
					CacheKeyPolicy: &composite.CacheKeyPolicy{
						IncludeHost: true,
					},
				},
			},
			updateExpected: false,
		},
		{
			desc: "cache policies are missing from spec, no update needed",
			sp: utils.ServicePort{
				BackendConfig: &backendconfigv1.BackendConfig{
					Spec: backendconfigv1.BackendConfigSpec{
						Cdn: &backendconfigv1.CDNConfig{
							Enabled: true,
						},
					},
				},
			},
			be: &composite.BackendService{
				EnableCDN: true,
				CdnPolicy: &composite.BackendServiceCdnPolicy{
					CacheKeyPolicy: &composite.CacheKeyPolicy{
						IncludeHost: true,
					},
				},
			},
			updateExpected: false,
		},
		{
			desc: "settings are identical, no update needed",
			sp: utils.ServicePort{
				BackendConfig: &backendconfigv1.BackendConfig{
					Spec: backendconfigv1.BackendConfigSpec{
						Cdn: &backendconfigv1.CDNConfig{
							Enabled: true,
							CachePolicy: &backendconfigv1.CacheKeyPolicy{
								IncludeHost: true,
							},
						},
					},
				},
			},
			be: &composite.BackendService{
				EnableCDN: true,
				CdnPolicy: &composite.BackendServiceCdnPolicy{
					CacheKeyPolicy: &composite.CacheKeyPolicy{
						IncludeHost: true,
					},
				},
			},
			updateExpected: false,
		},
		{
			desc: "cache settings are different, update needed",
			sp: utils.ServicePort{
				BackendConfig: &backendconfigv1.BackendConfig{
					Spec: backendconfigv1.BackendConfigSpec{
						Cdn: &backendconfigv1.CDNConfig{
							Enabled: true,
							CachePolicy: &backendconfigv1.CacheKeyPolicy{
								IncludeHost:        true,
								IncludeQueryString: false,
								IncludeProtocol:    false,
							},
						},
					},
				},
			},
			be: &composite.BackendService{
				EnableCDN: true,
				CdnPolicy: &composite.BackendServiceCdnPolicy{
					CacheKeyPolicy: &composite.CacheKeyPolicy{
						IncludeHost:        false,
						IncludeQueryString: true,
						IncludeProtocol:    true,
					},
				},
			},
			updateExpected: true,
		},
		{
			desc: "enabled setting is different, update needed",
			sp: utils.ServicePort{
				BackendConfig: &backendconfigv1.BackendConfig{
					Spec: backendconfigv1.BackendConfigSpec{
						Cdn: &backendconfigv1.CDNConfig{
							Enabled: true,
						},
					},
				},
			},
			be: &composite.BackendService{
				EnableCDN: false,
			},
			updateExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := EnsureCDN(tc.sp, tc.be)
			if result != tc.updateExpected {
				t.Errorf("%v: expected %v but got %v", tc.desc, tc.updateExpected, result)
			}
		})
	}
}
