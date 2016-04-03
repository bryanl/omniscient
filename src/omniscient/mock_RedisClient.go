package omniscient

import "github.com/stretchr/testify/mock"

import "time"

type MockRedisClient struct {
	mock.Mock
}

func (_m *MockRedisClient) Delete(keys ...string) (int64, error) {
	ret := _m.Called(keys)

	var r0 int64
	if rf, ok := ret.Get(0).(func(...string) int64); ok {
		r0 = rf(keys...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...string) error); ok {
		r1 = rf(keys...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) HGetAllMap(key string) (map[string]string, error) {
	ret := _m.Called(key)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(string) map[string]string); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) HMSet(key string, field string, value string, pairs ...string) (string, error) {
	ret := _m.Called(key, field, value, pairs)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, ...string) string); ok {
		r0 = rf(key, field, value, pairs...)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, ...string) error); ok {
		r1 = rf(key, field, value, pairs...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) LPush(key string, values ...string) (int64, error) {
	ret := _m.Called(key, values)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, ...string) int64); ok {
		r0 = rf(key, values...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = rf(key, values...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) LRange(key string, start int64, stop int64) ([]string, error) {
	ret := _m.Called(key, start, stop)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, int64, int64) []string); ok {
		r0 = rf(key, start, stop)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64, int64) error); ok {
		r1 = rf(key, start, stop)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) LRem(key string, count int64, value interface{}) (int64, error) {
	ret := _m.Called(key, count, value)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, int64, interface{}) int64); ok {
		r0 = rf(key, count, value)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64, interface{}) error); ok {
		r1 = rf(key, count, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) Ping() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	ret := _m.Called(key, value, expiration)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, interface{}, time.Duration) string); ok {
		r0 = rf(key, value, expiration)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}, time.Duration) error); ok {
		r1 = rf(key, value, expiration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
