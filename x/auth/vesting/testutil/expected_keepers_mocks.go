// Code generated by MockGen. DO NOT EDIT.
// Source: x/auth/vesting/types/expected_keepers.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
)

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// BlockedAddr mocks base method.
func (m *MockBankKeeper) BlockedAddr(addr types.AccAddress) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockedAddr", addr)
	ret0, _ := ret[0].(bool)
	return ret0
}

// BlockedAddr indicates an expected call of BlockedAddr.
func (mr *MockBankKeeperMockRecorder) BlockedAddr(addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockedAddr", reflect.TypeOf((*MockBankKeeper)(nil).BlockedAddr), addr)
}

// IsSendEnabledCoins mocks base method.
func (m *MockBankKeeper) IsSendEnabledCoins(ctx context.Context, coins ...types.Coin) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range coins {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsSendEnabledCoins", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsSendEnabledCoins indicates an expected call of IsSendEnabledCoins.
func (mr *MockBankKeeperMockRecorder) IsSendEnabledCoins(ctx interface{}, coins ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, coins...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSendEnabledCoins", reflect.TypeOf((*MockBankKeeper)(nil).IsSendEnabledCoins), varargs...)
}

// SendCoins mocks base method.
func (m *MockBankKeeper) SendCoins(ctx context.Context, fromAddr, toAddr types.AccAddress, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoins", ctx, fromAddr, toAddr, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoins indicates an expected call of SendCoins.
func (mr *MockBankKeeperMockRecorder) SendCoins(ctx, fromAddr, toAddr, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoins", reflect.TypeOf((*MockBankKeeper)(nil).SendCoins), ctx, fromAddr, toAddr, amt)
}

// MockAccountsModKeeper is a mock of AccountsModKeeper interface.
type MockAccountsModKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountsModKeeperMockRecorder
}

// MockAccountsModKeeperMockRecorder is the mock recorder for MockAccountsModKeeper.
type MockAccountsModKeeperMockRecorder struct {
	mock *MockAccountsModKeeper
}

// NewMockAccountsModKeeper creates a new mock instance.
func NewMockAccountsModKeeper(ctrl *gomock.Controller) *MockAccountsModKeeper {
	mock := &MockAccountsModKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountsModKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountsModKeeper) EXPECT() *MockAccountsModKeeperMockRecorder {
	return m.recorder
}

// IsAccountsModuleAccount mocks base method.
func (m *MockAccountsModKeeper) IsAccountsModuleAccount(ctx context.Context, accountAddr []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAccountsModuleAccount", ctx, accountAddr)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAccountsModuleAccount indicates an expected call of IsAccountsModuleAccount.
func (mr *MockAccountsModKeeperMockRecorder) IsAccountsModuleAccount(ctx, accountAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAccountsModuleAccount", reflect.TypeOf((*MockAccountsModKeeper)(nil).IsAccountsModuleAccount), ctx, accountAddr)
}

// SendModuleMessageUntyped mocks base method.
func (m *MockAccountsModKeeper) SendModuleMessageUntyped(ctx context.Context, sender []byte, msg protoiface.MessageV1) (protoiface.MessageV1, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendModuleMessageUntyped", ctx, sender, msg)
	ret0, _ := ret[0].(protoiface.MessageV1)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendModuleMessageUntyped indicates an expected call of SendModuleMessageUntyped.
func (mr *MockAccountsModKeeperMockRecorder) SendModuleMessageUntyped(ctx, sender, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendModuleMessageUntyped", reflect.TypeOf((*MockAccountsModKeeper)(nil).SendModuleMessageUntyped), ctx, sender, msg)
}

// CurrentAccountNumber mocks base method.
func (m *MockAccountsModKeeper) CurrentAccountNumber(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentAccountNumber", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrentAccountNumber indicates an expected call of CurrentAccountNumber.
func (mr *MockAccountsModKeeperMockRecorder) ç(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentAccountNumber", reflect.TypeOf((*MockAccountsModKeeper)(nil).CurrentAccountNumber), ctx)
}

// NextAccountNumber mocks base method.
func (m *MockAccountsModKeeper) NextAccountNumber(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextAccountNumber", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextAccountNumber indicates an expected call of NextAccountNumber.
func (mr *MockAccountsModKeeperMockRecorder) NextAccountNumber(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextAccountNumber", reflect.TypeOf((*MockAccountsModKeeper)(nil).NextAccountNumber), ctx)
}
