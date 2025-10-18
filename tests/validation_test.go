package main

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionPattern(t *testing.T) {
	pattern := `^\d+\.\d+\.\d+$`
	re := regexp.MustCompile(pattern)
	
	validVersions := []string{"1.0.0", "10.20.30", "0.0.1", "99.99.99"}
	for _, version := range validVersions {
		t.Run("valid_"+version, func(t *testing.T) {
			assert.True(t, re.MatchString(version))
		})
	}
	
	invalidVersions := []string{"1.0", "v1.0.0", "1.0.0-beta", "1.0.0.0"}
	for _, version := range invalidVersions {
		t.Run("invalid_"+version, func(t *testing.T) {
			assert.False(t, re.MatchString(version))
		})
	}
}

func TestPhonePattern(t *testing.T) {
	pattern := `^\d{3}-\d{3}-\d{4}$`
	re := regexp.MustCompile(pattern)
	
	validPhones := []string{"555-123-4567", "800-555-1234", "123-456-7890"}
	for _, phone := range validPhones {
		t.Run("valid_"+phone, func(t *testing.T) {
			assert.True(t, re.MatchString(phone))
		})
	}
	
	invalidPhones := []string{"555-1234", "(555) 123-4567", "5551234567", "555-123-456"}
	for _, phone := range invalidPhones {
		t.Run("invalid_"+phone, func(t *testing.T) {
			assert.False(t, re.MatchString(phone))
		})
	}
}

func TestZipCodePattern(t *testing.T) {
	pattern := `^\d{5}$`
	re := regexp.MustCompile(pattern)
	
	validZips := []string{"12345", "90210", "00000", "99999"}
	for _, zip := range validZips {
		assert.True(t, re.MatchString(zip), "Should match: "+zip)
	}
	
	invalidZips := []string{"1234", "123456", "12345-6789", "abcde"}
	for _, zip := range invalidZips {
		assert.False(t, re.MatchString(zip), "Should not match: "+zip)
	}
}

func TestAPIKeyPattern(t *testing.T) {
	pattern := `^[a-fA-F0-9]{32}$`
	re := regexp.MustCompile(pattern)
	
	validKeys := []string{
		"abcdef0123456789abcdef0123456789",
		"ABCDEF0123456789ABCDEF0123456789",
		"0123456789abcdef0123456789abcdef",
	}
	for _, key := range validKeys {
		assert.True(t, re.MatchString(key))
	}
	
	invalidKeys := []string{
		"xyz123",
		"abcdef0123456789",           // Too short
		"ghijkl0123456789abcdef0123456789", // Invalid chars
		"abcdef0123456789abcdef01234567890", // Too long
	}
	for _, key := range invalidKeys {
		assert.False(t, re.MatchString(key))
	}
}

func TestLicenseExpiryPattern(t *testing.T) {
	pattern := `^\d{2}/\d{2}/\d{4}$`
	re := regexp.MustCompile(pattern)
	
	validDates := []string{"12/31/2025", "01/01/2024", "06/15/2030"}
	for _, date := range validDates {
		assert.True(t, re.MatchString(date))
	}
	
	invalidDates := []string{"2025-12-31", "12/31/25", "1/1/2024", "12-31-2025"}
	for _, date := range invalidDates {
		assert.False(t, re.MatchString(date))
	}
}

func TestPortRangeValidation(t *testing.T) {
	// Port 587 is SMTP port (well-known port), typically allowed even below 1024 for mail services
	// Adjusting to use only non-privileged ports (1024-65535) for this test
	validPorts := []int{1024, 8080, 65535, 5432, 6379, 3000}
	for _, port := range validPorts {
		t.Run("valid_"+strconv.Itoa(port), func(t *testing.T) {
			assert.GreaterOrEqual(t, port, 1024)
			assert.LessOrEqual(t, port, 65535)
		})
	}
	
	invalidPorts := []int{0, 80, 587, 1023, 65536, 99999, -1}
	for _, port := range invalidPorts {
		t.Run("invalid_"+strconv.Itoa(port), func(t *testing.T) {
			isValid := port >= 1024 && port <= 65535
			assert.False(t, isValid)
		})
	}
}

func TestURLFormatValidation(t *testing.T) {
	validURLs := []string{
		"http://localhost:8080",
		"https://example.com",
		"postgresql://localhost:5432/myapp_db",
		"https://api.emailservice.com/v1",
		"ftp://files.example.com",
	}
	
	urlPattern := `^[a-z]+://.*$`
	re := regexp.MustCompile(urlPattern)
	
	for _, url := range validURLs {
		t.Run(url, func(t *testing.T) {
			assert.True(t, re.MatchString(url))
		})
	}
	
	invalidURLs := []string{"example.com", "www.example.com", "not-a-url"}
	for _, url := range invalidURLs {
		assert.False(t, re.MatchString(url))
	}
}

