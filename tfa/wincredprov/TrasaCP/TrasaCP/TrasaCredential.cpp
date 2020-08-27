//
// THIS CODE AND INFORMATION IS PROVIDED "AS IS" WITHOUT WARRANTY OF
// ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A
// PARTICULAR PURPOSE.
//
// Copyright (c) 2018 Seknox Cybersecurity Pvt. Ltd.
//
//

#ifndef WIN32_NO_STATUS
#include <ntstatus.h>
#define WIN32_NO_STATUS
#endif
#include <unknwn.h>
#include <wincred.h>
#include "TrasaCredential.h"
//#include "TrasaProvider.h"
#include "guid.h"



#include "TrasaUtils.h"

# include <thread>
#include <cstdio>
#include <Windows.h>


using CSharpForm = LPWSTR(__stdcall*)(LPWSTR userName);


// TrasaCredential ////////////////////////////////////////////////////////

TrasaCredential::TrasaCredential():
    _cRef(1),
    _pCredProvCredentialEvents(NULL)
{
    DllAddRef();

    ZeroMemory(_rgCredProvFieldDescriptors, sizeof(_rgCredProvFieldDescriptors));
    ZeroMemory(_rgFieldStatePairs, sizeof(_rgFieldStatePairs));
    ZeroMemory(_rgFieldStrings, sizeof(_rgFieldStrings));
}

TrasaCredential::~TrasaCredential()
{
    if (_rgFieldStrings[SFI_PASSWORD])
    {
        // CoTaskMemFree (below) deals with NULL, but StringCchLength does not.
        size_t lenPassword;
        HRESULT hr = StringCchLengthW(_rgFieldStrings[SFI_PASSWORD], 128, &(lenPassword));
        if (SUCCEEDED(hr))
        {
            SecureZeroMemory(_rgFieldStrings[SFI_PASSWORD], lenPassword * sizeof(*_rgFieldStrings[SFI_PASSWORD]));
        }
        else
        {
            // TODO: Determine how to handle count error here.
        }
    }
    for (int i = 0; i < ARRAYSIZE(_rgFieldStrings); i++)
    {
        CoTaskMemFree(_rgFieldStrings[i]);
        CoTaskMemFree(_rgCredProvFieldDescriptors[i].pszLabel);
    }

    DllRelease();
}

// Initializes one credential with the field information passed in.
// Set the value of the SFI_USERNAME field to pwzUsername.
HRESULT TrasaCredential::Initialize(
                                      CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus,
                                      const CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR* rgcpfd,
                                      const FIELD_STATE_PAIR* rgfsp,
                                      DWORD dwFlags,
                                      PCWSTR pwzUsername,
                                      PCWSTR pwzPassword
                                      )
{
    HRESULT hr = S_OK;
    _cpus = cpus;
    _dwFlags = dwFlags;
    // Copy the field descriptors for each field. This is useful if you want to vary the 
    // field descriptors based on what Usage scenario the credential was created for.
    for (DWORD i = 0; SUCCEEDED(hr) && i < ARRAYSIZE(_rgCredProvFieldDescriptors); i++)
    {
        _rgFieldStatePairs[i] = rgfsp[i];
        hr = FieldDescriptorCopy(rgcpfd[i], &_rgCredProvFieldDescriptors[i]);
    }

    // Initialize the String values of all the fields.
	if (SUCCEEDED(hr))
	{
		hr = SHStrDupW(L"This system is protected with TRASA. Use your local or domain account for authentication.", &_rgFieldStrings[SFI_LARGE_TEXT]);
	}
    if (SUCCEEDED(hr))
    {
        hr = SHStrDupW(pwzUsername, &_rgFieldStrings[SFI_USERNAME]);
    }
    if (SUCCEEDED(hr))
    {
        hr = SHStrDupW(pwzPassword ? pwzPassword : L"", &_rgFieldStrings[SFI_PASSWORD]);
    }
	/*if (SUCCEEDED(hr))
    {
        hr = SHStrDupW( L"", &_rgFieldStrings[SFI_TOTP]);
    }*/
    if (SUCCEEDED(hr))
    {
        hr = SHStrDupW(L"Submit", &_rgFieldStrings[SFI_SUBMIT_BUTTON]);
    }

    return S_OK;
}

