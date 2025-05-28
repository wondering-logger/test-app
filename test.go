package managers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"serverless_platform/internal/serverless_controller/models"
	"serverless_platform/internal/serverless_controller/registry/database"
	"serverless_platform/internal/shared/enums"
	sharedtypes "serverless_platform/internal/shared/types"
	"serverless_platform/mocks"
)

// TestGetServerlessAppReport_Success tests the successful retrieval of a serverless app report
func TestGetServerlessAppReport_Success(t *testing.T) {
	// Create mock instances for our dependencies
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager with our mocks
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	// Define test data - creating a sample deployment
	testDeploymentID := "test-deployment-123"
	testComponentName := "test-component"
	testSystemID := sharedtypes.SystemId("test-system-id")

	// Create a test deployment DTO
	testDeployment := &database.DeploymentDto{
		ID:            testDeploymentID,
		ComponentName: testComponentName,
		DeploymentEnvironment: database.DeploymentEnvironmentDto{
			Environment:    enums.EnvironmentDv,
			SubEnvironment: sharedtypes.SubEnvironment("test-subenv"),
			Region:         enums.RegionEast,
			Cloud:          enums.CloudProviderAzure,
		},
		Namespace:   "test-namespace",
		ClusterName: "test-cluster",
		NetworkZone: enums.NetworkZoneTrusted,
	}

	// Create a test component DTO
	testComponent := &database.ComponentDto{
		Name:           testComponentName,
		SystemId:       testSystemID,
		DeploymentType: enums.DeploymentTypeACS,
		ComponentType:  enums.ComponentTypeWeb,
		Owners:         []string{"test-owner@example.com"},
	}

	// Set up mock expectations
	// First, we expect the deployment repo to be called with the deployment ID
	mockDeploymentRepo.On("Get", mock.Anything, testDeploymentID).Return(testDeployment, nil)

	// Then, we expect the component repo to be called with the component name
	mockComponentRepo.On("Get", mock.Anything, testComponentName).Return(testComponent, nil)

	// Call the method we're testing
	ctx := context.Background()
	report, err := rm.GetServerlessAppReport(ctx, testDeploymentID)

	// Assert the results
	assert.NoError(t, err, "Expected no error when retrieving report")
	assert.NotNil(t, report, "Expected report to not be nil")

	// Verify the report contains the expected data
	assert.Equal(t, testDeploymentID, report.ID)
	assert.Equal(t, "Dv", report.Environment) // Note: the enum.String() returns "Dv"
	assert.Equal(t, "test-subenv", report.SubEnvironment)
	assert.Equal(t, "Azure", report.CloudProvider) // Note: the enum.String() returns "Azure"
	assert.Equal(t, "East", report.Region) // Note: the enum.String() returns "East"
	assert.Equal(t, "test-cluster", report.Cluster)
	assert.Equal(t, "test-system-id", report.SystemEntryId)
	assert.Equal(t, "Active", report.Status)

	// Verify the report has data items
	assert.Len(t, report.Data, 2, "Expected 2 data items in the report")

	// Check the data items
	componentNameFound := false
	namespaceFound := false
	for _, dataItem := range report.Data {
		if dataItem.Key == "componentName" {
			assert.Equal(t, testComponentName, dataItem.Value)
			componentNameFound = true
		}
		if dataItem.Key == "namespace" {
			assert.Equal(t, "test-namespace", dataItem.Value)
			namespaceFound = true
		}
	}
	assert.True(t, componentNameFound, "Expected to find componentName in data items")
	assert.True(t, namespaceFound, "Expected to find namespace in data items")

	// Verify that the ReportGeneratedOn timestamp is recent (within the last minute)
	assert.WithinDuration(t, time.Now(), report.ReportGeneratedOn, time.Minute)

	// Verify that all mock expectations were met
	mockDeploymentRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

// TestGetServerlessAppReport_DeploymentNotFound tests when the deployment is not found
func TestGetServerlessAppReport_DeploymentNotFound(t *testing.T) {
	// Create mock instances
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	testDeploymentID := "non-existent-deployment"

	// Set up mock to return nil deployment (not found)
	mockDeploymentRepo.On("Get", mock.Anything, testDeploymentID).Return(nil, nil)

	// Call the method
	ctx := context.Background()
	report, err := rm.GetServerlessAppReport(ctx, testDeploymentID)

	// Assert the results
	assert.Error(t, err, "Expected error when deployment not found")
	assert.Nil(t, report, "Expected report to be nil")
	assert.Contains(t, err.Error(), "deployment not found")

	// Verify mocks
	mockDeploymentRepo.AssertExpectations(t)
	mockComponentRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

// TestGetServerlessAppReport_DeploymentRepoError tests when the deployment repo returns an error
func TestGetServerlessAppReport_DeploymentRepoError(t *testing.T) {
	// Create mock instances
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	testDeploymentID := "test-deployment-123"
	expectedError := errors.New("database connection error")

	// Set up mock to return an error
	mockDeploymentRepo.On("Get", mock.Anything, testDeploymentID).Return(nil, expectedError)

	// Call the method
	ctx := context.Background()
	report, err := rm.GetServerlessAppReport(ctx, testDeploymentID)

	// Assert the results
	assert.Error(t, err, "Expected error when deployment repo fails")
	assert.Nil(t, report, "Expected report to be nil")
	assert.Contains(t, err.Error(), "error retrieving deployment")
	assert.Contains(t, err.Error(), "database connection error")

	// Verify mocks
	mockDeploymentRepo.AssertExpectations(t)
	mockComponentRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}