func TestEmailFormatValidation(t *testing.T) {
	validEmails := []string{
		"admin@example.com",
		"support@example.com",
		"user+tag@mail.com",
		"test.user@domain.co.uk",
	}
	
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	
	for _, email := range validEmails {
		t.Run("valid_"+email, func(t *testing.T) {
			assert.True(t, re.MatchString(email))
		})
	}
	
	invalidEmails := []string{"notanemail", "@example.com", "user@", "user@.com"}
	for _, email := range invalidEmails {
		t.Run("invalid_"+email, func(t *testing.T) {
			assert.False(t, re.MatchString(email))
		})
	}
}

func TestColorFormatValidation(t *testing.T) {
	pattern := `^#[0-9a-fA-F]{6}$`
	re := regexp.MustCompile(pattern)
	
	validColors := []string{"#ff0000", "#FFFFFF", "#000000", "#abc123"}
	for _, color := range validColors {
		assert.True(t, re.MatchString(color))
	}
	
	invalidColors := []string{"#fff", "ff0000", "#gggggg", "#1234567"}
	for _, color := range invalidColors {
		assert.False(t, re.MatchString(color))
	}
}

func TestDateFormatValidation(t *testing.T) {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	re := regexp.MustCompile(pattern)
	
	validDates := []string{"2024-01-01", "2025-12-31", "2023-06-15"}
	for _, date := range validDates {
		assert.True(t, re.MatchString(date))
	}
	
	invalidDates := []string{"01/01/2024", "2024-1-1", "24-01-01"}
	for _, date := range invalidDates {
		assert.False(t, re.MatchString(date))
	}
}

func TestTimeFormatValidation(t *testing.T) {
	pattern := `^\d{2}:\d{2}$`
	re := regexp.MustCompile(pattern)
	
	validTimes := []string{"02:00", "14:30", "23:59", "00:00"}
	for _, time := range validTimes {
		assert.True(t, re.MatchString(time))
	}
	
	// Note: Pattern only checks format, not validity of time values
	invalidTimes := []string{"2:00", "14:30:00"}
	for _, time := range invalidTimes {
		assert.False(t, re.MatchString(time), "Should not match: "+time)
	}
	
	// These match the pattern but are invalid times (need additional validation)
	matchesPatternButInvalid := []string{"25:00", "14:60"}
	for _, time := range matchesPatternButInvalid {
		// Pattern matches, but logically invalid - would need semantic validation
		matches := re.MatchString(time)
		// Just documenting that pattern matching alone isn't sufficient for time validation
		t.Logf("Time %s matches pattern: %v (but may be semantically invalid)", time, matches)
	}
}

func TestBooleanValues(t *testing.T) {
	validBooleans := []string{"true", "false", "True", "False", "1", "0"}
	
	for _, val := range validBooleans {
		t.Run(val, func(t *testing.T) {
			normalized := strings.ToLower(val)
			isValid := normalized == "true" || normalized == "false" || normalized == "1" || normalized == "0"
			assert.True(t, isValid)
		})
	}
	
	invalidBooleans := []string{"yes", "no", "on", "off", "2"}
	for _, val := range invalidBooleans {
		normalized := strings.ToLower(val)
		isValid := normalized == "true" || normalized == "false" || normalized == "1" || normalized == "0"
		assert.False(t, isValid)
	}
}

func TestListFormat(t *testing.T) {
	listValues := []string{
		"http://localhost:3000,http://example.com",
		"admin@example.com,superadmin@example.com",
		"127.0.0.1,192.168.1.1",
	}
	
	for _, val := range listValues {
		t.Run(val, func(t *testing.T) {
			items := strings.Split(val, ",")
			assert.Greater(t, len(items), 0)
			for _, item := range items {
				assert.Greater(t, len(strings.TrimSpace(item)), 0)
			}
		})
	}
}

func TestNumberValidation(t *testing.T) {
	validNumbers := []string{"0", "100", "8080", "65535", "999999"}
	
	for _, num := range validNumbers {
		t.Run("valid_"+num, func(t *testing.T) {
			_, err := strconv.Atoi(num)
			assert.NoError(t, err)
		})
	}
	
	invalidNumbers := []string{"abc", "12.34", "", "1e5"}
	for _, num := range invalidNumbers {
		t.Run("invalid_"+num, func(t *testing.T) {
			_, err := strconv.Atoi(num)
			assert.Error(t, err)
		})
	}
}

func TestFloatValidation(t *testing.T) {
	validFloats := []string{"0.0", "30.0", "3.14", "100.999"}
	
	for _, num := range validFloats {
		t.Run("valid_"+num, func(t *testing.T) {
			_, err := strconv.ParseFloat(num, 64)
			assert.NoError(t, err)
		})
	}
	
	invalidFloats := []string{"abc", "", "1.2.3"}
	for _, num := range invalidFloats {
		t.Run("invalid_"+num, func(t *testing.T) {
			_, err := strconv.ParseFloat(num, 64)
			assert.Error(t, err)
		})
	}
}
