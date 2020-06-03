package password

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"hash"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

// Hash functions.
const (
	HashSHA256 = uint(1)
)

var encoding = base64.StdEncoding

// Service is a high-level interface used to validate,
// verify and hash passwords.
type Service interface {
	//Validate(password string) error
	Hash(password string) string
	Verify(password, hash string) bool

	// SetIterationCount(cnt uint)
	// SetSaltSize(cnt uint)
	// SetKeySize(cnt uint)
}

func New() Service {
	return &service{
		iterationCount: uint(15000),
		saltSize:       uint(128 / 8),
		keySize:        uint(256 / 8),
		hash:           HashSHA256,
	}
}

type service struct {
	iterationCount uint
	saltSize       uint
	keySize        uint
	hash           uint
}

// SetIterationCount sets the number of iterations the hasher will use.
// func (s *service) SetIterationCount(cnt uint) {
// 	s.iterationCount = cnt
// }

// SetSaltSize sets the size of the password salt.
//
// This must be divisible by 8.
// func (s *service) SetSaltSize(cnt uint) {
// 	s.saltSize = cnt / 8
// }

// SetKeySize sets the size of the password sub key.
//
// This must be divisible by 8.
// func (s *service) SetKeySize(cnt uint) {
// 	s.keySize = cnt / 8
// }

func (s *service) Hash(password string) string {
	salt := make([]byte, s.saltSize)
	rand.Read(salt)

	alg := getHashFunc(s.hash)
	subKey := pbkdf2.Key([]byte(password), salt, int(s.iterationCount), int(s.keySize), alg)
	fmt.Printf("\t%s\n", encoding.EncodeToString(subKey))

	output := make([]byte, 13+len(salt)+len(subKey))
	output[0] = 0x01 // format marker

	writeHeader(output, 1, s.hash)
	writeHeader(output, 5, s.iterationCount)
	writeHeader(output, 9, s.saltSize)

	copy(output[13:], salt)
	copy(output[13+len(salt):], subKey)

	return encoding.EncodeToString(output)
}

func writeHeader(buffer []byte, offset int, value uint) {
	buffer[offset+0] = byte(value >> 24)
	buffer[offset+1] = byte(value >> 16)
	buffer[offset+2] = byte(value >> 8)
	buffer[offset+3] = byte(value >> 0)
}

func (s *service) Verify(password, hash string) (success bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("failed to verify password: %v", r)
			success = false
		}
	}()

	hashedData, err := encoding.DecodeString(hash)
	if err != nil {
		log.Printf("failed to decode hashed password: %v", err)
		return false
	}

	// Read header info
	alg := getHashFunc(readHeader(hashedData, 1))
	iterCnt := int(readHeader(hashedData, 5))
	saltLen := int(readHeader(hashedData, 9))

	// Read the salt: it must be >= 128 bites
	if saltLen < int(s.saltSize) {
		return false
	}

	salt := make([]byte, saltLen)
	copy(salt[:], hashedData[13:13+saltLen])

	// Read the subkey: must be >= than the service's key size.
	subKeyLen := len(hashedData) - 13 - saltLen
	if subKeyLen < int(s.keySize) {
		return false
	}

	expectedSubKey := make([]byte, subKeyLen)
	copy(expectedSubKey[:], hashedData[13+saltLen:13+saltLen+subKeyLen])

	// Hash the incoming password.
	actualSubKey := pbkdf2.Key([]byte(password), salt, iterCnt, subKeyLen, alg)

	return subtle.ConstantTimeCompare(actualSubKey, expectedSubKey) == 1
}

func readHeader(buffer []byte, offset int) uint {
	return uint(buffer[offset+0])<<24 |
		uint(buffer[offset+1])<<16 |
		uint(buffer[offset+2])<<8 |
		uint(buffer[offset+3])
}

func getHashFunc(v uint) func() hash.Hash {
	switch v {
	case HashSHA256:
		return sha256.New
	default:
		panic(fmt.Errorf("unrecognized hash func: %d", v))
	}
}
