package processes_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"serverless_platform/cmd/serverless_controller/configs"
	"serverless_platform/internal/serverless_controller/processes"
	"serverless_platform/internal/serverless_controller/services/naas"
	naasModels "serverless_platform/internal/serverless_controller/services/naas/models"
	"serverless_platform/internal/shared/enums"
	sharedModels "serverless_platform/internal/shared/models"
)

// Test 1: Happy Path Deploy - Namespace doesn't exist, gets created
func TestNaaSProcess_HappyPathDeploy_NamespaceCreated(t *testing.T) {
	// Setup
	mockNonProdClient := &MockNaaSClient{}
	naasProcess := createTestNaaSProcess(mockNonProdClient, mockNonProdClient)

	// Create ALM request with Deploy action
	almRequest := createBaseALMRequestWithAction(enums.ActionDeploy)
	almRequest.DeploymentActivity.Containers = []sharedModels.ContainerData{
		createTestContainer("test-container", "100m", "256Mi", "200m", "512Mi"),
	}

	// Mock GetNamespace to return "not found"
	mockNonProdClient.On("GetNamespace", mock.Anything, "test-system-id", mock.Anything).
		Return(nil, errors.New("namespace not found"))

	// Mock CreateNamespace to succeed
	expectedCreateRequest := mock.MatchedBy(func(req naasModels.NaaS_Request) bool {
		return req.Metadata.Name != "" &&
			   req.Spec.SystemID == "test-system-id" &&
			   req.Spec.ResourceQuota["requests.cpu"] == "280m" && // 100m * 1.4 * 2
			   req.Spec.ResourceQuota["requests.memory"] == "716Mi" // 256Mi * 1.4 * 2
	})

	mockNonProdClient.On("CreateNamespace", mock.Anything, expectedCreateRequest).
		Return(&naasModels.NaaS_Response{
			Status: naasModels.NaaS_Status{
				ClusterName:   "test-cluster-001",
				NamespaceName: "test-namespace",
				IsCreated:     true,
			},
			NaaS_Request: naasModels.NaaS_Request{
				Metadata: naasModels.NaaS_Request_Metadata{
					Name: "test-namespace",
				},
			},
		}, nil)

	// Execute
	response, err := naasProcess.CreateOrUpdateNaaSResources(context.Background(), almRequest)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "test-cluster-001", response.ClusterName)
	assert.Equal(t, "test-namespace", response.NamespaceName)

	// Verify method calls
	mockNonProdClient.AssertCalled(t, "GetNamespace", mock.Anything, "test-system-id", mock.Anything)
	mockNonProdClient.AssertCalled(t, "CreateNamespace", mock.Anything, mock.Anything)
	mockNonProdClient.AssertNotCalled(t, "UpdateNamespace", mock.Anything, mock.Anything)
}

// Test 2: Delete Action with Valid Namespace - Should skip update call
func TestNaaSProcess_DeleteAction_ValidNamespace_SkipsUpdate(t *testing.T) {
	// Setup
	mockNonProdClient := &MockNaaSClient{}
	naasProcess := createTestNaaSProcess(mockNonProdClient, mockNonProdClient)

	// Create ALM request with Delete action
	almRequest := createBaseALMRequestWithAction(enums.ActionDelete)
	almRequest.DeploymentActivity.Containers = []sharedModels.ContainerData{
		createTestContainer("test-container", "100m", "256Mi", "200m", "512Mi"),
	}

	// Mock GetNamespace to return existing namespace
	existingNamespace := &naasModels.NaaS_Response{
		Status: naasModels.NaaS_Status{
			ClusterName:   "existing-cluster-001",
			NamespaceName: "existing-namespace",
			IsCreated:     true,
		},
		NaaS_Request: naasModels.NaaS_Request{
			Metadata: naasModels.NaaS_Request_Metadata{
				Name: "existing-namespace",
			},
			Spec: naasModels.NaaS_Spec{
				ResourceQuota: map[string]string{
					"requests.cpu":    "100m", // Different from what would be calculated
					"requests.memory": "200Mi",
				},
			},
		},
	}

	mockNonProdClient.On("GetNamespace", mock.Anything, "test-system-id", mock.Anything).
		Return(existingNamespace, nil)

	// Execute
	response, err := naasProcess.CreateOrUpdateNaaSResources(context.Background(), almRequest)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "existing-cluster-001", response.ClusterName)
	assert.Equal(t, "existing-namespace", response.NamespaceName)

	// Verify method calls - should NOT call CreateNamespace or UpdateNamespace
	mockNonProdClient.AssertCalled(t, "GetNamespace", mock.Anything, "test-system-id", mock.Anything)
	mockNonProdClient.AssertNotCalled(t, "CreateNamespace", mock.Anything, mock.Anything)
	mockNonProdClient.AssertNotCalled(t, "UpdateNamespace", mock.Anything, mock.Anything)
}

