package consumer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	envelope "github.com/utilitywarehouse/energy-contracts/pkg/generated"
	"github.com/utilitywarehouse/energy-contracts/pkg/generated/platform"
	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/consumer"
	"github.com/utilitywarehouse/energy-onboarding-lynnluong/internal/models"
	"github.com/uw-labs/substrate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// MockProjector is a mock implementation of the store.Projector interface
type MockProjector struct {
	mock.Mock
}

func (m *MockProjector) Begin() {
	m.Called()
}

func (m *MockProjector) AddToBatch(svc models.EnergyService) error {
	args := m.Called(svc)
	return args.Error(0)
}

func (m *MockProjector) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestConsumer_HandleBatch(t *testing.T) {
	mockProjector := new(MockProjector)

	// Test data
	serviceID := "test-service-id"
	event := &platform.ElectricityServiceRequestReceivedEvent{
		ServiceId: serviceID,
	}

	anyMsg, err := anypb.New(event)
	assert.NoError(t, err)

	envelope := &envelope.Envelope{
		Message: anyMsg,
	}

	msgData, err := proto.Marshal(envelope)
	assert.NoError(t, err)

	// Set up the mock expectations
	mockProjector.On("Begin").Return()
	mockProjector.On("AddToBatch", mock.AnythingOfType("models.EnergyService")).Return(nil)
	mockProjector.On("Commit", mock.Anything).Return(nil)

	// Create a Consumer with the mock projector
	cons := &consumer.Consumer{
		Db: mockProjector,
	}

	// Call HandleBatch with the test message
	ctx := context.Background()
	err = cons.HandleBatch(ctx, []substrate.Message{testStruct{data: msgData}})

	// Assert that there were no errors and the expected calls were made
	assert.NoError(t, err)
	mockProjector.AssertExpectations(t)
}

type testStruct struct {
	data []byte
}

func (t testStruct) Data() []byte {
	return t.data
}