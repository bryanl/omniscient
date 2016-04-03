package omniscient

import "github.com/stretchr/testify/mock"

type MockNoteRepository struct {
	mock.Mock
}

func (_m *MockNoteRepository) Create(content string) (*Note, error) {
	ret := _m.Called(content)

	var r0 *Note
	if rf, ok := ret.Get(0).(func(string) *Note); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Note)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(content)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockNoteRepository) Retrieve(id string) (*Note, error) {
	ret := _m.Called(id)

	var r0 *Note
	if rf, ok := ret.Get(0).(func(string) *Note); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Note)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockNoteRepository) Update(id string, content string) (*Note, error) {
	ret := _m.Called(id, content)

	var r0 *Note
	if rf, ok := ret.Get(0).(func(string, string) *Note); ok {
		r0 = rf(id, content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Note)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, content)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockNoteRepository) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *MockNoteRepository) List() ([]Note, error) {
	ret := _m.Called()

	var r0 []Note
	if rf, ok := ret.Get(0).(func() []Note); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Note)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
