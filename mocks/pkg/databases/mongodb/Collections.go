// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	helpers "payment-service/internal/pkg/helpers"

	mock "github.com/stretchr/testify/mock"

	mongodb "payment-service/internal/pkg/databases/mongodb"

	options "go.mongodb.org/mongo-driver/mongo/options"
)

// Collections is an autogenerated mock type for the Collections type
type Collections struct {
	mock.Mock
}

// Aggregate provides a mock function with given fields: payload, ctx
func (_m *Collections) Aggregate(payload mongodb.Aggregate, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for Aggregate")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.Aggregate, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// Close provides a mock function with given fields: ctx
func (_m *Collections) Close(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountData provides a mock function with given fields: payload, ctx
func (_m *Collections) CountData(payload mongodb.CountData, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for CountData")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.CountData, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// FindAllData provides a mock function with given fields: payload, ctx
func (_m *Collections) FindAllData(payload mongodb.FindAllData, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindAllData")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.FindAllData, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// FindMany provides a mock function with given fields: payload, ctx
func (_m *Collections) FindMany(payload mongodb.FindMany, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindMany")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.FindMany, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// FindOne provides a mock function with given fields: payload, ctx
func (_m *Collections) FindOne(payload mongodb.FindOne, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.FindOne, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// FindOneAndUpdate provides a mock function with given fields: payload, rd, ctx
func (_m *Collections) FindOneAndUpdate(payload mongodb.FindOneAndUpdate, rd options.ReturnDocument, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, rd, ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindOneAndUpdate")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.FindOneAndUpdate, options.ReturnDocument, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, rd, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// InsertOne provides a mock function with given fields: payload, ctx
func (_m *Collections) InsertOne(payload mongodb.InsertOne, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for InsertOne")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.InsertOne, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// UpdateOne provides a mock function with given fields: payload, ctx
func (_m *Collections) UpdateOne(payload mongodb.UpdateOne, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.UpdateOne, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// UpsertOne provides a mock function with given fields: payload, ctx
func (_m *Collections) UpsertOne(payload mongodb.UpdateOne, ctx context.Context) <-chan helpers.Result {
	ret := _m.Called(payload, ctx)

	if len(ret) == 0 {
		panic("no return value specified for UpsertOne")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(mongodb.UpdateOne, context.Context) <-chan helpers.Result); ok {
		r0 = rf(payload, ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// NewCollections creates a new instance of Collections. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCollections(t interface {
	mock.TestingT
	Cleanup(func())
}) *Collections {
	mock := &Collections{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
