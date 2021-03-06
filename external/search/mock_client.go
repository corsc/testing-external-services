package search

import (
	"github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// Search provides a mock function with given fields: term
func (_m *MockClient) Search(term string) (Results, error) {
	ret := _m.Called(term)

	var r0 Results
	if rf, ok := ret.Get(0).(func(string) Results); ok {
		r0 = rf(term)
	} else {
		r0 = ret.Get(0).(Results)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(term)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
