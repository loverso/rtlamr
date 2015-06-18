// RTLAMR - An rtl-sdr receiver for smart meters operating in the 900MHz ISM band.
// Copyright (C) 2015 Douglas Hall
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package scmplus

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/bemasher/rtlamr/crc"
	"github.com/bemasher/rtlamr/decode"
	"github.com/bemasher/rtlamr/parse"
)

func NewPacketConfig(symbolLength int) (cfg decode.PacketConfig) {
	cfg.DataRate = 32768
	cfg.CenterFreq = 912600155
	cfg.SymbolLength = symbolLength
	cfg.PreambleSymbols = 24
	cfg.PacketSymbols = 17 * 8
	cfg.Preamble = "010101010001011010100011"

	return
}

type Parser struct {
	decode.Decoder
	crc.CRC
}

func (p Parser) Dec() decode.Decoder {
	return p.Decoder
}

func (p Parser) Cfg() decode.PacketConfig {
	return p.Decoder.Cfg
}

func NewParser(symbolLength, decimation int, fastMag bool) (p Parser) {
	p.Decoder = decode.NewDecoder(NewPacketConfig(symbolLength), decimation, fastMag)
	p.CRC = crc.NewCRC("CCITT", 0xFFFF, 0x1021, 0x1D0F)
	return
}

func (p Parser) Parse(indices []int) (msgs []parse.Message) {
	seen := make(map[string]bool)

	for _, pkt := range p.Decoder.Slice(indices) {
		s := string(pkt)
		if seen[s] {
			continue
		}
		seen[s] = true

		data := parse.NewDataFromBytes(pkt)

		if data.Bytes[3] != 0x1E {
			continue
		}

		// // If the checksum fails, bail.
		// if residue := p.Checksum(data.Bytes[3:]); residue != p.Residue {
		// 	continue
		// }

		var scm SCMPlus
		scm.Preamble = binary.BigEndian.Uint32(data.Bytes[0:4]) >> 8
		scm.PacketTypeID = data.Bytes[3]
		scm.EndpointType = data.Bytes[4]
		scm.EndpointID = binary.BigEndian.Uint32(data.Bytes[5:9])
		scm.Consumption = binary.BigEndian.Uint32(data.Bytes[9:13])
		scm.Tamper = binary.BigEndian.Uint16(data.Bytes[13:15])
		scm.PacketCRC = binary.BigEndian.Uint16(data.Bytes[15:17])

		// // If the meter id is 0, bail.
		// if scm.ERTSerialNumber == 0 {
		// 	continue
		// }

		msgs = append(msgs, scm)
	}

	return
}

// Standard Consumption Message
type SCMPlus struct {
	Preamble     uint32 // Training and Frame sync.
	PacketTypeID uint8
	EndpointType uint8
	EndpointID   uint32
	Consumption  uint32
	Tamper       uint16
	PacketCRC    uint16
}

func (scm SCMPlus) MsgType() string {
	return "SCMPlus"
}

func (scm SCMPlus) MeterID() uint32 {
	return scm.EndpointID
}

func (scm SCMPlus) MeterType() uint8 {
	return scm.EndpointType
}

func (scm SCMPlus) String() string {
	var fields []string

	fields = append(fields, fmt.Sprintf("Preamble: 0x%02X", scm.Preamble))
	fields = append(fields, fmt.Sprintf("PacketTypeID: 0x%02X", scm.PacketTypeID))
	fields = append(fields, fmt.Sprintf("EndpointType: %02d", scm.EndpointType))
	fields = append(fields, fmt.Sprintf("EndpointID: %10d", scm.EndpointID))
	fields = append(fields, fmt.Sprintf("Consumption: %10d", scm.Consumption))
	fields = append(fields, fmt.Sprintf("Tamper: 0x%04X", scm.Tamper))
	fields = append(fields, fmt.Sprintf("PacketCRC: 0x%04X", scm.PacketCRC))

	return "{" + strings.Join(fields, " ") + "}"
}

func (scm SCMPlus) Record() (r []string) {
	r = append(r, fmt.Sprintf("0x%02X\n", scm.Preamble))
	r = append(r, fmt.Sprintf("0x%02X\n", scm.PacketTypeID))
	r = append(r, fmt.Sprintf("0x%02X\n", scm.EndpointType))
	r = append(r, fmt.Sprintf("0x%02X\n", scm.EndpointID))
	r = append(r, fmt.Sprintf("0x%02X\n", scm.Consumption))
	r = append(r, fmt.Sprintf("0x%02X\n", scm.Tamper))
	r = append(r, fmt.Sprintf("0x%02X\n", scm.PacketCRC))

	return
}
