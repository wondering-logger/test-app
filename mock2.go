// TestListServerlessAppReports_Success tests successful retrieval of paginated reports
func TestListServerlessAppReports_Success(t *testing.T) {
	// Create mock instances
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	// Create test data - let's create 5 deployments for testing pagination
	testDeployments := []*database.DeploymentDto{
		{
			ID:            "deployment-1",
			ComponentName: "component-1",
			DeploymentEnvironment: database.DeploymentEnvironmentDto{
				Environment:    enums.EnvironmentDv,
				SubEnvironment: sharedtypes.SubEnvironment("subenv-1"),
				Region:         enums.RegionEast,
				Cloud:          enums.CloudProviderAzure,
			},
			Namespace:   "namespace-1",
			ClusterName: "cluster-1",
			NetworkZone: enums.NetworkZoneTrusted,
		},
		{
			ID:            "deployment-2",
			ComponentName: "component-2",
			DeploymentEnvironment: database.DeploymentEnvironmentDto{
				Environment:    enums.EnvironmentPd,
				SubEnvironment: sharedtypes.SubEnvironment("subenv-2"),
				Region:         enums.RegionWest,
				Cloud:          enums.CloudProviderOnPrem,
			},
			Namespace:   "namespace-2",
			ClusterName: "cluster-2",
			NetworkZone: enums.NetworkZoneSemiTrusted,
		},
		{
			ID:            "deployment-3",
			ComponentName: "component-3",
			DeploymentEnvironment: database.DeploymentEnvironmentDto{
				Environment:    enums.EnvironmentDv,
				SubEnvironment: sharedtypes.SubEnvironment("subenv-3"),
				Region:         enums.RegionCentral,
				Cloud:          enums.CloudProviderAzure,
			},
			Namespace:   "namespace-3",
			ClusterName: "cluster-3",
			NetworkZone: enums.NetworkZoneUntrusted,
		},
	}

	// Create corresponding components
	testComponents := map[string]*database.ComponentDto{
		"component-1": {
			Name:           "component-1",
			SystemId:       sharedtypes.SystemId("system-1"),
			DeploymentType: enums.DeploymentTypeACS,
			ComponentType:  enums.ComponentTypeWeb,
		},
		"component-2": {
			Name:           "component-2",
			SystemId:       sharedtypes.SystemId("system-2"),
			DeploymentType: enums.DeploymentTypeFaaS,
			ComponentType:  enums.ComponentTypeWorker,
		},
		"component-3": {
			Name:           "component-3",
			SystemId:       sharedtypes.SystemId("system-3"),
			DeploymentType: enums.DeploymentTypeACS,
			ComponentType:  enums.ComponentTypeCronJob,
		},
	}

	// Set up mocks
	mockDeploymentRepo.On("List", mock.Anything).Return(testDeployments, nil)

	// Mock component repo to return the corresponding component for each deployment
	for componentName, component := range testComponents {
		mockComponentRepo.On("Get", mock.Anything, componentName).Return(component, nil)
	}

	// Test case 1: Get first page with 2 items per page
	ctx := context.Background()
	reportList, err := rm.ListServerlessAppReports(ctx, 1, 2)

	// Assert results for first page
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, reportList, "Expected report list to not be nil")
	assert.Equal(t, 3, reportList.TotalCount, "Expected total count to be 3")
	assert.Equal(t, 1, reportList.PageNumber, "Expected page number to be 1")
	assert.Equal(t, 2, reportList.ItemsPerPage, "Expected items per page to be 2")
	assert.Len(t, reportList.Items, 2, "Expected 2 items on first page")

	// Verify first page items
	assert.Equal(t, "deployment-1", reportList.Items[0].ID)
	assert.Equal(t, "deployment-2", reportList.Items[1].ID)

	// Test case 2: Get second page
	reportList2, err := rm.ListServerlessAppReports(ctx, 2, 2)

	assert.NoError(t, err, "Expected no error for second page")
	assert.NotNil(t, reportList2, "Expected report list to not be nil")
	assert.Equal(t, 3, reportList2.TotalCount, "Expected total count to be 3")
	assert.Equal(t, 2, reportList2.PageNumber, "Expected page number to be 2")
	assert.Equal(t, 2, reportList2.ItemsPerPage, "Expected items per page to be 2")
	assert.Len(t, reportList2.Items, 1, "Expected 1 item on second page")
	assert.Equal(t, "deployment-3", reportList2.Items[0].ID)

	// Verify all mock expectations were met
	mockDeploymentRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

