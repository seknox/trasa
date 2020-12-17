package consts

//Failed reasons
type FailedReason string

const (
	REASON_MALFORMED_REQUEST_RECEIVED FailedReason = "REASON_MALFORMED_REQUEST_RECEIVED"

	REASON_INVALID_SERVICE_HOSTNAME    FailedReason = "REASON_INVALID_SERVICE_HOSTNAME"
	REASON_INVALID_SERVICE_CREDS       FailedReason = "INVALID_SERVICE_CREDS"
	REASON_INVALID_SERVICE_ID          FailedReason = "INVALID_SERVICE_ID"
	REASON_USER_DISABLED               FailedReason = "USER_DISABLED"
	REASON_INVALID_USER_CREDS          FailedReason = "INVALID_USER_CREDS"
	REASON_USER_NOT_FOUND              FailedReason = "USER_NOT_FOUND"
	REASON_ORG_NOT_FOUND               FailedReason = "ORG_NOT_FOUND"
	REASON_INVALID_PUBLIC_KEY          FailedReason = "INVALID_PUBLIC_KEY"
	REASON_INVALID_CERTIFICATE         FailedReason = "INVALID_CERTIFICATE"
	REASON_IDENTITY_PROVIDER_NOT_FOUND FailedReason = "IDENTITY_PROVIDER_NOT_FOUND"
	REASON_INVALID_TOKEN               FailedReason = "INVALID_TOKEN"
	REASON_SPOOFED_LOGIN               FailedReason = "SPOOFED_LOGIN"
	REASON_INVALID_TOTP                FailedReason = "INVALID_TOTP"
	REASON_U2F_FAILED                  FailedReason = "U2F_FAILED"
	REASON_U2FY_FAILED                 FailedReason = "U2FY_FAILED"
	REASON_TIME_POLICY_FAILED          FailedReason = "TIME_POLICY_FAILED"
	REASON_IP_POLICY_FAILED            FailedReason = "IP_POLICY_FAILED"
	REASON_POLICY_EXPIRED              FailedReason = "POLICY_EXPIRED"
	REASON_ADHOC_POLICY_FAILED         FailedReason = "ADHOC_POLICY_FAILED"
	REASON_NO_POLICY_ASSIGNED          FailedReason = "NO_POLICY_ASSIGNED"

	REASON_INVALID_PRIVILEGE      FailedReason = "INVALID_PRIVILEGE"
	REASON_LDAP_AUTH_FAILED       FailedReason = "LDAP_AUTH_FAILED"
	REASON_UNKNOWN                FailedReason = "UNKNOWN"
	REASON_TRASA_ERROR            FailedReason = "TRASA_ERROR"
	REASON_HOST_NOT_REACHABLE     FailedReason = "HOST_NOT_REACHABLE"
	REASON_DYNAMIC_SERVICE_FAILED FailedReason = "DYNAMIC_SERVICE_FAILED"
	REASON_COUNTRY_POLICY_FAILED  FailedReason = "COUNTRY_POLICY_FAILED"

	REASON_DEVICE_POLICY_FAILED FailedReason = "DEVICE_POLICY_FAILED"
	REASON_DEVICE_NOT_ENROLLED  FailedReason = "DEVICE_NOT_ENROLLED"

	REASON_PASSWORD_POLICY_FAILED    FailedReason = "PASSWORD_POLICY_FAILED"
	REASON_INVALID_HOST_KEY          FailedReason = "INVALID_HOST_KEY"
	REASON_FILE_TRANSFER_NOT_ALLOWED FailedReason = "FILE_TRANSFER_NOT_ALLOWED"

	REASON_FAILED_TO_GENERATE_TOKEN FailedReason = "FAILED_TO_GENERATE_TOKEN"
)

//Login methods
