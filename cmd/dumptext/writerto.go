package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

// slightly different than https://pkg.go.dev/io#WriterTo
type writerTo interface {
	WriteTo(io.Writer, []byte) error
}

// Writes raw bytes in native endianness.
type nativeEndianWriterTo struct{}

func (nativeEndianWriterTo) WriteTo(w io.Writer, b []byte) error {
	return binary.Write(w, binary.NativeEndian, b)
}

// Writes hex escaped bytes (Example: `\xde\xad\xbe\xef`).
type escapedHexBytesWriterTo struct{}

func (wt escapedHexBytesWriterTo) WriteTo(w io.Writer, b []byte) error {
	for _, v := range b {
		fmt.Fprint(w, "\\x"+wt.encode(v))
	}
	return nil
}

const hextable = "0123456789abcdef"

func (escapedHexBytesWriterTo) encode(v byte) string {
	return string(hextable[v>>4]) + string(hextable[v&0x0f])
}

// Writes a formatted hexdump.
type hexdumpWriterTo struct{}

func (wt hexdumpWriterTo) WriteTo(w io.Writer, b []byte) error {
	dumper := hex.Dumper(w)
	defer dumper.Close()
	_, err := dumper.Write(b)
	return err
}
