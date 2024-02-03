package tests

import (
	"os"
	"testing"
	"time"

	"github.com/ivinayakg/shorte.live/api/helpers"
)

func TestHelpers(t *testing.T) {
	t.Run("TestValidShortString", func(t *testing.T) {
		str1 := "helloworld"
		str2 := "#helloworld/"

		result1 := helpers.NotValidShortString(&str1)
		result2 := helpers.NotValidShortString(&str2)

		if result1 == true {
			t.Errorf("Expected %v, got %v", true, result1)
		}
		if result2 != true {
			t.Errorf("Expected %v, got %v", true, result2)
		}
		t.Log("TestValidShortString passed")
	})
	t.Run("TestTimeRemaining", func(t *testing.T) {
		tests := []struct {
			name     string
			duration time.Duration
			want     string
		}{
			{
				name:     "Test with positive duration",
				duration: time.Hour*25 + time.Minute*30 + time.Second*45,
				want:     "Time remaining = 01d 01h 30m 45s",
			},
			{
				name:     "Test with zero duration",
				duration: 0,
				want:     "Time has expired",
			},
			{
				name:     "Test with negative duration",
				duration: -time.Hour,
				want:     "Time has expired",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := helpers.TimeRemaining(tt.duration); got != tt.want {
					t.Errorf("TimeRemaining() = %v, want %v", got, tt.want)
				}
			})
		}
		t.Log("TestTimeRemaining passed")
	})
	t.Run("TestRemoverDomainError", func(t *testing.T) {
		// Set the DOMAIN environment variable for the test
		os.Setenv("DOMAIN", "example.com")

		tests := []struct {
			name string
			url  string
			want bool
		}{
			{
				name: "Test with domain",
				url:  "http://example.com/path",
				want: false,
			},
			{
				name: "Test with www domain",
				url:  "http://www.example.com/path",
				want: false,
			},
			{
				name: "Test with https domain",
				url:  "https://example.com/path",
				want: false,
			},
			{
				name: "Test with different domain",
				url:  "http://different.com/path",
				want: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := helpers.RemoverDomainError(tt.url); got != tt.want {
					t.Errorf("RemoverDomainError() = %v, want %v", got, tt.want)
				}
			})
		}

		t.Log("TestRemoverDomainError passed")
	})
	t.Run("TestEnforceHTTP", func(t *testing.T) {
		tests := []struct {
			name string
			url  string
			want string
		}{
			{
				name: "Test with http URL",
				url:  "http://example.com",
				want: "http://example.com",
			},
			{
				name: "Test with https URL",
				url:  "https://example.com",
				want: "https://example.com",
			},
			{
				name: "Test with non-http URL",
				url:  "example.com",
				want: "https://example.com",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := helpers.EnforceHTTP(tt.url); got != tt.want {
					t.Errorf("EnforceHTTP() = %v, want %v", got, tt.want)
				}
			})
		}
		t.Log("TestEnforceHTTP passed")
	})
}
