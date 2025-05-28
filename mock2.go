package mocks
// go test ./internal/serverless_controller/managers -v -run TestGetServerlessAppReport
import (
	"context"
	"serverless_platform/internal/serverless_controller/registry/database"

	"github.com/stretchr/testify/mock"
)

// MockDeploymentDocRepo is a mock implementation of DeploymentDocRepo
type MockDeploymentDocRepo struct {
	mock.Mock
}

// Get mocks the Get method of DeploymentDocRepo
func (m *MockDeploymentDocRepo) Get(ctx context.Context, id string) (*database.DeploymentDto, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*database.DeploymentDto), args.Error(1)
}

// Create mocks the Create method of DeploymentDocRepo
func (m *MockDeploymentDocRepo) Create(ctx context.Context, dto *database.DeploymentDto) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

// Update mocks the Update method of DeploymentDocRepo
func (m *MockDeploymentDocRepo) Update(ctx context.Context, dto *database.DeploymentDto) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

// List mocks the List method of DeploymentDocRepo
func (m *MockDeploymentDocRepo) List(ctx context.Context) ([]*database.DeploymentDto, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*database.DeploymentDto), args.Error(1)
}

// ListForComponent mocks the ListForComponent method of DeploymentDocRepo
func (m *MockDeploymentDocRepo) ListForComponent(ctx context.Context, componentName string) ([]database.DeploymentDto, error) {
	args := m.Called(ctx, componentName)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]database.DeploymentDto), args.Error(1)
}

// Delete mocks the Delete method of DeploymentDocRepo
func (m *MockDeploymentDocRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
