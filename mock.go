package mocks

import (
	"context"
	"serverless_platform/internal/serverless_controller/registry/database"

	"github.com/stretchr/testify/mock"
)

// MockComponentDocRepo is a mock implementation of ComponentDocRepo
type MockComponentDocRepo struct {
	mock.Mock
}

// Get mocks the Get method of ComponentDocRepo
func (m *MockComponentDocRepo) Get(ctx context.Context, componentName string) (*database.ComponentDto, error) {
	args := m.Called(ctx, componentName)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*database.ComponentDto), args.Error(1)
}

// Create mocks the Create method of ComponentDocRepo
func (m *MockComponentDocRepo) Create(ctx context.Context, dto *database.ComponentDto) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

// Update mocks the Update method of ComponentDocRepo
func (m *MockComponentDocRepo) Update(ctx context.Context, dto *database.ComponentDto) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

// List mocks the List method of ComponentDocRepo
func (m *MockComponentDocRepo) List(ctx context.Context) ([]*database.ComponentDto, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*database.ComponentDto), args.Error(1)
}

// Delete mocks the Delete method of ComponentDocRepo
func (m *MockComponentDocRepo) Delete(ctx context.Context, componentName string) error {
	args := m.Called(ctx, componentName)
	return args.Error(0)
}
