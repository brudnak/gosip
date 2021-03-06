package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Changes -conf -coll -mods Top

// Changes represent SharePoint Changes API queryable collection struct
// Always use NewChanges constructor instead of &Changes{}
type Changes struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ChangesResp - changes response type with helper processor methods
type ChangesResp struct {
	Data        func() []*ChangeInfo         // gets change info collection
	GetNextPage func() (*ChangesResp, error) // gets next page changes collection
	data        []*ChangeInfo
}

// ChangeInfo - changes API response payload structure
type ChangeInfo struct {
	ChangeToken       *StringValue `json:"ChangeToken"`
	ChangeType        int          `json:"ChangeType"` // see more https://docs.microsoft.com/en-us/previous-versions/office/sharepoint-csom/ee543793(v%3Doffice.15)
	Editor            string       `json:"Editor"`
	EditorEmailHint   string       `json:"EditorEmailHint"`
	ItemID            int          `json:"ItemId"`
	ListID            string       `json:"ListId"`
	ServerRelativeURL string       `json:"ServerRelativeUrl"`
	SharedByUser      string       `json:"SharedByUser"`
	SharedWithUsers   string       `json:"SharedWithUsers"`
	SiteID            string       `json:"SiteId"`
	Time              time.Time    `json:"Time"`
	UniqueID          string       `json:"UniqueId"`
	WebID             string       `json:"WebId"`
}

// ChangeQuery ...
type ChangeQuery struct {
	ChangeTokenStart      string // Specifies the start date and start time for changes that are returned through the query
	ChangeTokenEnd        string // Specifies the end date and end time for changes that are returned through the query
	Add                   bool   // Specifies whether add changes are included in the query
	Alert                 bool   // Specifies whether changes to alerts are included in the query
	ContentType           bool   // Specifies whether changes to content types are included in the query
	DeleteObject          bool   // Specifies whether deleted objects are included in the query
	Field                 bool   // Specifies whether changes to fields are included in the query
	File                  bool   // Specifies whether changes to files are included in the query
	Folder                bool   // Specifies whether changes to folders are included in the query
	Group                 bool   // Specifies whether changes to groups are included in the query
	GroupMembershipAdd    bool   // Specifies whether adding users to groups is included in the query
	GroupMembershipDelete bool   // Specifies whether deleting users from the groups is included in the query
	Item                  bool   // Specifies whether general changes to list items are included in the query
	List                  bool   // Specifies whether changes to lists are included in the query
	Move                  bool   // Specifies whether move changes are included in the query
	Navigation            bool   // Specifies whether changes to the navigation structure of a site collection are included in the query
	Rename                bool   // Specifies whether renaming changes are included in the query
	Restore               bool   // Specifies whether restoring items from the recycle bin or from backups is included in the query
	RoleAssignmentAdd     bool   // Specifies whether adding role assignments is included in the query
	RoleAssignmentDelete  bool   // Specifies whether adding role assignments is included in the query
	RoleDefinitionAdd     bool   // Specifies whether adding role assignments is included in the query
	RoleDefinitionDelete  bool   // Specifies whether adding role assignments is included in the query
	RoleDefinitionUpdate  bool   // Specifies whether adding role assignments is included in the query
	SecurityPolicy        bool   // Specifies whether modifications to security policies are included in the query
	Site                  bool   // Specifies whether changes to site collections are included in the query
	SystemUpdate          bool   // Specifies whether updates made using the item SystemUpdate method are included in the query
	Update                bool   // Specifies whether update changes are included in the query
	User                  bool   // Specifies whether changes to users are included in the query
	View                  bool   // Specifies whether changes to views are included in the query
	Web                   bool   // Specifies whether changes to Web sites are included in the query
}

// NewChanges - Changes struct constructor function
func NewChanges(client *gosip.SPClient, endpoint string, config *RequestConfig) *Changes {
	return &Changes{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// GetCurrentToken gets current change token for this parent entity
func (changes *Changes) GetCurrentToken() (string, error) {
	endpoint := fmt.Sprintf("%s?$select=CurrentChangeToken", changes.endpoint)
	client := NewHTTPClient(changes.client)
	data, err := client.Get(endpoint, changes.config)
	if err != nil {
		return "", err
	}
	data = NormalizeODataItem(data)
	res := &struct {
		CurrentChangeToken StringValue `json:"CurrentChangeToken"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.CurrentChangeToken.StringValue, nil
}

// GetChanges gets changes in scope of the parent container using provided change query
func (changes *Changes) GetChanges(changeQuery *ChangeQuery) (*ChangesResp, error) {
	endpoint := toURL(fmt.Sprintf("%s/GetChanges", changes.endpoint), changes.modifiers)
	client := NewHTTPClient(changes.client)
	metadata := map[string]interface{}{}
	if changeQuery != nil {
		optsRaw, _ := json.Marshal(changeQuery)
		_ = json.Unmarshal(optsRaw, &metadata)
	}
	metadata["__metadata"] = map[string]string{"type": "SP.ChangeQuery"}
	if changeQuery != nil && changeQuery.ChangeTokenStart != "" {
		metadata["ChangeTokenStart"] = map[string]string{"StringValue": changeQuery.ChangeTokenStart}
	}
	if changeQuery != nil && changeQuery.ChangeTokenEnd != "" {
		metadata["ChangeTokenEnd"] = map[string]string{"StringValue": changeQuery.ChangeTokenEnd}
	}
	for k, v := range metadata {
		if v == false || v == "" || v == nil {
			delete(metadata, k)
		}
	}
	query := map[string]interface{}{"query": metadata}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	data, err := client.Post(endpoint, bytes.NewBuffer(body), changes.config)
	if err != nil {
		return nil, err
	}

	result := &ChangesResp{}
	result.Data = func() []*ChangeInfo {
		result.data = nil
		collection, _ := normalizeODataCollection(data)
		var changesInfo []*ChangeInfo
		for _, changeItem := range collection {
			c := &ChangeInfo{}
			if err := json.Unmarshal(changeItem, &c); err == nil {
				changesInfo = append(changesInfo, c)
			}
		}
		result.data = changesInfo
		return changesInfo
	}
	result.GetNextPage = func() (*ChangesResp, error) {
		if result.data == nil {
			result.Data()
		}
		if len(result.data) == 0 {
			return nil, fmt.Errorf("can't get next page of an empty collection")
		}
		changeQuery.ChangeTokenStart = result.data[len(result.data)-1].ChangeToken.StringValue
		return changes.GetChanges(changeQuery)
	}

	return result, nil
}
