// TestGetServerlessAppReport_ComponentRepoError tests when the component repo returns an error
func TestGetServerlessAppReport_ComponentRepoError(t *testing.T) {
	// Create mock instances
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	testDeploymentID := "test-deployment-123"
	testComponentName := "test-component"

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

	expectedError := errors.New("component repository error")

	// Set up mocks
	// First call succeeds
	mockDeploymentRepo.On("Get", mock.Anything, testDeploymentID).Return(testDeployment, nil)
	// Second call fails
	mockComponentRepo.On("Get", mock.Anything, testComponentName).Return(nil, expectedError)

	// Call the method
	ctx := context.Background()
	report, err := rm.GetServerlessAppReport(ctx, testDeploymentID)

	// Assert the results
	assert.Error(t, err, "Expected error when component repo fails")
	assert.Nil(t, report, "Expected report to be nil")
	assert.Contains(t, err.Error(), "error retrieving component")
	assert.Contains(t, err.Error(), "component repository error")

	// Verify mocks
	mockDeploymentRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}
