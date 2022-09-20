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

type ChunkMarkerHeader struct {
	Icon FourOctets
	Name FourOctets
}

type ChunkHeader struct {
	ChunkMarkerHeader
	OctetCount uint32
}

func ReadHeader(reader io.Reader) error {
	headerSpace := make([]byte, len(FileHeader()))

	reader.Read(headerSpace)

	if bytes.Compare(headerSpace, FileHeader()) != 0 {
		return fmt.Errorf("file header was unexpected")
	}

	return nil
}

func ReadChunkMarkerHeader(reader io.Reader) (ChunkMarkerHeader, error) {
	var icon FourOctets
	if err := binary.Read(reader, binary.BigEndian, &icon); err != nil {
		return ChunkMarkerHeader{}, err
	}
	var name FourOctets
	if err := binary.Read(reader, binary.BigEndian, &name); err != nil {
		return ChunkMarkerHeader{}, err
	}

	return ChunkMarkerHeader{
		Icon: icon,
		Name: name,
	}, nil
}

func ReadChunkHeader(reader io.Reader) (ChunkHeader, error) {
	chunkMarker, chunkMarkerErr := ReadChunkMarkerHeader(reader)
	if chunkMarkerErr != nil {
		return ChunkHeader{}, chunkMarkerErr
	}

	var chunkOctetCount uint32
	if err := binary.Read(reader, binary.BigEndian, &chunkOctetCount); err != nil {
		return ChunkHeader{}, err
	}

	return ChunkHeader{ChunkMarkerHeader: chunkMarker, OctetCount: chunkOctetCount}, nil
}

func ReadChunk(reader io.Reader) (ChunkHeader, []byte, error) {
	header, headerErr := ReadChunkHeader(reader)
	if headerErr != nil {
		return ChunkHeader{}, nil, headerErr
	}
	octetTarget := make([]byte, header.OctetCount)

	octetCountRead, readErr := reader.Read(octetTarget)
	if readErr != nil {
		return ChunkHeader{}, nil, readErr
	}

	if octetCountRead != int(header.OctetCount) {
		return ChunkHeader{}, nil, fmt.Errorf("wrong octet count read")
	}

	return header, octetTarget, nil
}

func ReadInternalChunkMarker(reader io.Reader) (ChunkHeader, error) {
	return ReadChunkHeader(reader)
}
