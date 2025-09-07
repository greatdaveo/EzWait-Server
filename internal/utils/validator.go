package utils

import (
	"regexp"
	"strings"
	"time"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}

	return strings.Join(messages, "; ")
}

// For Email Validation
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// For UK Phone Number Validation
func ValidateUKPhone(phone string) bool {
	// To replace spaces and special characters
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// UK phone patterns
	patterns := []string{
		`^07\d{9}$`,   // Mobile
		`^01\d{8,9}$`, // Landline
		`^02\d{8,9}$`, // Landline
		`^03\d{9}$`,   // 03 numbers
	}

	for _, pattern := range patterns {
		if regexp.MustCompile(pattern).MatchString(phone) {
			return true
		}
	}
	return false
}

// For UK Postcode validation
func ValidateUKPostcode(postcode string) bool {
	postcodeRegex := regexp.MustCompile(`^[A-Z]{1,2}[0-9][A-Z0-9]? ?[0-9][A-Z]{2}$`)
	return postcodeRegex.MatchString(strings.ToUpper(postcode))
}

// For Password strength validation
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false, "Password must contain at least one uppercase letter"
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false, "Password must contain at least one lowercase letter"
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false, "Password must contain at least one number"
	}

	return true, ""
}

// For Booking time validation
func ValidateBookingTime(startTime, endTime time.Time, bookingDay time.Time) (bool, string) {
	now := time.Now()

	// To check if booking is in the future
	if bookingDay.Before(now.Truncate(24 * time.Hour)) {
		return false, "Booking must be for a future date"
	}

	// To check if start time is before end time
	if startTime.After(endTime) || startTime.Equal(endTime) {
		return false, "Start time must be before end time"
	}

	// To check if booking is within business hours (8 AM to 8 PM)
	startHour := startTime.Hour()
	endHour := endTime.Hour()

	if startHour < 8 || endHour > 20 {
		return false, "Bookings must be between 8 AM and 8 PM"
	}

	return true, ""
}

// To sanitize input strings
func SanitizeString(input string) string {
	// To remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlRegex.ReplaceAllString(input, "")

	// To trim whitespace
	input = strings.TrimSpace(input)

	return input
}
