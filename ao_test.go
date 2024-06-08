package aogo

import (
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMU is a mock implementation of the IMU interface
type MockMU struct {
	mock.Mock
}

func (m *MockMU) SendMessage(process string, data string, tags []types.Tag, anchor string, s *signer.Signer) (string, error) {
	args := m.Called(process, data, tags, anchor, s)
	return args.String(0), args.Error(1)
}

func (m *MockMU) SpawnProcess(module string, data string, tags []types.Tag, s *signer.Signer) (string, error) {
	args := m.Called(module, data, tags, s)
	return args.String(0), args.Error(1)
}

func (m *MockMU) Monitor() {
	m.Called()
}

// MockCU is a mock implementation of the ICU interface
type MockCU struct {
	mock.Mock
}

func (m *MockCU) LoadResult(process string, message string) (*Response, error) {
	args := m.Called(process, message)
	return args.Get(0).(*Response), args.Error(1)
}

func (m *MockCU) DryRun(message Message) (*Response, error) {
	args := m.Called(message)
	return args.Get(0).(*Response), args.Error(1)
}

func TestNewAO(t *testing.T) {
	ao, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, ao)
	assert.NotNil(t, ao.mu)
	assert.NotNil(t, ao.cu)
}

func TestNewAOWithCustomMU(t *testing.T) {
	customURL := "https://custom-mu.url"
	ao, err := New(WithMU(customURL))
	assert.NoError(t, err)
	assert.Equal(t, customURL, ao.mu.(*MU).url)
}

func TestNewAOWithCustomCU(t *testing.T) {
	customURL := "https://custom-cu.url"
	ao, err := New(WithCU(customURL))
	assert.NoError(t, err)
	assert.Equal(t, customURL, ao.cu.(*CU).url)
}

func TestSpawnProcess(t *testing.T) {
	mockMU := new(MockMU)
	ao := &AO{mu: mockMU, cu: newCU(CU_URL)}
	mockSigner := &signer.Signer{}
	mockTags := []types.Tag{}
	mockMU.On("SpawnProcess", "testModule", "testData", mockTags, mockSigner).Return("processID", nil)

	processID, err := ao.SpawnProcess("testModule", "testData", mockTags, mockSigner)
	assert.NoError(t, err)
	assert.Equal(t, "processID", processID)
	mockMU.AssertExpectations(t)
}

func TestSendMessage(t *testing.T) {
	mockMU := new(MockMU)
	ao := &AO{mu: mockMU, cu: newCU(CU_URL)}
	mockSigner := &signer.Signer{}
	mockTags := []types.Tag{}
	mockMU.On("SendMessage", "testProcess", "testData", mockTags, "", mockSigner).Return("messageID", nil)

	messageID, err := ao.SendMessage("testProcess", "testData", mockTags, "", mockSigner)
	assert.NoError(t, err)
	assert.Equal(t, "messageID", messageID)
	mockMU.AssertExpectations(t)
}

func TestLoadResult(t *testing.T) {
	mockCU := new(MockCU)
	ao := &AO{mu: newMU(MU_URL), cu: mockCU}
	mockResponse := &Response{}
	mockCU.On("LoadResult", "testProcess", "testMessage").Return(mockResponse, nil)

	response, err := ao.LoadResult("testProcess", "testMessage")
	assert.NoError(t, err)
	assert.Equal(t, mockResponse, response)
	mockCU.AssertExpectations(t)
}
