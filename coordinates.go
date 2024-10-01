package main

import (
	"fmt"
)


func Encode3D(x, y, z uint32, method string = "morton") uint64 {  
	switch method {
		case "morton":
			return EncodeMorton3D(x, y, z)
		}
}

func Decode3D(morton uint64, method string = "morton") (x, y, z uint32) {
	switch method {
		case "morton":
			return DecodeMorton3D(morton)
		}
}


// EncodeMorton3D interleaves the bits of x, y, and z to produce a Morton code.
func EncodeMorton3D(x, y, z uint32) uint64 {
	return interleaveBits64(uint64(x), uint64(y), uint64(z))
}

// DecodeMorton3D de-interleaves the bits of a Morton code to produce x, y, and z.
func DecodeMorton3D(morton uint64) (x, y, z uint32) {
	x = uint32(deinterleaveBits64(morton))
	y = uint32(deinterleaveBits64(morton >> 1))
	z = uint32(deinterleaveBits64(morton >> 2))
	return
}

// interleaveBits64 interleaves the bits of three 64-bit numbers.
func interleaveBits64(x, y, z uint64) uint64 {
	x = spreadBits(x)
	y = spreadBits(y)
	z = spreadBits(z)
	return x | (y << 1) | (z << 2)
}

// spreadBits spreads the bits of a 64-bit number by inserting two zeros between each bit.
func spreadBits(n uint64) uint64 {
	n &= 0x1fffff // Limit to 21 bits to avoid overflow
	n = (n | n<<32) & 0x1f00000000ffff
	n = (n | n<<16) & 0x1f0000ff0000ff
	n = (n | n<<8) & 0x100f00f00f00f00f
	n = (n | n<<4) & 0x10c30c30c30c30c3
	n = (n | n<<2) & 0x1249249249249249
	return n
}

// deinterleaveBits64 compacts the bits of a Morton code to retrieve the original number.
func deinterleaveBits64(n uint64) uint64 {
	n &= 0x1249249249249249
	n = (n ^ (n >> 2)) & 0x10c30c30c30c30c3
	n = (n ^ (n >> 4)) & 0x100f00f00f00f00f
	n = (n ^ (n >> 8)) & 0x1f0000ff0000ff
	n = (n ^ (n >> 16)) & 0x1f00000000ffff
	n = (n ^ (n >> 32)) & 0x1fffff
	return n
}
