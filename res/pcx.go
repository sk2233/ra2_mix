package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// PCXHeader represents the PCX file header structure
type PCXHeader struct {
	Manufacturer uint8        // Should be 10
	Version      uint8        // Should be 5
	Encoding     uint8        // Should be 1 (RLE)
	BitsPerPixel uint8        // Should be 8
	XMin         uint16       // Left margin
	YMin         uint16       // Top margin
	XMax         uint16       // Right margin
	YMax         uint16       // Bottom margin
	XDPi         uint16       // Horizontal resolution
	YDPi         uint16       // Vertical resolution
	ColorMap     [16][3]uint8 // Color palette (16 colors * 3 bytes)
	Reserved     uint8        // Should be 0
	Planes       uint8        // Number of color planes (1 or 3)
	BytesPerLine uint16       // Bytes per scan line
	PalType      uint16       // Palette type
	HScreenSize  uint16       // Horizontal screen size
	VScreenSize  uint16       // Vertical screen size
	Reserved2    [54]uint8    // Padding to make header 128 bytes
}

// PCXDecoder handles PCX file decoding
type PCXDecoder struct {
	header PCXHeader
	data   []byte
}

// NewPCXDecoder creates a new PCX decoder from file data
func NewPCXDecoder(data []byte) (*PCXDecoder, error) {
	if len(data) < 128 {
		return nil, errors.New("file too small to contain PCX header")
	}

	decoder := &PCXDecoder{data: data}

	// Parse header using binary.Read
	reader := bytes.NewReader(data[:128])
	if err := binary.Read(reader, binary.LittleEndian, &decoder.header); err != nil {
		return nil, fmt.Errorf("failed to read PCX header: %v", err)
	}

	return decoder, nil
}

// IsValid validates the PCX file format
func (d *PCXDecoder) IsValid() bool {
	h := d.header
	return h.Manufacturer == 10 &&
		h.Version == 5 &&
		h.Encoding == 1 &&
		h.BitsPerPixel == 8 &&
		(h.Planes == 1 || h.Planes == 3)
}

// GetDimensions returns the image dimensions
func (d *PCXDecoder) GetDimensions() (width, height int) {
	width = int(d.header.XMax - d.header.XMin + 1)
	height = int(d.header.YMax - d.header.YMin + 1)
	return
}

// GetImageData returns the compressed image data (after header)
func (d *PCXDecoder) GetImageData() []byte {
	return d.data[128:]
}

// GetPalette returns the 256-color palette from the end of the file
func (d *PCXDecoder) GetPalette() []byte {
	if len(d.data) < 768 {
		return nil
	}
	return d.data[len(d.data)-768:]
}

// decodeRLE performs RLE decompression on PCX data
func (d *PCXDecoder) decodeRLE(src []byte, expectedSize int) ([]byte, error) {
	var result []byte
	i := 0

	for i < len(src) && len(result) < expectedSize {
		b := src[i]
		i++

		if b < 0xC0 {
			// Single byte value
			result = append(result, b)
		} else {
			// RLE encoded: count in lower 6 bits, value in next byte
			if i >= len(src) {
				return nil, errors.New("unexpected end of data in RLE stream")
			}
			count := int(b & 0x3F)
			value := src[i]
			i++

			for j := 0; j < count && len(result) < expectedSize; j++ {
				result = append(result, value)
			}
		}
	}

	if len(result) != expectedSize {
		return nil, fmt.Errorf("decoded size mismatch: expected %d, got %d", expectedSize, len(result))
	}

	return result, nil
}

// Decode decodes the PCX image data into raw pixel data
func (d *PCXDecoder) Decode() ([]byte, error) {
	if !d.IsValid() {
		return nil, errors.New("invalid PCX file format")
	}

	width, height := d.GetDimensions()
	bytesPerLine := int(d.header.BytesPerLine)
	planes := int(d.header.Planes)

	// Calculate expected size of decompressed data
	expectedSize := bytesPerLine * height * planes

	// Decode RLE compressed data
	imageData := d.GetImageData()
	// Remove palette from image data if present
	if planes == 1 && len(imageData) > expectedSize {
		imageData = imageData[:len(imageData)-768]
	}

	decompressed, err := d.decodeRLE(imageData, expectedSize)
	if err != nil {
		return nil, fmt.Errorf("RLE decode failed: %v", err)
	}

	if planes == 1 {
		// Single plane (indexed color) - extract actual image data
		result := make([]byte, width*height)
		srcIdx := 0
		dstIdx := 0

		for y := 0; y < height; y++ {
			// Copy only the actual width, skip padding
			copy(result[dstIdx:dstIdx+width], decompressed[srcIdx:srcIdx+width])
			srcIdx += bytesPerLine
			dstIdx += width
		}

		return result, nil
	} else {
		// Three planes (RGB) - interleave the planes
		result := make([]byte, width*height*3)

		for y := 0; y < height; y++ {
			rOffset := y * bytesPerLine
			gOffset := rOffset + bytesPerLine*height
			bOffset := gOffset + bytesPerLine*height

			for x := 0; x < width; x++ {
				resultIdx := (y*width + x) * 3
				result[resultIdx] = decompressed[rOffset+x]   // R
				result[resultIdx+1] = decompressed[gOffset+x] // G
				result[resultIdx+2] = decompressed[bOffset+x] // B
			}
		}

		return result, nil
	}
}

// LoadPCXFile loads and decodes a PCX file
func LoadPCXFile(filename string) ([]byte, int, int, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, 0, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, 0, nil, fmt.Errorf("failed to read file: %v", err)
	}

	decoder, err := NewPCXDecoder(data)
	if err != nil {
		return nil, 0, 0, nil, err
	}

	imageData, err := decoder.Decode()
	if err != nil {
		return nil, 0, 0, nil, err
	}

	width, height := decoder.GetDimensions()
	palette := decoder.GetPalette()

	return imageData, width, height, palette, nil
}
