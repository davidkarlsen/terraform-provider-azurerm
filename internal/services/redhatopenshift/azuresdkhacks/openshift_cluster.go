// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

// This azuresdkhack exists because the go-azure-sdk 2023-09-04 openshiftclusters package does not yet include
// PlatformWorkloadIdentityProfile support (GA'd in ARO February 2026).
// TODO: remove this hack once the upstream SDK is updated to include PlatformWorkloadIdentityProfile.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04/openshiftclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// OpenShiftClustersClient wraps the SDK client to add PlatformWorkloadIdentityProfile support.
type OpenShiftClustersClient struct {
	client *openshiftclusters.OpenShiftClustersClient
}

// NewOpenShiftClustersWorkaroundClient creates a new workaround client.
func NewOpenShiftClustersWorkaroundClient(c *openshiftclusters.OpenShiftClustersClient) OpenShiftClustersClient {
	return OpenShiftClustersClient{client: c}
}

// OpenShiftCluster extends the SDK type to include PlatformWorkloadIdentityProfile.
type OpenShiftCluster struct {
	Id         *string                     `json:"id,omitempty"`
	Location   string                      `json:"location"`
	Name       *string                     `json:"name,omitempty"`
	Properties *OpenShiftClusterProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData      `json:"systemData,omitempty"`
	Tags       *map[string]string          `json:"tags,omitempty"`
	Type       *string                     `json:"type,omitempty"`
}

// OpenShiftClusterProperties extends the SDK type to include PlatformWorkloadIdentityProfile.
type OpenShiftClusterProperties struct {
	ApiserverProfile                *openshiftclusters.APIServerProfile        `json:"apiserverProfile,omitempty"`
	ClusterProfile                  *openshiftclusters.ClusterProfile          `json:"clusterProfile,omitempty"`
	ConsoleProfile                  *openshiftclusters.ConsoleProfile          `json:"consoleProfile,omitempty"`
	IngressProfiles                 *[]openshiftclusters.IngressProfile        `json:"ingressProfiles,omitempty"`
	MasterProfile                   *openshiftclusters.MasterProfile           `json:"masterProfile,omitempty"`
	NetworkProfile                  *openshiftclusters.NetworkProfile          `json:"networkProfile,omitempty"`
	PlatformWorkloadIdentityProfile *PlatformWorkloadIdentityProfile           `json:"platformWorkloadIdentityProfile,omitempty"`
	ProvisioningState               *openshiftclusters.ProvisioningState       `json:"provisioningState,omitempty"`
	ServicePrincipalProfile         *openshiftclusters.ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	WorkerProfiles                  *[]openshiftclusters.WorkerProfile         `json:"workerProfiles,omitempty"`
	WorkerProfilesStatus            *[]openshiftclusters.WorkerProfile         `json:"workerProfilesStatus,omitempty"`
}

// OpenShiftClusterUpdate extends the SDK type to include PlatformWorkloadIdentityProfile.
type OpenShiftClusterUpdate struct {
	Properties *OpenShiftClusterUpdateProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData            `json:"systemData,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
}

// OpenShiftClusterUpdateProperties mirrors OpenShiftClusterProperties for PATCH operations.
type OpenShiftClusterUpdateProperties struct {
	PlatformWorkloadIdentityProfile *PlatformWorkloadIdentityProfile           `json:"platformWorkloadIdentityProfile,omitempty"`
	ServicePrincipalProfile         *openshiftclusters.ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
}

// PlatformWorkloadIdentityProfile holds the map of named platform workload identities.
type PlatformWorkloadIdentityProfile struct {
	PlatformWorkloadIdentities *map[string]PlatformWorkloadIdentity `json:"platformWorkloadIdentities,omitempty"`
}

// PlatformWorkloadIdentity holds a single user-assigned managed identity resource ID.
type PlatformWorkloadIdentity struct {
	ResourceId *string `json:"resourceId,omitempty"`
}

// CreateOrUpdateOperationResponse is the response type for CreateOrUpdate.
type CreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// GetOperationResponse holds the response for a Get operation.
type GetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *OpenShiftCluster
}

// UpdateOperationResponse is the response type for Update.
type UpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// CreateOrUpdate performs a PUT for the cluster.
func (c OpenShiftClustersClient) CreateOrUpdate(ctx context.Context, id openshiftclusters.ProviderOpenShiftClusterId, input OpenShiftCluster) (result CreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.client.Client)
	if err != nil {
		return
	}

	return
}

// CreateOrUpdateThenPoll performs CreateOrUpdate then polls until it's completed.
func (c OpenShiftClustersClient) CreateOrUpdateThenPoll(ctx context.Context, id openshiftclusters.ProviderOpenShiftClusterId, input OpenShiftCluster) error {
	result, err := c.CreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CreateOrUpdate: %+v", err)
	}

	return nil
}

// Get retrieves an OpenShiftCluster with the extended properties.
func (c OpenShiftClustersClient) Get(ctx context.Context, id openshiftclusters.ProviderOpenShiftClusterId) (result GetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       id.ID(),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model OpenShiftCluster
	if err = resp.Unmarshal(&model); err != nil {
		return
	}

	result.Model = &model
	return
}

// Update performs a PATCH for the cluster.
func (c OpenShiftClustersClient) Update(ctx context.Context, id openshiftclusters.ProviderOpenShiftClusterId, input OpenShiftClusterUpdate) (result UpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       id.ID(),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.client.Client)
	if err != nil {
		return
	}

	return
}

// UpdateThenPoll performs Update then polls until it's completed.
func (c OpenShiftClustersClient) UpdateThenPoll(ctx context.Context, id openshiftclusters.ProviderOpenShiftClusterId, input OpenShiftClusterUpdate) error {
	result, err := c.Update(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Update: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Update: %+v", err)
	}

	return nil
}
