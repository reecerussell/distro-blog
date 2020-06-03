package password

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
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

var (
	// DefaultOptions are the service's default password validation options.
	DefaultOptions = &Options{
		RequiredLength:         6,
		RequireUppercase:       true,
		RequireLowercase:       true,
		RequireNonAlphanumeric: false,
		RequireDigit:           true,
		RequiredUniqueChars:    0,
	}

	// DefaultKeySize is the default key size.
	DefaultKeySize = uint(256 / 8)

	// DefaultSaltSize is the default salt size.
	DefaultSaltSize = uint(128 / 8)

	// DefaultHash is the default hash function.
	DefaultHash = HashSHA256

	// DefaultIterationCount is the default number of iterations
	// used by the hasher.
	DefaultIterationCount = uint(15000)
)

// Service is a high-level interface used to validate,
// verify and hash passwords.
type Service interface {
	Hash(password string) string
	Verify(password, hash string) bool
	Validate(password string) error

	// SetIterationCount(cnt uint)
	// SetSaltSize(cnt uint)
	// SetKeySize(cnt uint)
	SetValidationOptions(opts *Options)
}

// Options contains a set of critrea that a
// password has to hit, to be considered valid.
type Options struct {
	// RequiredLength is the minimum length a password has to be.
	RequiredLength int

	// RequireUppercase is a flag which demands at least one
	// character to be uppercase.
	RequireUppercase bool

	// RequireLowercase is a flag which demands at least one
	// character to be lowercase.
	RequireLowercase bool

	// RequireNonAlphanumeric is a flag which demands at least one
	// character to be non alphanumeric.
	RequireNonAlphanumeric bool

	// RequireDigit is a flag which demands at least one
	// character to be a digit.
	RequireDigit bool

	// RequiredUniqueChars is an integer value, that determines
	// how many unique characters are required in a password.
	RequiredUniqueChars int
}

func New() Service {
	return &service{
		iterationCount:    DefaultIterationCount,
		saltSize:          DefaultSaltSize,
		keySize:           DefaultKeySize,
		hash:              DefaultHash,
		validationOptions: DefaultOptions,
	}
}

type service struct {
	iterationCount    uint
	saltSize          uint
	keySize           uint
	hash              uint
	validationOptions *Options
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

// SetValidationOptions sets the service's validation options.
func (s *service) SetValidationOptions(opts *Options) {
	s.validationOptions = opts
}

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

func (s *service) Validate(password string) error {
	l := len(password)
	if l < 1 {
		return errors.New("password is required")
	}

	if l < s.validationOptions.RequiredLength {
		return fmt.Errorf("password must be atleast %d characters long", s.validationOptions.RequiredLength)
	}

	var (
		hasNonAlphanumeric bool
		hasDigit           bool
		hasLower           bool
		hasUpper           bool
		uniqueChars        []byte
	)

	for i := 0; i < l; i++ {
		c := byte(password[i])

		if !hasNonAlphanumeric && !isLetterOrDigit(c) {
			hasNonAlphanumeric = true
		}

		if !hasDigit && isDigit(c) {
			hasDigit = true
		}

		if !hasLower && isLower(c) {
			hasLower = true
		}

		if !hasUpper && isUpper(c) {
			hasUpper = true
		}

		d := true

		for _, dc := range uniqueChars {
			if dc == c {
				d = false
			}
		}

		if d {
			uniqueChars = append(uniqueChars, c)
		}
	}

	if s.validationOptions.RequireNonAlphanumeric && !hasNonAlphanumeric {
		return errors.New("password requires an non-alphanumeric character")
	}

	if s.validationOptions.RequireDigit && !hasDigit {
		return errors.New("password requires a digit")
	}

	if s.validationOptions.RequireLowercase && !hasLower {
		return errors.New("password requires a lowercase letter")
	}

	if s.validationOptions.RequireUppercase && !hasUpper {
		return errors.New("password requires an uppercase letter")
	}

	if s.validationOptions.RequiredUniqueChars >= 1 && len(uniqueChars) < s.validationOptions.RequiredUniqueChars {
		return fmt.Errorf("password requires atleast %d unique characters", s.validationOptions.RequiredUniqueChars)
	}

	return nil
}

// isDigit returns a flag indicating whether the supplied character
// is a digit - true if the character is a digit, otherwise false.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLower returns a flag indicating whether the supplied character is
// a lower case ASCII letter - true if the character is a lower case
// ASCII letter, otherwise false.
func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// isUpper returns a flag indicating whether the supplied character is
// an upper case ASCII letter - true if the character is an upper case ASCII
// letter, otherwise false.
func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

// isLetterOrDigit returns a flag indicating whether the supplied character is
// an ASCII letter or digit - true if the character is an ASCII letter or digit,
// otherwise false.
func isLetterOrDigit(c byte) bool {
	return isUpper(c) || isLower(c) || isDigit(c)
}
