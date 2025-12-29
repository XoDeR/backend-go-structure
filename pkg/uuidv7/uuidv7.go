package uuidv7

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
)

func New() uuid.UUID {
	return NewWithTime(time.Now())
}

func NewWithTime(t time.Time) uuid.UUID {
	var u uuid.UUID

	unixMs := uint64(t.UnixMilli())

	// Timestamp: first 48 bits (6 bytes)
	binary.BigEndian.PutUint32(u[0:4], uint32(unixMs>>16))
	binary.BigEndian.PutUint16(u[4:6], uint16(unixMs))

	// Random bytes for the rest
	_, _ = rand.Read(u[6:16])

	// Set version (4 bits): 0111 = version 7
	u[6] = (u[6] & 0x0f) | 0x70

	// Set variant (2 bits): 10 = RFC 4122 variant
	u[8] = (u[8] & 0x3f) | 0x80

	return u
}

func ExtractTime(u uuid.UUID) time.Time {
	// Verify this is a UUID v7 (version bits should be 0111)
	if u[6]>>4 != 0x07 {
		return time.Time{}
	}

	// Extract 48-bit timestamp from first 6 bytes
	msHigh := uint64(binary.BigEndian.Uint32(u[0:4]))
	msLow := uint64(binary.BigEndian.Uint16(u[4:6]))
	unixMs := (msHigh << 16) | msLow

	return time.UnixMilli(int64(unixMs))
}

func IsV7(u uuid.UUID) bool {
	return u[6]>>4 == 0x07
}

func Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}

func MustParse(s string) UUID {
	return uuid.MustParse(s)
}

// UUID is a re-export of uuid.UUID for convenience
// Allows to use uuidv7.UUID throughout the codebase
type UUID = uuid.UUID

// Nil is the nil UUID
var Nil = uuid.Nil
