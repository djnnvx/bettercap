package uds

import (
	"context"
	"testing"

	canmod "github.com/bettercap/bettercap/v2/modules/can"
	"go.einride.tech/can"
	"github.com/stretchr/testify/assert"
)

type mockCAN struct {
	lastFrame can.Frame
}

func (m *mockCAN) SendFrame(ctx context.Context, frame can.Frame) error {
	m.lastFrame = frame
	return nil
}

func TestUDSRequestBytes(t *testing.T) {
	req := UDSRequest{ServiceID: 0x27, Params: []byte{0x01, 0x02, 0x03}}
	assert.Equal(t, []byte{0x27, 0x01, 0x02, 0x03}, req.Bytes())
}

func TestSendUDSMessage(t *testing.T) {
	mock := &mockCAN{}
	canModule := &canmod.CANModule{}
	canModule.SendFrame = mock.SendFrame

	uds := &UDS{CAN: canModule}
	ctx := context.Background()
	req := UDSRequest{ServiceID: UDS_SecurityAccess, Params: []byte{0x01, 0x11, 0x22}}
	arbID := uint32(0x7E0)

	err := uds.SendUDSMessage(ctx, arbID, req)
	assert.NoError(t, err)
	assert.Equal(t, arbID, mock.lastFrame.ID)
	assert.Equal(t, []byte{0x27, 0x01, 0x11, 0x22}, mock.lastFrame.Data[:mock.lastFrame.Length])
}