// LogonUI calls this in order to give us a callback in case we need to notify it of anything.
HRESULT TrasaCredential::Advise(ICredentialProviderCredentialEvents* pcpce)
{
    if (_pCredProvCredentialEvents != NULL)
    {
        _pCredProvCredentialEvents->Release();
    }
    _pCredProvCredentialEvents = pcpce;
    _pCredProvCredentialEvents->AddRef();
    return S_OK;
}

// LogonUI calls this to tell us to release the callback.
HRESULT TrasaCredential::UnAdvise()
{
    if (_pCredProvCredentialEvents)
    {
        _pCredProvCredentialEvents->Release();
    }
    _pCredProvCredentialEvents = NULL;
    return S_OK;
}

// LogonUI calls this function when our tile is selected (zoomed).
// If you simply want fields to show/hide based on the selected state,
// there's no need to do anything here - you can set that up in the 
// field definitions.  But if you want to do something
// more complicated, like change the contents of a field when the tile is
// selected, you would do it here.
HRESULT TrasaCredential::SetSelected(BOOL* pbAutoLogon)  
{
    *pbAutoLogon = FALSE;  

    return S_OK;
}

// Similarly to SetSelected, LogonUI calls this when your tile was selected
// and now no longer is. The most common thing to do here (which we do below)
// is to clear out the password field.
HRESULT TrasaCredential::SetDeselected()
{
    HRESULT hr = S_OK;
    if (_rgFieldStrings[SFI_PASSWORD])
    {
        // CoTaskMemFree (below) deals with NULL, but StringCchLength does not.
        size_t lenPassword;
        hr = StringCchLengthW(_rgFieldStrings[SFI_PASSWORD], 128, &(lenPassword));
        if (SUCCEEDED(hr))
        {
            SecureZeroMemory(_rgFieldStrings[SFI_PASSWORD], lenPassword * sizeof(*_rgFieldStrings[SFI_PASSWORD]));

            CoTaskMemFree(_rgFieldStrings[SFI_PASSWORD]);
            hr = SHStrDupW(L"", &_rgFieldStrings[SFI_PASSWORD]);
        }

        if (SUCCEEDED(hr) && _pCredProvCredentialEvents)
        {
            _pCredProvCredentialEvents->SetFieldString(this, SFI_PASSWORD, _rgFieldStrings[SFI_PASSWORD]);
        }
    }

    return hr;
}

// Gets info for a particular field of a tile. Called by logonUI to get information to 
// display the tile.
HRESULT TrasaCredential::GetFieldState(
    DWORD dwFieldID,
    CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs,
    CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE* pcpfis
    )
{
    HRESULT hr;

    // Validate paramters.
    if ((dwFieldID < ARRAYSIZE(_rgFieldStatePairs)) && pcpfs && pcpfis)
    {
        *pcpfs = _rgFieldStatePairs[dwFieldID].cpfs;
        *pcpfis = _rgFieldStatePairs[dwFieldID].cpfis;

        hr = S_OK;
    }
    else
    {
        hr = E_INVALIDARG;
    }
    return hr;
}

// Sets ppwz to the string value of the field at the index dwFieldID.
HRESULT TrasaCredential::GetStringValue(
    DWORD dwFieldID, 
    PWSTR* ppwz
    )
{
    HRESULT hr;

    // Check to make sure dwFieldID is a legitimate index.
    if (dwFieldID < ARRAYSIZE(_rgCredProvFieldDescriptors) && ppwz) 
    {
        // Make a copy of the string and return that. The caller
        // is responsible for freeing it.
        hr = SHStrDupW(_rgFieldStrings[dwFieldID], ppwz);
    }
    else
    {
        hr = E_INVALIDARG;
    }

    return hr;
}

