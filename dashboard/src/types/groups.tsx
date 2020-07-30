
export type GroupType ={
    groupID     :string
    orgID       :string
    groupType   :string
    groupName   :string
    status      :boolean
    memberCount :number
    createdAt   :number
    updatedAt   :number
}


export type PolicyType ={
    PolicyID         :string,
    OrgID            :string,
    PolicyName       :string,
    DayAndTime       :DayAndTimePolicy[],
    TfaRequired      :boolean,
    RecordSession    :boolean,
    FileTransfer     :boolean,
    IPSource         :string,
    AllowedCountries :string,
    DevicePolicy     :DevicePolicy,
    RiskThreshold    :number,
    CreatedAt        :number,
    UpdatedAt        :number,
    Expiry           :string,
    IsExpired        :boolean,
    UsedBy           :number,
}


export type DayAndTimePolicy={

}

type DevicePolicy={}