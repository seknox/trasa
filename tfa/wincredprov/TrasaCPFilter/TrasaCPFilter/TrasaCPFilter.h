#pragma once
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

#pragma once

#include "helpers.h"
#include "common.h"
#include "dll.h"
#include "resource.h"

class TrasaCPFilter : public ICredentialProviderFilter
{
public:
	//This section contains some COM boilerplate code 

	// IUnknown 
	STDMETHOD_(ULONG, AddRef)()
	{
		return _cRef++;
	}

	STDMETHOD_(ULONG, Release)()
	{
		LONG cRef = _cRef--;
		if (!cRef)
		{
			delete this;
		}
		return cRef;
	}

	STDMETHOD(QueryInterface)(REFIID riid, void** ppv)
	{
		HRESULT hr;
		if (IID_IUnknown == riid || IID_ICredentialProviderFilter == riid)
		{
			*ppv = this;
			reinterpret_cast<IUnknown*>(*ppv)->AddRef();
			hr = S_OK;
		}
		else
		{
			*ppv = NULL;
			hr = E_NOINTERFACE;
		}
		return hr;
	}
	//#pragma warning(disable:4100)

public:
	//Implementation of ICredentialProviderFilter 
	IFACEMETHODIMP Filter(CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus,
		DWORD dwFlags,
		GUID* rgclsidProviders,
		BOOL* rgbAllow,
		DWORD cProviders);

	IFACEMETHODIMP UpdateRemoteCredential(const CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcsIn,
		CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcsOut);

	friend HRESULT TrasaCPFilterCreateInstance(__in REFIID riid, __deref_out void** ppv);

protected:
	TrasaCPFilter();
	__override ~TrasaCPFilter();

private:
	LONG _cRef;
};