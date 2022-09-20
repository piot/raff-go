/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

// Package raff contains chunk and file headers
package raff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// FileHeader is octets to be inserted at the start of a file.
func FileHeader() []byte {
	return []byte{0xF0, 0x9F, 0xA6, 0x8A, 'R', 'A', 'F', 'F', 0x0a}
}

type FourOctets uint32

func IconToString(f FourOctets) string {
	c1 := byte(f >> 24)
	c2 := byte((f >> 16) & 0xff)
	c3 := byte((f >> 8) & 0xff)
	c4 := byte(f & 0xff)

	octets := []byte{c1, c2, c3, c4}

	return string(octets)
}

func NameToString(f FourOctets) string {
	c1 := byte(f >> 24)
	c2 := byte((f >> 16) & 0xff)
	c3 := byte((f >> 8) & 0xff)
	c4 := byte(f & 0xff)

	return string(c1) + string(c2) + string(c3) + string(c4)
}

func write(writer io.Writer, octets []byte) error {
	writtenOctetCount, err := writer.Write(octets)
	if err != nil {
		return err
	}

	if writtenOctetCount != len(octets) {
		return fmt.Errorf("raff: couldn't write everything")
	}

	return err
}

// WriteHeader writes a file header.
func WriteHeader(writer io.Writer) error {
	return write(writer, FileHeader())
}

// WriteChunk writes a octet slice to file with an extended header.
func WriteChunk(writer io.Writer, icon FourOctets, name FourOctets, octets []byte) error {
	var temp bytes.Buffer

	binary.Write(&temp, binary.BigEndian, icon)
	binary.Write(&temp, binary.BigEndian, name)

	chunkCount := uint32(len(octets))

	if err := binary.Write(&temp, binary.BigEndian, chunkCount); err != nil {
		return err
	}

	if err := write(writer, temp.Bytes()); err != nil {
		return err
	}

	if err := write(writer, octets); err != nil {
		return err
	}

	return nil
}

// WriteInternalChunkMarker writes a octet slice to file with an extended header.
func WriteInternalChunkMarker(writer io.Writer, icon FourOctets) error {
	if err := binary.Write(writer, binary.BigEndian, icon); err != nil {
		return err
	}

	return nil
}
