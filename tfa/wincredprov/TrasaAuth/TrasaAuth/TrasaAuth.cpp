/*++
This code is subjected to Copyright.
Owner of this code is sakshyam@seknox.com
Seknox Cybersecurity Pvt. Ltd.
--*/


#if ( _MSC_VER >= 800 )
#pragma warning ( 3 : 4100 ) // enable "Unreferenced formal parameter"
#pragma warning ( 3 : 4219 ) // enable "trailing ',' used for variable argument list"
#endif


#ifndef WIN32_NO_STATUS
#include <ntstatus.h>
#define WIN32_NO_STATUS
#endif
#include <unknwn.h>

///////////////////
#include "stdafx.h"
#define _UNICODE
#define UNICODE
#include <windef.h>
#include <windows.h>
#include <lmcons.h>
#include <lmaccess.h>
#include <lmapibuf.h>
#include <subauth.h>
#include <stdio.h>
#include <winnt.h>
#include <string>
#include <winhttp.h>
#include <iostream>

#include <vector>
#include <strsafe.h>
#include <tchar.h>

#include "trasautils.h"

//#pragma warning(disable : 4146)
//#pragma comment(lib, "winhttp.lib")





NTSTATUS
NTAPI
Msv1_0SubAuthenticationRoutine(
	IN NETLOGON_LOGON_INFO_CLASS LogonLevel,
	IN PVOID LogonInformation,
	IN ULONG Flags,
	IN PUSER_ALL_INFORMATION UserAll,
	OUT PULONG WhichFields,
	OUT PULONG UserFlags,
	OUT PBOOLEAN Authoritative,
	OUT PLARGE_INTEGER LogoffTime,
	OUT PLARGE_INTEGER KickoffTime
)

{



	NTSTATUS Status;
	SYSTEMTIME CurrentTime;
	WCHAR buf[256];
	PNETLOGON_LOGON_IDENTITY_INFO Identity =
		(PNETLOGON_LOGON_IDENTITY_INFO)LogonInformation;
	//NETLOGON_LOGON_IDENTITY_INFO val = (NETLOGON_LOGON_IDENTITY_INFO)LogonInformation;


	/////////////////////////////////////////////
	////////////// http request /////////////////
	/////////////////////////////////////////////

	// get value of username
	WCHAR usernameFromLogon[256];
	WCHAR workstationFromLogon[256];
	swprintf_s(usernameFromLogon, RTL_NUMBER_OF(usernameFromLogon), L"%wZ", &Identity->UserName);
	swprintf_s(workstationFromLogon, RTL_NUMBER_OF(workstationFromLogon), L"%wZ", &Identity->Workstation);
	// change WCHAR to string
	std::wstring un(usernameFromLogon);
	std::string userName(un.begin(), un.end());

	std::wstring ws(workstationFromLogon);
	std::string workstation(ws.begin(), ws.end());


	// default status is account locked out.
	Status = STATUS_ACCOUNT_LOCKED_OUT;

	// change status success here
	if (sendRequest(userName, workstation).compare("success") == 0)
	{
		Status = STATUS_SUCCESS;
	}

	

	*Authoritative = TRUE;
	*UserFlags = 0;
	*WhichFields = 0;

	/////////////////////////////////////
	GetLocalTime(&CurrentTime);

	if (!Identity) {
		//WriteLogFile(TEXT("No identity\r\n"));
		return Status;
	}


	swprintf_s(buf, RTL_NUMBER_OF(buf),
		L"%02d/%02d/%d %02d:%02d:%02d: Logon (level=%d) %wZ\\%wZ (%wZ) from %wZ\r\n",
		CurrentTime.wMonth, CurrentTime.wDay, CurrentTime.wYear,
		CurrentTime.wHour, CurrentTime.wMinute, CurrentTime.wSecond,
		LogonLevel,
		&Identity->LogonDomainName, &Identity->UserName,
		&UserAll->FullName, &Identity->Workstation);
	WriteLogFile(buf);

	switch (LogonLevel) {
	case NetlogonInteractiveInformation:
	case NetlogonServiceInformation:
	case NetlogonNetworkInformation:

		//
		// If you care you can determine what to do here
		//
		*Authoritative = FALSE;

		if (LogoffTime) {
			LogoffTime->HighPart = 0x7FFFFFFF;
			LogoffTime->LowPart = 0xFFFFFFFF;
		}

		if (KickoffTime) {
			KickoffTime->HighPart = 0x7FFFFFFF;
			KickoffTime->LowPart = 0xFFFFFFFF;
		}
		break;

	default:
		return STATUS_INVALID_INFO_CLASS;
	}

	return Status;
}




