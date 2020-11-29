package models

import (
	"github.com/seknox/trasa/server/consts"
)

type ScimUser struct {
	Schemas          []string                   `json:"schemas"`
	ID               string                     `json:"id"`
	ExternalID       string                     `json:"externalId"`
	UserName         string                     `json:"userName"`
	UserRole         string                     `json:"userRole"`
	Name             ScimUserName               `json:"name"`
	Emails           []ScimUserEmails           `json:"emails"`
	Password         string                     `json:"password"`
	Groups           []ScimUserGroups           `json:"groups"`
	X509Certificates []ScimUserX509Certificates `json:"x509Certificates"`
	Active           bool                       `json:"active"`
	Meta             ScimMeta                   `json:"meta"`
}

type ScimListUser struct {
	Schemas      []string   `json:"schemas"`
	TotalResults int        `json:"totalResults"`
	ItemsPerPage int        `json:"itemsPerPage"`
	StartIndex   int        `json:"startIndex"`
	Resources    []ScimUser `json:"Resources"`
}

type ScimListGroup struct {
	Schemas      []string    `json:"schemas"`
	TotalResults int         `json:"totalResults"`
	ItemsPerPage int         `json:"itemsPerPage"`
	StartIndex   int         `json:"startIndex"`
	Resources    []ScimGroup `json:"Resources"`
}

type ScimUserName struct {
	GivenName  string `json:"givenName"`
	MiddleName string `json:"middleName"`
	FamilyName string `json:"familyName"`
}

type ScimUserEmails struct {
	Primary bool   `json:"primary"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type ScimUserX509Certificates struct {
	Value string `json:"value"`
}

type ScimUserGroups struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type ScimConflict struct {
	Schemas []string `json:"schemas"`
	Detail  string   `json:"detail"`
	Status  int      `json:"status"`
}

func (c ScimConflict) New(detail string) ScimConflict {
	c.Schemas = []string{consts.SCIM_ERR}
	c.Detail = detail
	c.Status = 409

	return c
}

type ScimGroup struct {
	Schemas     []string           `json:"schemas"`
	ID          string             `json:"id"`
	DisplayName string             `json:"displayName"`
	Members     []ScimGroupMembers `json:"members"`
	Meta        ScimMeta           `json:"meta"`
}

type ScimGroupMembers struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type ScimMeta struct {
	ResourceType string `json:"resourceType"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
	Version      string `json:"version"`
	Location     string `json:"location"`
}

type ScimGroupPatch struct {
	Schemas    []string            `json:"schemas"`
	Operations []ScimGroupPatchOps `json:"Operations"`
}

type ScimGroupPatchOps struct {
	Op    string             `json:"op"`
	Path  string             `json:"path"`
	Value []ScimGroupMembers `json:"value,omitempty"`
}
