package models

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

type UserGroup struct {
	MapID     string `json:"mapID"`
	GroupID   string `json:"groupID"`
	OrgID     string `json:"orgID"`
	UserID    string `json:"userID"`
	Status    bool   `json:"status"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updated_at"`
}

type ServiceGroupMap struct {
	MapID         string `json:"mapID"`
	GroupID       string `json:"groupID"`
	OrgID         string `json:"orgID"`
	AuthserviceID string `json:"authserviceID"`
	Status        bool   `json:"status"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updated_at"`
}

type GroupMap struct {
	MapID   string `json:"mapID"`
	GroupID string `json:"groupID"`
	OrgID   string `json:"orgID"`
	// EntityVal can be either userID, serviceID or policy string for current group in context.
	EntityVal string `json:"entityVal"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updated_at"`
}