// Gets the image to show in the user tile.
HRESULT TrasaCredential::GetBitmapValue(
    DWORD dwFieldID, 
    HBITMAP* phbmp
    )
{
    HRESULT hr;
    if ((SFI_TILEIMAGE == dwFieldID) && phbmp)
    {
        HBITMAP hbmp = LoadBitmap(HINST_THISDLL, MAKEINTRESOURCE(IDB_TILE_IMAGE));
        if (hbmp != NULL)
        {
            hr = S_OK;
            *phbmp = hbmp;
        }
        else
        {
            hr = HRESULT_FROM_WIN32(GetLastError());
        }
    }
    else
    {
        hr = E_INVALIDARG;
    }

    return hr;
}

// Sets pdwAdjacentTo to the index of the field the submit button should be 
// adjacent to. We recommend that the submit button is placed next to the last
// field which the user is required to enter information in. Optional fields
// should be below the submit button.
HRESULT TrasaCredential::GetSubmitButtonValue(
    DWORD dwFieldID,
    DWORD* pdwAdjacentTo
    )
{
    HRESULT hr;

    // Validate parameters.
    if ((SFI_SUBMIT_BUTTON == dwFieldID) && pdwAdjacentTo)
    {
        // pdwAdjacentTo is a pointer to the fieldID you want the submit button to appear next to.
        *pdwAdjacentTo = SFI_PASSWORD;
        hr = S_OK;
    }
    else
    {
        hr = E_INVALIDARG;
    }
    return hr;
}

// Sets the value of a field which can accept a string as a value.
// This is called on each keystroke when a user types into an edit field.
HRESULT TrasaCredential::SetStringValue(
    DWORD dwFieldID, 
    PCWSTR pwz      
    )
{
    HRESULT hr;

    // Validate parameters.
    if (dwFieldID < ARRAYSIZE(_rgCredProvFieldDescriptors) && 
       (CPFT_EDIT_TEXT == _rgCredProvFieldDescriptors[dwFieldID].cpft || 
        CPFT_PASSWORD_TEXT == _rgCredProvFieldDescriptors[dwFieldID].cpft)) 
    {
        PWSTR* ppwzStored = &_rgFieldStrings[dwFieldID];
        CoTaskMemFree(*ppwzStored);
        hr = SHStrDupW(pwz, ppwzStored);
    }
    else
    {
        hr = E_INVALIDARG;
    }

    return hr;
}

//------------- 
// The following methods are for logonUI to get the values of various UI elements and then communicate
// to the credential about what the user did in that field.  However, these methods are not implemented
// because our tile doesn't contain these types of UI elements
HRESULT TrasaCredential::GetCheckboxValue(
    DWORD dwFieldID, 
    BOOL* pbChecked,
    PWSTR* ppwzLabel
    )
{
    UNREFERENCED_PARAMETER(dwFieldID);
    UNREFERENCED_PARAMETER(pbChecked);
    UNREFERENCED_PARAMETER(ppwzLabel);

    return E_NOTIMPL;
}

HRESULT TrasaCredential::GetComboBoxValueCount(
    DWORD dwFieldID, 
    DWORD* pcItems, 
    DWORD* pdwSelectedItem
    )
{
    UNREFERENCED_PARAMETER(dwFieldID);
    UNREFERENCED_PARAMETER(pcItems);
    UNREFERENCED_PARAMETER(pdwSelectedItem);
    return E_NOTIMPL;
}

HRESULT TrasaCredential::GetComboBoxValueAt(
    DWORD dwFieldID, 
    DWORD dwItem,
    PWSTR* ppwzItem
    )
{
    UNREFERENCED_PARAMETER(dwFieldID);
    UNREFERENCED_PARAMETER(dwItem);
    UNREFERENCED_PARAMETER(ppwzItem);
    return E_NOTIMPL;
}

HRESULT TrasaCredential::SetCheckboxValue(
    DWORD dwFieldID, 
    BOOL bChecked
    )
{
    UNREFERENCED_PARAMETER(dwFieldID);
    UNREFERENCED_PARAMETER(bChecked);

    return E_NOTIMPL;
}

