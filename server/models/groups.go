package models

//Group can be user or service group
type Group struct {
	GroupID     string `json:"groupID"`
	OrgID       string `json:"orgID"`
	GroupType   string `json:"groupType"`
	GroupName   string `json:"groupName"`
	Status      bool   `json:"status"`
	MemberCount int    `json:"memberCount"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

//UserGroupMap is a database relation map between user and group
type UserGroupMap struct {
	MapID     string `json:"mapID"`
	GroupID   string `json:"groupID"`
	OrgID     string `json:"orgID"`
	UserID    string `json:"userID"`
	Status    bool   `json:"status"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updated_at"`
}

//ServiceGroupMap is a database relation map between Service and Group
type ServiceGroupMap struct {
	MapID         string `json:"mapID"`
	GroupID       string `json:"groupID"`
	OrgID         string `json:"orgID"`
	AuthserviceID string `json:"authserviceID"`
	Status        bool   `json:"status"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updated_at"`
}
