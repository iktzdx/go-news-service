//nolint:forcetypeassert,wrapcheck
package posts

import (
	"github.com/stretchr/testify/mock"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type MockBoundaryRepoPort struct {
	mock.Mock
}

func (mockRepo *MockBoundaryRepoPort) FindPostByID(id int64) (storage.Data, error) {
	args := mockRepo.Called(id)

	return args.Get(0).(storage.Data), args.Error(1)
}

func (mockRepo *MockBoundaryRepoPort) List(opts storage.SearchOpts) (storage.BulkData, error) {
	args := mockRepo.Called(opts)

	return args.Get(0).(storage.BulkData), args.Error(1)
}

func (mockRepo *MockBoundaryRepoPort) Create(data storage.Data) (int64, error) {
	args := mockRepo.Called(data)

	return args.Get(0).(int64), args.Error(1)
}

func (mockRepo *MockBoundaryRepoPort) Update(change storage.Data) (int64, error) {
	args := mockRepo.Called(change)

	return args.Get(0).(int64), args.Error(1)
}

func (mockRepo *MockBoundaryRepoPort) Delete(id int64) (int64, error) {
	args := mockRepo.Called(id)

	return args.Get(0).(int64), args.Error(1)
}