// TestListServerlessAppReports_EmptyList tests when there are no deployments
func TestListServerlessAppReports_EmptyList(t *testing.T) {
	// Create mock instances
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	// Set up mock to return empty list
	mockDeploymentRepo.On("List", mock.Anything).Return([]*database.DeploymentDto{}, nil)

	// Call the method
	ctx := context.Background()
	reportList, err := rm.ListServerlessAppReports(ctx, 1, 10)

	// Assert results
	assert.NoError(t, err, "Expected no error for empty list")
	assert.NotNil(t, reportList, "Expected report list to not be nil")
	assert.Equal(t, 0, reportList.TotalCount, "Expected total count to be 0")
	assert.Equal(t, 1, reportList.PageNumber, "Expected page number to be 1")
	assert.Equal(t, 10, reportList.ItemsPerPage, "Expected items per page to be 10")
	assert.Len(t, reportList.Items, 0, "Expected 0 items")

	// Verify mocks
	mockDeploymentRepo.AssertExpectations(t)
	mockComponentRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

// TestListServerlessAppReports_InvalidPagination tests pagination edge cases
func TestListServerlessAppReports_InvalidPagination(t *testing.T) {
	// Create mock instances
	mockComponentRepo := new(mocks.MockComponentDocRepo)
	mockDeploymentRepo := new(mocks.MockDeploymentDocRepo)

	// Create the report manager
	rm := NewReportManager(mockComponentRepo, mockDeploymentRepo)

	// Create test data
	testDeployments := []*database.DeploymentDto{
		{
			ID:            "deployment-1",
			ComponentName: "component-1",
			DeploymentEnvironment: database.DeploymentEnvironmentDto{
				Environment:    enums.EnvironmentDv,
				SubEnvironment: sharedtypes.SubEnvironment("subenv-1"),
				Region:         enums.RegionEast,
				Cloud:          enums.CloudProviderAzure,
			},
			Namespace:   "namespace-1",
			ClusterName: "cluster-1",
			NetworkZone: enums.NetworkZoneTrusted,
		},
	}

	testComponent := &database.ComponentDto{
		Name:           "component-1",
		SystemId:       sharedtypes.SystemId("system-1"),
		DeploymentType: enums.DeploymentTypeACS,
		ComponentType:  enums.ComponentTypeWeb,
	}

	// Set up mocks
	mockDeploymentRepo.On("List", mock.Anything).Return(testDeployments, nil)
	mockComponentRepo.On("Get", mock.Anything, "component-1").Return(testComponent, nil)

	// Test case 1: Invalid page number (0 or negative)
	ctx := context.Background()
	reportList, err := rm.ListServerlessAppReports(ctx, 0, 10)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 1, reportList.PageNumber, "Expected page number to default to 1")

	// Test case 2: Invalid items per page (0 or negative)
	reportList2, err := rm.ListServerlessAppReports(ctx, 1, -5)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 10, reportList2.ItemsPerPage, "Expected items per page to default to 10")

	// Test case 3: Page number beyond available data
	reportList3, err := rm.ListServerlessAppReports(ctx, 5, 10)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 1, reportList3.TotalCount, "Expected total count to be 1")
	assert.Len(t, reportList3.Items, 0, "Expected 0 items when page is beyond available data")

	// Verify mocks
	mockDeploymentRepo.AssertExpectations(t)
}
