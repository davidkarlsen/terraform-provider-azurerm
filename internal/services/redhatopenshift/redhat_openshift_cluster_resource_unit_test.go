// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04/openshiftclusters"
)

func TestExpandOpenshiftPlatformWorkloadIdentityProfile(t *testing.T) {
	cases := []struct {
		name     string
		input    []PlatformWorkloadIdentityProfile
		expected *openshiftclusters.PlatformWorkloadIdentityProfile
	}{
		{
			name:     "empty input returns nil",
			input:    []PlatformWorkloadIdentityProfile{},
			expected: nil,
		},
		{
			name: "single identity",
			input: []PlatformWorkloadIdentityProfile{
				{
					PlatformWorkloadIdentities: []PlatformWorkloadIdentity{
						{
							Name:       "cloud-controller-manager",
							ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myIdentity",
						},
					},
				},
			},
			expected: &openshiftclusters.PlatformWorkloadIdentityProfile{
				PlatformWorkloadIdentities: &map[string]openshiftclusters.PlatformWorkloadIdentity{
					"cloud-controller-manager": {
						ResourceId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myIdentity"),
					},
				},
			},
		},
		{
			name: "multiple identities",
			input: []PlatformWorkloadIdentityProfile{
				{
					PlatformWorkloadIdentities: []PlatformWorkloadIdentity{
						{
							Name:       "cloud-controller-manager",
							ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ccm",
						},
						{
							Name:       "ingress",
							ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ingress",
						},
					},
				},
			},
			expected: &openshiftclusters.PlatformWorkloadIdentityProfile{
				PlatformWorkloadIdentities: &map[string]openshiftclusters.PlatformWorkloadIdentity{
					"cloud-controller-manager": {
						ResourceId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ccm"),
					},
					"ingress": {
						ResourceId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ingress"),
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandOpenshiftPlatformWorkloadIdentityProfile(tc.input)
			if tc.expected == nil {
				if result != nil {
					t.Fatalf("expected nil, got %+v", result)
				}
				return
			}
			if result == nil {
				t.Fatalf("expected non-nil result, got nil")
			}
			if result.PlatformWorkloadIdentities == nil {
				t.Fatalf("expected PlatformWorkloadIdentities to be non-nil")
			}
			expected := *tc.expected.PlatformWorkloadIdentities
			actual := *result.PlatformWorkloadIdentities
			if len(actual) != len(expected) {
				t.Fatalf("expected %d identities, got %d", len(expected), len(actual))
			}
			for name, expectedIdentity := range expected {
				actualIdentity, ok := actual[name]
				if !ok {
					t.Fatalf("expected identity %q not found in result", name)
				}
				if pointer.From(actualIdentity.ResourceId) != pointer.From(expectedIdentity.ResourceId) {
					t.Fatalf("expected ResourceId %q for identity %q, got %q",
						pointer.From(expectedIdentity.ResourceId), name, pointer.From(actualIdentity.ResourceId))
				}
			}
		})
	}
}

func TestFlattenOpenShiftPlatformWorkloadIdentityProfile(t *testing.T) {
	cases := []struct {
		name     string
		input    *openshiftclusters.PlatformWorkloadIdentityProfile
		expected []PlatformWorkloadIdentityProfile
	}{
		{
			name:     "nil input returns empty slice",
			input:    nil,
			expected: []PlatformWorkloadIdentityProfile{},
		},
		{
			name: "single identity",
			input: &openshiftclusters.PlatformWorkloadIdentityProfile{
				PlatformWorkloadIdentities: &map[string]openshiftclusters.PlatformWorkloadIdentity{
					"cloud-controller-manager": {
						ResourceId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ccm"),
					},
				},
			},
			expected: []PlatformWorkloadIdentityProfile{
				{
					PlatformWorkloadIdentities: []PlatformWorkloadIdentity{
						{
							Name:       "cloud-controller-manager",
							ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ccm",
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := flattenOpenShiftPlatformWorkloadIdentityProfile(tc.input)
			if len(result) != len(tc.expected) {
				t.Fatalf("expected %d profiles, got %d", len(tc.expected), len(result))
			}
			if len(tc.expected) == 0 {
				return
			}
			actualIdentities := result[0].PlatformWorkloadIdentities
			expectedIdentities := tc.expected[0].PlatformWorkloadIdentities
			if len(actualIdentities) != len(expectedIdentities) {
				t.Fatalf("expected %d identities, got %d", len(expectedIdentities), len(actualIdentities))
			}
			// Build a map for comparison since order from map iteration is not guaranteed
			actualMap := make(map[string]PlatformWorkloadIdentity)
			for _, identity := range actualIdentities {
				actualMap[identity.Name] = identity
			}
			for _, expectedIdentity := range expectedIdentities {
				actualIdentity, ok := actualMap[expectedIdentity.Name]
				if !ok {
					t.Fatalf("expected identity %q not found in result", expectedIdentity.Name)
				}
				if actualIdentity.ResourceId != expectedIdentity.ResourceId {
					t.Fatalf("expected ResourceId %q for identity %q, got %q",
						expectedIdentity.ResourceId, expectedIdentity.Name, actualIdentity.ResourceId)
				}
			}
		})
	}
}
