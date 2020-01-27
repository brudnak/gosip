package api

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func TestWeb(t *testing.T) {
	t.Parallel()
	checkClient(t)

	sp := NewSP(spClient)
	web := sp.Web()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web"

	t.Run("Constructor", func(t *testing.T) {
		web := NewWeb(spClient, endpoint, nil)
		if _, err := web.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if web.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				web.ToURL(),
			)
		}
	})

	t.Run("ToURLWithModifiers", func(t *testing.T) {
		apiURL, _ := url.Parse(web.ToURL())
		query := url.Values{
			"$select": []string{"Title,Webs/Title"},
			"$expand": []string{"Webs"},
		}
		apiURL.RawQuery = query.Encode()
		expectedURL := apiURL.String()

		resURL := web.Select("Title,Webs/Title").Expand("Webs").ToURL()
		if resURL != expectedURL {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				expectedURL,
				resURL,
			)
		}
	})

	t.Run("FromURL", func(t *testing.T) {
		w := web.FromURL("site_url")
		if w.endpoint != "site_url" {
			t.Error("can't get site from url")
		}
	})

	t.Run("GetTitle", func(t *testing.T) {
		data, err := web.Select("Title").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().Title == "" {
			t.Error("can't get web title property")
		}

		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("NoTitle", func(t *testing.T) {
		data, err := web.Select("Id").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().Title != "" {
			t.Error("can't get web title property")
		}
	})

	t.Run("CurrentUser", func(t *testing.T) {
		if spClient.AuthCnfg.GetStrategy() == "addin" {
			t.Skip("is not applicable for Addin Only auth strategy")
		}

		data, err := web.CurrentUser().Select("LoginName").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().LoginName == "" {
			t.Error("can't get current user")
		}
	})

	t.Run("EnsureFolder", func(t *testing.T) {
		guid := uuid.New().String()
		if _, err := web.EnsureFolder("Shared Documents/" + guid + "/doc1/doc2/doc3/doc4"); err != nil {
			t.Error(err)
		}
		if err := web.GetFolder("Shared Documents/" + guid).Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("EnsureUser", func(t *testing.T) {
		user, err := sp.Web().CurrentUser().Get()
		if err != nil {
			t.Error(err)
		}
		if _, err := sp.Web().EnsureUser(user.Data().Email); err != nil {
			t.Error(err)
		}
	})

	t.Run("Roles", func(t *testing.T) {
		if _, err := sp.Web().Roles().HasUniqueAssignments(); err != nil {
			t.Error(err)
		}
	})

	t.Run("AvailableContentTypes", func(t *testing.T) {
		resp, err := sp.Web().AvailableContentTypes().Get()
		if err != nil {
			t.Error(err)
		}
		if len(resp.Data()) == 0 {
			t.Error("can't get available content types")
		}
	})

}
