package api

import (
	"testing"
)

func TestSP(t *testing.T) {
	checkClient(t)

	t.Run("ToURL", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.ToURL() != spClient.AuthCnfg.GetSiteURL() {
			t.Errorf(
				"incorrect site URL, expected \"%s\", received \"%s\"",
				spClient.AuthCnfg.GetSiteURL(),
				sp.ToURL(),
			)
		}
	})

	t.Run("Web", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Web() == nil {
			t.Errorf("failed to get Web object")
		}
	})

	t.Run("Site", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Site() == nil {
			t.Errorf("failed to get Site object")
		}
	})

	t.Run("Utility", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Utility() == nil {
			t.Errorf("failed to get Site object")
		}
	})

	t.Run("Metadata", func(t *testing.T) {
		sp := NewSP(spClient)
		if _, err := sp.Metadata(); err != nil {
			t.Error(err)
		}
	})

}
