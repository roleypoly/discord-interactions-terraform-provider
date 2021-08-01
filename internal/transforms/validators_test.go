package transforms_test

import (
	"testing"

	"github.com/roleypoly/terraform-provider-discord-interactions/internal/transforms"
)

func TestNameValidator(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{
			name:     "my-cool-command",
			expected: true,
		},
		{
			name:     "MyUncoolCommand",
			expected: false,
		},
		{
			name:     "thisislongerthan32characterssoitshouldfail",
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			warns, errs := transforms.ValidateName(tC.name, "name")
			result := len(errs) == 0

			if result != tC.expected {
				t.Errorf("did not match expectation, got: %v, %v", warns, errs)
			}
		})
	}
}

func TestDescriptionValidator(t *testing.T) {
	testCases := []struct {
		desc     string
		expected bool
	}{
		{
			desc:     "",
			expected: false,
		},
		{
			desc:     "hello world!",
			expected: true,
		},
		{
			desc:     "hello world!hello world!hello world!hello world!hello world!hello world!hello world!hello world!hello world!hello world!",
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run("description: "+tC.desc, func(t *testing.T) {
			warns, errs := transforms.ValidateDescription(tC.desc, "description")
			result := len(errs) == 0

			if result != tC.expected {
				t.Errorf("did not match expectation, got: %v, %v", warns, errs)
			}
		})
	}
}

func TestSnowflake(t *testing.T) {
	testCases := []struct {
		desc      string
		snowflake string
		expected  bool
	}{
		{
			desc:      "good",
			snowflake: "386659935687147521",
			expected:  true,
		},
		{
			desc:      "bad",
			snowflake: "aaaaaaaaaa",
			expected:  false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			warns, errs := transforms.ValidateSnowflake(tC.snowflake, "id")
			result := len(errs) == 0

			if result != tC.expected {
				t.Errorf("did not match expectation, got: %v, %v", warns, errs)
			}
		})
	}
}