HRESULT TrasaCredential::SetComboBoxSelectedValue(
    DWORD dwFieldId,
    DWORD dwSelectedItem
    )
{
    UNREFERENCED_PARAMETER(dwFieldId);
    UNREFERENCED_PARAMETER(dwSelectedItem);
    return E_NOTIMPL;
}

HRESULT TrasaCredential::CommandLinkClicked(DWORD dwFieldID)
{
    UNREFERENCED_PARAMETER(dwFieldID);
    return E_NOTIMPL;
}
//------ end of methods for controls we don't have in our tile ----//







void SeparateUserAndDomainName(

	__in wchar_t* domain_slash_username,

	__out wchar_t* username,

	__in int sizeUsername,

	__out_opt wchar_t* domain,

	__in_opt int sizeDomain

)

{

	int pos;

	for (pos = 0; domain_slash_username[pos] != L'\\' && domain_slash_username[pos] != NULL; pos++);



	if (domain_slash_username[pos] != NULL)

	{

		int i;

		for (i = 0; i < pos && i < sizeDomain; i++)

			domain[i] = domain_slash_username[i];

domain[i] = L'\0';



for (i = 0; domain_slash_username[pos + i + 1] != NULL && i < sizeUsername; i++)

	username[i] = domain_slash_username[pos + i + 1];

username[i] = L'\0';

	}

	else

	{

	int i;

	for (i = 0; i < pos && i < sizeUsername; i++)

		username[i] = domain_slash_username[i];

	username[i] = L'\0';

	}

}


//WCHAR suf[100];
//swprintf(suf, 100, L"user of %s | %s\r\n", useru, domainu);
//WriteLogFile(suf);
//WriteLogFile(L"\r\n---\r\n");

BOOL CheckCreds(LPWSTR useru, LPWSTR domainu, LPWSTR passu, HANDLE hToken) {
	if (LogonUserW(useru, domainu, passu, LOGON32_LOGON_NETWORK, LOGON32_PROVIDER_DEFAULT, &hToken))
	{
		return true;
	}
	else {
		return false;
	}
}

