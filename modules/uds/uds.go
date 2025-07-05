package uds

import (
	"context"
	"fmt"

	canmod "github.com/bettercap/bettercap/v2/modules/can"
	"go.einride.tech/can"
)

// UDS represents the UDS protocol module layered on top of CAN.
type UDS struct {
	CAN *canmod.CANModule
}

// UDS service IDs
const (
	UDS_DiagnosticSessionControl   byte = 0x10
	UDS_ECUReset                   byte = 0x11
	UDS_SecurityAccess             byte = 0x27
	UDS_CommunicationControl       byte = 0x28
	UDS_ReadMemoryByAddress        byte = 0x23
	UDS_WriteMemoryByAddress       byte = 0x3D
	// ... other service IDs as needed
)

// UDSRequest represents a UDS request message.
type UDSRequest struct {
	ServiceID byte
	Params    []byte
}

// Bytes returns the UDS message as a CAN payload.
func (r UDSRequest) Bytes() []byte {
	return append([]byte{r.ServiceID}, r.Params...)
}

// SendUDSMessage sends a UDS message via CAN using the provided arbitration ID.
func (u *UDS) SendUDSMessage(ctx context.Context, arbID uint32, req UDSRequest) error {
	frame := can.Frame{
		ID:     arbID,
		Length: uint8(len(req.Bytes())),
	}
	copy(frame.Data[:], req.Bytes())
	if u.CAN == nil || u.CAN.SendFrame == nil {
		return fmt.Errorf("CAN module or SendFrame not available")
	}
	return u.CAN.SendFrame(ctx, frame)
}

// FuzzUDSMessage is a placeholder for future fuzzing support.
func (u *UDS) FuzzUDSMessage(req UDSRequest) UDSRequest {
	// TODO: Integrate with bettercap's mutators for fuzzing.
	return req
}

// New returns a new UDS instance.
func New(c *canmod.CANModule) *UDS {
	return &UDS{
		CAN: c,
	}
}
