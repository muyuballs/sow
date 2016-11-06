package main

import (
	"io"
)

const (
	socksVer5 = 0x05

	socksCmdConnect            = 0x01
	cmdNotSupport              = 0x07
	NO_AUTHENTICATION_REQUIRED = 0x00
	NO_ACCEPTABLE_METHODS      = 0xFF
	ATYP_IP4                   = 0x01
	ATYP_DOMAIN                = 0x03
	ATYP_IP6                   = 0x04
)

const (
	socksVer4        = 0x04
	REQUEST_REJECTED = 0x5B
)

/*
			  X'00' NO AUTHENTICATION REQUIRED
	          X'01' GSSAPI
	          X'02' USERNAME/PASSWORD
	          X'03' to X'7F' IANA ASSIGNED
	          X'80' to X'FE' RESERVED FOR PRIVATE METHODS
	          X'FF' NO ACCEPTABLE METHODS
*/

func ReadByte(r io.Reader) (rel byte, err error) {
	buf := make([]byte, 1)
	_, err = io.ReadFull(r, buf)
	if err == nil {
		rel = buf[0]
	}
	return
}
