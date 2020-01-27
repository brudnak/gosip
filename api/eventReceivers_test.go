package api

import (
	"testing"
)

func TestEventReceivers(t *testing.T) {
	t.Parallel()
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("Get/Site", func(t *testing.T) {
		receivers, err := sp.Site().EventReceivers().Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if receivers[0].ReceiverID == "" {
			t.Error("can't get event receivers")
		}
	})

	t.Run("Get/Web", func(t *testing.T) {
		receivers, err := sp.Web().EventReceivers().Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if receivers[0].ReceiverID == "" {
			t.Error("can't get event receivers")
		}
	})

}