// Collect the username and password into a serialized credential for the correct usage scenario 
// (logon/unlock is what's demonstrated in this sample).  LogonUI then passes these credentials 
// back to the system to log on.
HRESULT TrasaCredential::GetSerialization(
	CREDENTIAL_PROVIDER_GET_SERIALIZATION_RESPONSE* pcpgsr,
	CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcs,
	PWSTR* ppwzOptionalStatusText,
	CREDENTIAL_PROVIDER_STATUS_ICON* pcpsiOptionalStatusIcon
)
{
	UNREFERENCED_PARAMETER(ppwzOptionalStatusText);
	UNREFERENCED_PARAMETER(pcpsiOptionalStatusIcon);

	HRESULT hr;

	WCHAR dwsz[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	DWORD dcch = ARRAYSIZE(dwsz);

	WCHAR dlwsz[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	DWORD dlcch = ARRAYSIZE(dlwsz);

	//WCHAR lwsz[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	//DWORD lcch = ARRAYSIZE(lwsz);

	DWORD cb = 0;
	BYTE* rgb = NULL;




	WCHAR useru[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	WCHAR domainu[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	WCHAR domainlu[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	WCHAR sepDomain[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	WCHAR finalDomain[CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1];
	ULONG userbuf = CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1;
	ULONG domainbuf = CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1;
	HANDLE hToken = NULL;


	SeparateUserAndDomainName(_rgFieldStrings[SFI_USERNAME], useru, userbuf, sepDomain, domainbuf);

	if (GetComputerNameExW(ComputerNameDnsDomain, dwsz, &dcch)) {
		wcsncpy_s(domainu, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1, dwsz, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1);
	}

	if (GetComputerNameExW(ComputerNameDnsHostname, dlwsz, &dlcch)) {
		wcsncpy_s(domainlu, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1, dlwsz, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1);
	}

	if ((wcscmp(domainu, sepDomain) == 0 || wcscmp(domainlu, sepDomain) == 0) ) {
		wcsncpy_s(finalDomain, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1, sepDomain, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1);

		if (LogonUserW(useru, sepDomain, _rgFieldStrings[SFI_PASSWORD], LOGON32_LOGON_NETWORK, LOGON32_PROVIDER_DEFAULT, &hToken))
		{
			HMODULE mod = LoadLibraryA("TrasaTfaPrompt.dll");
			CSharpForm form = reinterpret_cast<CSharpForm>(GetProcAddress(mod, "TrasaTfaPrompt"));
			//std::thread f1(form);
			LPWSTR resp = L"failed";
			WCHAR uubuf[100];
			swprintf(uubuf, 100, L"%s\\%s", sepDomain, useru);
			LPWSTR userName = uubuf;
			resp = form(userName);
			if (wcscmp(resp, L"success") != 0) {
				//const wchar_t *inv = L"invalid invalid";
				//PWSTR* in =	"invalid invalid";
				*ppwzOptionalStatusText = L"failed 2fa";
				*pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
				::MessageBox(NULL, "Failed 2FA", "Alert", 0);

				return S_OK;
			}
		}
		/* // removing else because its not needed ? @sshahcodes 24th June 2019.
		else {
			*ppwzOptionalStatusText = L"invalid username or password";
			*pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
			::MessageBox(NULL, "Invalid Username or Password for provided domain", "Alert", 0);
			return S_OK;
		}*/
	}
	else {
		// if we are here, it means no domain was provider
	//	WriteLogFile(L"\r\n--no domain provided--\r\n");
		if (CheckCreds(useru, domainu, _rgFieldStrings[SFI_PASSWORD], &hToken)) {
			HMODULE mod = LoadLibraryA("TrasaTfaPrompt.dll");
			CSharpForm form = reinterpret_cast<CSharpForm>(GetProcAddress(mod, "TrasaTfaPrompt"));
			//std::thread f1(form);
			LPWSTR resp = L"failed";
			
			wcsncpy_s(finalDomain, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1, domainu, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1);
			WCHAR uubuf[100];
			swprintf(uubuf, 100, L"%s\\%s", domainu, useru);
			// sending only username and discarding domain name. @sshahcodes 24th June 2019.
			LPWSTR userName = uubuf; // _rgFieldStrings[SFI_USERNAME];
			resp = form(userName);
			if (wcscmp(resp, L"success") != 0) {
				//const wchar_t *inv = L"invalid invalid";
				//PWSTR* in =	"invalid invalid";
				*ppwzOptionalStatusText = L"failed 2fa";
				*pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
				::MessageBox(NULL, "Failed 2FA", "Alert", 0);

				return S_OK;
			}
		}
		else {
			if (CheckCreds(useru, domainlu, _rgFieldStrings[SFI_PASSWORD], &hToken)) {
				wcsncpy_s(finalDomain, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1, domainlu, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1);
				HMODULE mod = LoadLibraryA("TrasaTfaPrompt.dll");
				CSharpForm form = reinterpret_cast<CSharpForm>(GetProcAddress(mod, "TrasaTfaPrompt"));
				//std::thread f1(form);
				LPWSTR resp = L"failed";
				WCHAR uubuf[100];
				swprintf(uubuf, 100, L"%s\\%s", domainlu, useru);
				// sending only username and discarding domain name. @sshahcodes 24th June 2019.
				LPWSTR userName = uubuf; // _rgFieldStrings[SFI_USERNAME];
				resp = form(userName);
				if (wcscmp(resp, L"success") != 0) {
					//const wchar_t *inv = L"invalid invalid";
					//PWSTR* in =	"invalid invalid";
					*ppwzOptionalStatusText = L"failed 2fa";
					*pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
					::MessageBox(NULL, "Failed 2FA", "Alert", 0);

					return S_OK;
				}
			}
			// removing else because its not needed ? @sshahcodes 24th June 2019.
			//else {
				// DWORD last = GetLastError();
				// WCHAR buf[100];
				// swprintf(buf, 100, L"user pass failed for %s , %s, %s with status %d \r\n", useru, domainu, _rgFieldStrings[SFI_PASSWORD], last);
				// WriteLogFile(L"\r\n---\r\n");
				// WriteLogFile(buf);
				// WriteLogFile(L"\r\n---\r\n");
			/*	*ppwzOptionalStatusText = L"invalid username or password";
				*pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
				::MessageBox(NULL, "Invalid Username or Password for provided ", "Alert", 0);
				return S_OK;*/
			//}
		}
		
	}
		///////////////////////////////////////////////////

		///////////////////////////////////////
// at last we check if domain name is still empty. in this case, we assign local domain to finalDomain as last resort
	if (wcscmp(finalDomain, L"") == 0) {
		wcsncpy_s(finalDomain, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1, domainlu, CREDUI_MAX_DOMAIN_TARGET_LENGTH + 1);
	}


   /* if (GetComputerNameW(wsz, &cch))
    {*/
        PWSTR pwzProtectedPassword;

        hr = ProtectIfNecessaryAndCopyPassword(_rgFieldStrings[SFI_PASSWORD], _cpus, &pwzProtectedPassword);

		//{
			// Only CredUI scenarios should use CredPackAuthenticationBuffer.  Custom packing logic is necessary for
			// logon and unlock scenarios in order to specify the correct MessageType.
			if (CPUS_CREDUI == _cpus)
			{
				if (SUCCEEDED(hr))
				{
					PWSTR pwzDomainUsername = NULL;
					//PWSTR pwzDomainUsername = NULL;
					hr = DomainUsernameStringAlloc(finalDomain, useru, &pwzDomainUsername);
					if (SUCCEEDED(hr))
					{
						// We use KERB_INTERACTIVE_UNLOCK_LOGON in both unlock and logon scenarios.  It contains a
						// KERB_INTERACTIVE_LOGON to hold the creds plus a LUID that is filled in for us by Winlogon
						// as necessary.
						if (!CredPackAuthenticationBufferW((CREDUIWIN_PACK_32_WOW & _dwFlags) ? CRED_PACK_WOW_BUFFER : 0, pwzDomainUsername, pwzProtectedPassword, rgb, &cb))
						{
							if (ERROR_INSUFFICIENT_BUFFER == GetLastError())
							{
								rgb = (BYTE*)HeapAlloc(GetProcessHeap(), 0, cb);
								if (rgb)
								{
									// If the CREDUIWIN_PACK_32_WOW flag is set we need to return 32 bit buffers to our caller we do this by 
									// passing CRED_PACK_WOW_BUFFER to CredPacAuthenticationBufferW.
									if (!CredPackAuthenticationBufferW((CREDUIWIN_PACK_32_WOW & _dwFlags) ? CRED_PACK_WOW_BUFFER : 0, pwzDomainUsername, pwzProtectedPassword, rgb, &cb))
									{
										HeapFree(GetProcessHeap(), 0, rgb);
										hr = HRESULT_FROM_WIN32(GetLastError());
									}
									else
									{
										hr = S_OK;
									}
								}
								else
								{
									hr = E_OUTOFMEMORY;
								}
							}
							else
							{
								hr = E_FAIL;
							}
							HeapFree(GetProcessHeap(), 0, pwzDomainUsername);
						}
						else
						{
							hr = E_FAIL;
						}
					}
					CoTaskMemFree(pwzProtectedPassword);
				}
			}
			else
			{

				KERB_INTERACTIVE_UNLOCK_LOGON kiul;

				// Initialize kiul with weak references to our credential.
				hr = KerbInteractiveUnlockLogonInit(finalDomain, useru, pwzProtectedPassword, _cpus, &kiul);

				if (SUCCEEDED(hr))
				{
					// We use KERB_INTERACTIVE_UNLOCK_LOGON in both unlock and logon scenarios.  It contains a
					// KERB_INTERACTIVE_LOGON to hold the creds plus a LUID that is filled in for us by Winlogon
					// as necessary.
					hr = KerbInteractiveUnlockLogonPack(kiul, &pcpcs->rgbSerialization, &pcpcs->cbSerialization);
				}
			}

			if (SUCCEEDED(hr))
			{
				ULONG ulAuthPackage;
				hr = RetrieveNegotiateAuthPackage(&ulAuthPackage);
				if (SUCCEEDED(hr))
				{
					pcpcs->ulAuthenticationPackage = ulAuthPackage;
					pcpcs->clsidCredentialProvider = CLSID_TrasaProvider;

					// In CredUI scenarios, we must pass back the buffer constructed with CredPackAuthenticationBuffer.
					if (CPUS_CREDUI == _cpus)
					{
						pcpcs->rgbSerialization = rgb;
						pcpcs->cbSerialization = cb;
					}

				

					// At this point the credential has created the serialized credential used for logon
					// By setting this to CPGSR_RETURN_CREDENTIAL_FINISHED we are letting logonUI know
					// that we have all the information we need and it should attempt to submit the 
					// serialized credential.
					*pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
				}
				else
				{
					HeapFree(GetProcessHeap(), 0, rgb);
				}
			}
		// } comment out this portion if we are ever to enable our U2F handler in CP    
 /*   }
    else
    {
        DWORD dwErr = GetLastError();
        hr = HRESULT_FROM_WIN32(dwErr);
    }*/

    return hr;
}
struct REPORT_RESULT_STATUS_INFO
{
    NTSTATUS ntsStatus;
    NTSTATUS ntsSubstatus;
    PWSTR     pwzMessage;
    CREDENTIAL_PROVIDER_STATUS_ICON cpsi;
};

static const REPORT_RESULT_STATUS_INFO s_rgLogonStatusInfo[] =
{
    { STATUS_LOGON_FAILURE, STATUS_SUCCESS, L"Incorrect username or password combination....", CPSI_ERROR, },
    { STATUS_ACCOUNT_RESTRICTION, STATUS_ACCOUNT_DISABLED, L"The account is disabled.", CPSI_WARNING },
};



// ReportResult is completely optional.Its purpose is to allow a credential to customize the string
// and the icon displayed in the case of a logon failure.  For example, we have chosen to 
// customize the error shown in the case of bad username/password and in the case of the account
// being disabled.
HRESULT TrasaCredential::ReportResult(
	NTSTATUS ntsStatus,
	NTSTATUS ntsSubstatus,
	PWSTR* ppwszOptionalStatusText,
	CREDENTIAL_PROVIDER_STATUS_ICON* pcpsiOptionalStatusIcon
)
{
	// WriteLogFile(L"not called\r\n");
	// WriteLogFile(L"\r\n -- \r\n");
//	HRESULT hr = E_ACCESSDENIED;
	*ppwszOptionalStatusText = NULL;
	*pcpsiOptionalStatusIcon = CPSI_NONE;

	DWORD dwStatusInfo = (DWORD)-1;

	// Look for a match on status and substatus.
	for (DWORD i = 0; i < ARRAYSIZE(s_rgLogonStatusInfo); i++)
	{
		if (s_rgLogonStatusInfo[i].ntsStatus == ntsStatus && s_rgLogonStatusInfo[i].ntsSubstatus == ntsSubstatus)
		{
			dwStatusInfo = i;
			break;
		}
	}

	if ((DWORD)-1 != dwStatusInfo)
	{
		if (SUCCEEDED(SHStrDupW(s_rgLogonStatusInfo[dwStatusInfo].pwzMessage, ppwszOptionalStatusText)))
		{
			*pcpsiOptionalStatusIcon = s_rgLogonStatusInfo[dwStatusInfo].cpsi;
		}
	}

	if (_pCredProvCredentialEvents)
	{
		_pCredProvCredentialEvents->SetFieldString(this, SFI_PASSWORD, L"");
	}

	/*if (!SUCCEEDED(HRESULT_FROM_NT(ntsStatus)))
	{
		if (_pCredProvCredentialEvents)
		{
			_pCredProvCredentialEvents->SetFieldString(this, SFI_PASSWORD, L"");
		}
		
	}*/


	// Since NULL is a valid value for *ppwszOptionalStatusText and *pcpsiOptionalStatusIcon
	// this function can't fail.
	return S_OK;
}



//HWND mainhndld;
		// _pCredProvCredentialEvents->OnCreatingWindow(&mainhndld);

		// UINT_PTR timerID;
		// SetTimer(mainhndld,             // handle to main window 
		// 	timerID,            // timer identifier 
		// 	063333,                // five-minute interval 
		// 	(TIMERPROC)TrasaPrompt);

		//HMODULE mod = LoadLibraryA("tfaForm.dll");
		//CSharpForm form = reinterpret_cast<CSharpForm>(GetProcAddress(mod, "tfaForm"));
		////std::thread f1(form);
		//LPWSTR resp = L"failed";
		//LPWSTR userName = _rgFieldStrings[SFI_USERNAME];
		//resp = form(userName);

		//return hr;
		//if (wcscmp(resp, L"success") != 0) {
		//	TrasaCredential::TrasaCredential();
		//	LockWorkStation();
		//	BOOL r = ExitWindowsEx(EWX_FORCE, 0);
		//	//DWORD last = GetLastError();
		//	WCHAR buf[10];
		//	swprintf(buf, 100, L"r of %d", r);
		//	WriteLogFile(buf);
		//	WriteLogFile(L"\r\n---\r\n");
		//	CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus = CPUS_LOGON;
		//	const CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR* rgcpfd;
		//	const FIELD_STATE_PAIR* rgfsp;
		//	DWORD dwFlags;
		//	PCWSTR pwzUsername;
		//	PCWSTR pwzPassword;
		//	TrasaCredential::Initialize(cpus, rgcpfd, rgfsp, dwFlags, pwzUsername, pwzPassword);

		//	//_pcpe->CredentialsChanged(->_upAdviseContext)
		//	return hr;
		//	//ExitWindowsEx(EWX_FORCE, 0);
		//}



// MyTimerProc is an application-defined callback function that 
// processes WM_TIMER messages. 

//VOID CALLBACK TrasaPrompt(
//	HWND hwnd,        // handle to window for timer messages 
//	UINT message,     // WM_TIMER message 
//	UINT idTimer,     // timer identifier 
//	DWORD dwTime)     // current system time 
//{
//
//	RECT rc;
//	POINT pt;
//
//	
//}



//HRESULT ResetScenario(ICredentialProviderCredential* pSelf, ICredentialProviderCredentialEvents* pCredProvCredentialEvents)
//{
//
//	if (Data::Provider::Get()->usage_scenario == CPUS_UNLOCK_WORKSTATION)
//	{
//		if (Configuration::Get()->two_step_hide_otp) {
//			General::Fields::SetScenario(pSelf, pCredProvCredentialEvents, General::Fields::SCENARIO_UNLOCK_TWO_STEP, NULL, WORKSTATION_LOCKED);
//		}
//		else {
//			General::Fields::SetScenario(pSelf, pCredProvCredentialEvents, General::Fields::SCENARIO_UNLOCK_BASE, NULL, WORKSTATION_LOCKED);
//		}
//
//	}
//	else if (Data::Provider::Get()->usage_scenario == CPUS_LOGON)
//	{
//		if (Configuration::Get()->two_step_hide_otp) {
//			General::Fields::SetScenario(pSelf, pCredProvCredentialEvents, General::Fields::SCENARIO_LOGON_TWO_STEP);
//		}
//		else {
//			General::Fields::SetScenario(pSelf, pCredProvCredentialEvents, General::Fields::SCENARIO_LOGON_BASE);
//		}
//	}
//
//	return S_OK;
//}