// Test 3: Delete Action with Invalid Namespace - Should return error
func TestNaaSProcess_DeleteAction_InvalidNamespace_ReturnsError(t *testing.T) {
	// Setup
	mockNonProdClient := &MockNaaSClient{}
	naasProcess := createTestNaaSProcess(mockNonProdClient, mockNonProdClient)

	// Create ALM request with Delete action
	almRequest := createBaseALMRequestWithAction(enums.ActionDelete)
	almRequest.DeploymentActivity.Containers = []sharedModels.ContainerData{
		createTestContainer("test-container", "100m", "256Mi", "200m", "512Mi"),
	}

	// Mock GetNamespace to return "not found"
	mockNonProdClient.On("GetNamespace", mock.Anything, "test-system-id", mock.Anything).
		Return(nil, errors.New("namespace not found"))

	// Execute
	response, err := naasProcess.CreateOrUpdateNaaSResources(context.Background(), almRequest)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available for delete action")
	assert.Empty(t, response.ClusterName)
	assert.Empty(t, response.NamespaceName)

	// Verify method calls - should NOT call CreateNamespace or UpdateNamespace
	mockNonProdClient.AssertCalled(t, "GetNamespace", mock.Anything, "test-system-id", mock.Anything)
	mockNonProdClient.AssertNotCalled(t, "CreateNamespace", mock.Anything, mock.Anything)
	mockNonProdClient.AssertNotCalled(t, "UpdateNamespace", mock.Anything, mock.Anything)
}

// Helper functions
func createTestNaaSProcess(prodClient, nonProdClient *MockNaaSClient) processes.NaaSProcess {
	mockConfig := configs.Config{
		NaaSNamespaceAdminGroups: []string{"admin-group-1", "admin-group-2"},
	}

	// You'll need to modify NewNaaSProcess or create a test version that accepts mock clients
	return processes.NewNaaSProcessForTesting(mockConfig, prodClient, nonProdClient)
}

func createBaseALMRequestWithAction(action enums.Action) sharedModels.ALMRequest {
	return sharedModels.ALMRequest{
		ActivityMetaData: sharedModels.ActivityMetaData{
			Action: action,
		},
		ComponentMetadata: sharedModels.ComponentMetadata{
			ComponentName:  "test-component",
			Environment:    enums.EnvironmentDv,
			SubEnvironment: "test-subenv",
			Region:         enums.RegionEast,
			SystemEntryId:  "test-system-id",
			NetworkZone:    enums.NetworkZoneTrusted,
			CloudProvider:  enums.CloudProviderAzure,
		},
		DeploymentActivity: sharedModels.DeploymentActivity{
			Containers: []sharedModels.ContainerData{}, // Will be filled by test cases
		},
	}
}

func createTestContainer(name, reqCpu, reqMem, limitCpu, limitMem string) sharedModels.ContainerData {
	return sharedModels.ContainerData{
		Name: name,
		Resources: sharedModels.Resources{
			Requests: sharedModels.ResourceValues{
				Cpu:    reqCpu,
				Memory: reqMem,
			},
			Limits: sharedModels.ResourceValues{
				Cpu:    limitCpu,
				Memory: limitMem,
			},
		},
	}
}