NTSTATUS
NTAPI
Msv1_0SubAuthenticationFilter(
	IN NETLOGON_LOGON_INFO_CLASS LogonLevel,
	IN PVOID LogonInformation,
	IN ULONG Flags,
	IN PUSER_ALL_INFORMATION UserAll,
	OUT PULONG WhichFields,
	OUT PULONG UserFlags,
	OUT PBOOLEAN Authoritative,
	OUT PLARGE_INTEGER LogoffTime,
	OUT PLARGE_INTEGER KickoffTime
)
{

	NTSTATUS Status;
	SYSTEMTIME CurrentTime;
	WCHAR buf[256];
	PNETLOGON_LOGON_IDENTITY_INFO Identity =
		(PNETLOGON_LOGON_IDENTITY_INFO)LogonInformation;
	//NETLOGON_LOGON_IDENTITY_INFO val = (NETLOGON_LOGON_IDENTITY_INFO)LogonInformation;


	/////////////////////////////////////////////
	////////////// http request /////////////////
	/////////////////////////////////////////////

	// get value of username
	WCHAR usernameFromLogon[256];
	WCHAR workstationFromLogon[256];
	swprintf_s(usernameFromLogon, RTL_NUMBER_OF(usernameFromLogon), L"%wZ", &Identity->UserName);
	swprintf_s(workstationFromLogon, RTL_NUMBER_OF(workstationFromLogon), L"%wZ", &Identity->Workstation);
	// change WCHAR to string
	std::wstring un(usernameFromLogon);
	std::string userName(un.begin(), un.end());

	std::wstring ws(workstationFromLogon);
	std::string workstation(ws.begin(), ws.end());


	// default status is account locked out.
	Status = STATUS_ACCOUNT_LOCKED_OUT;

	// change status success here
	if (sendRequest(userName, workstation).compare("success") == 0)
	{
		Status = STATUS_SUCCESS;
	}



	*Authoritative = TRUE;
	*UserFlags = 0;
	*WhichFields = 0;

	/////////////////////////////////////
	GetLocalTime(&CurrentTime);

	if (!Identity) {
		//WriteLogFile(TEXT("No identity\r\n"));
		return Status;
	}


	swprintf_s(buf, RTL_NUMBER_OF(buf),
		L"%02d/%02d/%d %02d:%02d:%02d: Logon (level=%d) %wZ\\%wZ (%wZ) from %wZ\r\n",
		CurrentTime.wMonth, CurrentTime.wDay, CurrentTime.wYear,
		CurrentTime.wHour, CurrentTime.wMinute, CurrentTime.wSecond,
		LogonLevel,
		&Identity->LogonDomainName, &Identity->UserName,
		&UserAll->FullName, &Identity->Workstation);
	WriteLogFile(buf);

	switch (LogonLevel) {
	case NetlogonInteractiveInformation:
	case NetlogonServiceInformation:
	case NetlogonNetworkInformation:

		//
		// If you care you can determine what to do here
		//
		*Authoritative = FALSE;

		if (LogoffTime) {
			LogoffTime->HighPart = 0x7FFFFFFF;
			LogoffTime->LowPart = 0xFFFFFFFF;
		}

		if (KickoffTime) {
			KickoffTime->HighPart = 0x7FFFFFFF;
			KickoffTime->LowPart = 0xFFFFFFFF;
		}
		break;

	default:
		return STATUS_INVALID_INFO_CLASS;
	}

	return Status;
}