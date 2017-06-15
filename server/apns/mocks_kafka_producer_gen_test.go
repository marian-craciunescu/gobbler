// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/cosminrentea/gobbler/server/kafka (interfaces: Producer)

package apns

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Producer interface
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *_MockProducerRecorder
}

// Recorder for MockProducer (not exported)
type _MockProducerRecorder struct {
	mock *MockProducer
}

func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &_MockProducerRecorder{mock}
	return mock
}

func (_m *MockProducer) EXPECT() *_MockProducerRecorder {
	return _m.recorder
}

func (_m *MockProducer) Report(_param0 string, _param1 []byte, _param2 string) {
	_m.ctrl.Call(_m, "Report", _param0, _param1, _param2)
}

func (_mr *_MockProducerRecorder) Report(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Report", arg0, arg1, arg2)
}

func (_m *MockProducer) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProducerRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

func (_m *MockProducer) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProducerRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}
