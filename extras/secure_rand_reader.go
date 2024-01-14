package extras

import crand "crypto/rand"

// SecureRandReader is a random bytes reader that use a cryptographically secure pseudorandom number generator.
type SecureRandReader struct {
}

// Read implements io.Reader interface for SecureRandReader
// It populates the given byte slice with random bytes.
func (s *SecureRandReader) Read(p []byte) (n int, err error) {
	return crand.Read(p)
}

// NewSecureReader creates a new SecureRandReader
func NewSecureReader() *SecureRandReader {
	return &SecureRandReader{}
}
