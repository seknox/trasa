/* * * * * * * * * * * * * * * * * * * * *
**
** Copyright 2012 Dominik Pretzsch
**
**    Licensed under the Apache License, Version 2.0 (the "License");
**    you may not use this file except in compliance with the License.
**    You may obtain a copy of the License at
**
**        http://www.apache.org/licenses/LICENSE-2.0
**
**    Unless required by applicable law or agreed to in writing, software
**    distributed under the License is distributed on an "AS IS" BASIS,
**    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
**    See the License for the specific language governing permissions and
**    limitations under the License.
**
** * * * * * * * * * * * * * * * * * * */

#ifndef WIN32_NO_STATUS
#include <ntstatus.h>
#define WIN32_NO_STATUS
#endif
#include <unknwn.h>
#include "TrasaCPFilter.h"
#include "guid.h"


#define _UNICODE
#define UNICODE
#include <stdio.h>
#include <winnt.h>
#include <string>
#include <iostream>

#include <vector>
#include <strsafe.h>
#include <tchar.h>

#include "trasautils.h"

WCHAR buf[256];




// Boilerplate code to create our provider.
HRESULT TrasaCPFilterCreateInstance(__in REFIID riid, __deref_out void** ppv)
{
	HRESULT hr;

	TrasaCPFilter* pProvider = new TrasaCPFilter();

	if (pProvider)
	{
		hr = pProvider->QueryInterface(riid, ppv);
		pProvider->Release();
	}
	else
	{
		hr = E_OUTOFMEMORY;
	}

	return hr;
}


HRESULT TrasaCPFilter::Filter(CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus, DWORD dwFlags, GUID* rgclsidProviders, BOOL* rgbAllow, DWORD cProviders)
{


	//WriteLogFile(TEXT("No identity\r\n"));

	switch (cpus)
	{
	case CPUS_LOGON:
		for (DWORD i = 0; i < cProviders; i++)
		{
			if (i < dwFlags)
			{
			}
			//if (IsEqualGUID(rgclsidProviders[i], CLSID_PasswordCredentialProvider))
			// Only allow OTP CPs
			if (IsEqualGUID(rgclsidProviders[i], CLSID_TrasaProvider)) {
				rgbAllow[i] = TRUE;
			}
			else {
				rgbAllow[i] = FALSE;
			}
		}
		return S_OK;
		break;
	case CPUS_UNLOCK_WORKSTATION:
		for (DWORD i = 0; i < cProviders; i++)
		{
			if (i < dwFlags)
			{
			}
			//if (IsEqualGUID(rgclsidProviders[i], CLSID_PasswordCredentialProvider))
			// Only allow OTP CPs
			if (IsEqualGUID(rgclsidProviders[i], CLSID_TrasaProvider)) {
				rgbAllow[i] = TRUE;
			}
			else {
				rgbAllow[i] = FALSE;
			}
		}
		return S_OK;
		break;
	case CPUS_CREDUI:
		for (DWORD i = 0; i < cProviders; i++)
		{
			if (i < dwFlags)
			{
			}
			//if (IsEqualGUID(rgclsidProviders[i], CLSID_PasswordCredentialProvider))
			// Only allow OTP CPs
			if (IsEqualGUID(rgclsidProviders[i], CLSID_TrasaProvider) ) {
				rgbAllow[i] = TRUE;
			}
			else {
				rgbAllow[i] = FALSE;
			}
		}
		return S_OK;
		break;
	
	case CPUS_CHANGE_PASSWORD:
		return E_NOTIMPL;
	default:
		return E_INVALIDARG;
	}
}

TrasaCPFilter::TrasaCPFilter() :
	_cRef(1)
{
	DllAddRef();
}

TrasaCPFilter::~TrasaCPFilter()
{
	DllRelease();
}

HRESULT TrasaCPFilter::UpdateRemoteCredential(const CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpsIn, CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcsOut)
{

	HRESULT hr = E_NOTIMPL;

	if (!pcpsIn) // no point continuing has there are no credentials
		return E_NOTIMPL;

	pcpcsOut->ulAuthenticationPackage = pcpsIn->ulAuthenticationPackage;
	pcpcsOut->cbSerialization = pcpsIn->cbSerialization;
	pcpcsOut->rgbSerialization = pcpsIn->rgbSerialization;
	pcpcsOut->clsidCredentialProvider = CLSID_TrasaProvider;

	// if I need to copy the buffer contents I will use:
	if (pcpcsOut->cbSerialization > 0 && (pcpcsOut->rgbSerialization = (BYTE*)CoTaskMemAlloc(pcpsIn->cbSerialization)) != NULL)
	{
		CopyMemory(pcpcsOut->rgbSerialization, pcpsIn->rgbSerialization, pcpsIn->cbSerialization);
		WriteLogFile(L"\r\n returning S_OK \r\n");
		hr = S_OK;
	}



	return hr;
}


