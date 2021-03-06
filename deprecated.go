package sunrpc

import (
	"bytes"
	"io"
	"net"

	"github.com/rasky/go-xdr/xdr2"
)

// NOTE: This function is DEPRECATED and should not be exported anymore.
// Please use Client instead to implement a proper sunrpc client
func WriteCall(w io.Writer, program uint32, version uint32, proc uint32, args interface{}) error {
	var buf bytes.Buffer

	if _, err := xdr.Marshal(&buf, NewProcedureCall(program, version, proc)); err != nil {
		return err
	}

	// Write procedure arguments to the buffer
	if _, err := xdr.Marshal(&buf, args); err != nil {
		return err
	}

	// On TCP transport, we need to write a record marker
	// FIXME: this sniffing is really ugly; it'd be better to have a proper
	// client class that knows whether it's TCP or UDP.
	if _, ok := w.(*net.UDPConn); !ok {
		if err := WriteRecordMarker(w, uint32(buf.Len()), true); err != nil {
			return err
		}
	}

	// Send the payload
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
