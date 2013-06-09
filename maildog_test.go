/**
 *  Unit tests for maildog
 */

package mailman

import (
	"testing"
)

func AssertListContains (t *testing.T, item string, list string) {
	r := ListContains(item, list)
	if !r {
		t.Errorf("Failed test case: %s should be in %s", item, list)
	}
}

func AssertNotListContains (t *testing.T, item string, list string) {
	r := ListContains(item, list)
	if r {
		t.Errorf("Failed test case: %s should not be in %s", item, list)
	}
}

func TestListContains (t *testing.T) {
	AssertListContains(t, "boot@domain.com", "boot toot@example.com")
	AssertListContains(t, "boot@domain.com", "boot@domain.com boot@example.com")
	AssertNotListContains(t, "boot@domain.com", "toot boot@example.com")
	AssertListContains(t, "boot@domain.com", "boot boot@example.com")
}
