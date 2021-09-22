// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mockNewStore

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// NewStore is an autogenerated mock type for the NewStore type
type NewStore struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *NewStore) Create(_a0 context.Context, _a1 []byte) ([]byte, uuid.UUID, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, []byte) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 uuid.UUID
	if rf, ok := ret.Get(1).(func(context.Context, []byte) uuid.UUID); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(uuid.UUID)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, []byte) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *NewStore) Delete(_a0 context.Context, _a1 uuid.UUID) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *NewStore) Get(_a0 context.Context, _a1 uuid.UUID) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: _a0
func (_m *NewStore) GetAll(_a0 context.Context) ([][]byte, error) {
	ret := _m.Called(_a0)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(context.Context) [][]byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: _a0, _a1
func (_m *NewStore) GetUser(_a0 context.Context, _a1 uuid.UUID) ([][]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) [][]byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0, _a1, _a2
func (_m *NewStore) Update(_a0 context.Context, _a1 uuid.UUID, _a2 []byte) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, []byte) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, []byte) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}