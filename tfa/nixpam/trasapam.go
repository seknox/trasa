package main

/*
#cgo LDFLAGS: -lpam -fPIC
#include <stdlib.h>
#include <security/pam_appl.h>


// utility functions defined in trasapamHelper.go . retreives values from pam handel
char *get_username(pam_handle_t *pamh);
char *get_rhost(pam_handle_t *pamh);

char *get_trasaID(pam_handle_t *pamh, int flags);
char *get_tfaval(pam_handle_t *pamh, int flags);



*/
import "C"

import (
	"fmt"
	"unsafe"
)

var ServiceConfig configFile

//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {

	var err error
	ServiceConfig, err = readConfigFromFile()
	if err != nil {
		writeLog(fmt.Sprintf("[ReadConfigFromFile] %v. Bypassing 2fa verification.", err))

	}

	// get username
	cUsername := C.get_username(pamh)
	if cUsername == nil {
		return C.PAM_USER_UNKNOWN
	}
	defer C.free(unsafe.Pointer(cUsername))

	// get remote user ip address
	cRHost := C.get_rhost(pamh)
	if cRHost == nil {
		return C.PAM_USER_UNKNOWN
	}
	defer C.free(unsafe.Pointer(cRHost))

	if ServiceConfig.TrasaPAMConfig.Debug {
		rh := fmt.Sprintf("cRHost result: %s", C.GoString(cRHost))
		writeLog(rh)
	}

	// get trasaID
	cTrasaID := C.get_trasaID(pamh, flags)
	if cTrasaID == nil {
		return C.PAM_USER_UNKNOWN
	}
	defer C.free(unsafe.Pointer(cTrasaID))
	check, cErr := checkAndReturnPAMerr(C.GoString(cTrasaID))
	if check == true {
		return cErr
	}

	if ServiceConfig.TrasaPAMConfig.Debug {
		tid := fmt.Sprintf("cTrasaID result: %s", C.GoString(cTrasaID))
		writeLog(tid)
	}

	// get tfaval. totpCode or empty for u2f.
	cTfaval := C.get_tfaval(pamh, flags)
	defer C.free(unsafe.Pointer(cTfaval))
	check, cErr = checkAndReturnPAMerr(C.GoString(cTfaval))
	if check == true {
		return cErr
	}

	if ServiceConfig.TrasaPAMConfig.Debug {
		tf := fmt.Sprintf("cTfaval result: %s", C.GoString(cTfaval))
		writeLog(tf)
	}

	// call tfaReuqest flow here
	tfaResp := sendTfaReq(C.GoString(cUsername), C.GoString(cTrasaID), C.GoString(cTfaval), C.GoString(cRHost))

	if ServiceConfig.TrasaPAMConfig.Debug {
		f := fmt.Sprintf("final result: %v", tfaResp)
		writeLog(f)
	}

	if tfaResp == true {
		return C.PAM_SUCCESS
	}

	return C.PAM_AUTH_ERR

}

//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	return C.PAM_IGNORE
}

// since we return char * from c function, we use this function to check and return PAM module error
func checkAndReturnPAMerr(err string) (bool, C.int) {
	switch err {
	case "pam_auth_err":
		return true, C.PAM_AUTH_ERR
	case "pam_conv_err":
		return true, C.PAM_CONV_ERR
	case "cr":
		return true, C.PAM_CONV_ERR
	}

	return false, C.PAM_SUCCESS
}

func main() {

}
