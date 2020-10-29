package http

import (
	"strings"
)

// findAndReplaceUnamePass finds username and password field in string and replaces it with supplied values.
// We expect and recognize specefic hardcoded strings to be a email and password field.
// specifically username to be "trasauser@trasa.io" and password to be "trasapassword".
func findAndReplaceUnamePass(body, username, password string, shouldFIllUser, shouldFIllPass bool) string {
	//fmt.Println("body Bytes BEFORE: ", body)
	if shouldFIllUser && shouldFIllPass {
		// find and replace username field
		modeifiedBody := strings.Replace(body, "trasauser@trasa.io", username, -1)
		// find and replace username field
		modeifiedBody = strings.Replace(modeifiedBody, "trasapassword", password, -1)
		return modeifiedBody
	}

	if shouldFIllPass {

		modeifiedBody := strings.Replace(body, "trasapassword'", password, -1)
		return modeifiedBody
	}

	return body

}
