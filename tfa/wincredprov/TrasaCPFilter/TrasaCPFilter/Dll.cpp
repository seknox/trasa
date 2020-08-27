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

#include <windows.h>
#include <unknwn.h>
#include "Dll.h"
#include "helpers.h"

static LONG g_cRef = 0;   // global dll reference count
HINSTANCE g_hinst = NULL; // global dll hinstance

extern HRESULT TrasaCPFilterCreateInstance(__in REFIID riid, __deref_out void** ppv);
EXTERN_C GUID CLSID_CSample;

class CClassFactory : public IClassFactory
{
public:
	CClassFactory() : _cRef(1)
	{
	}

	// IUnknown
	IFACEMETHODIMP QueryInterface(__in REFIID riid, __deref_out void** ppv)
	{
		static const QITAB qit[] =
		{
			QITABENT(CClassFactory, IClassFactory),
			{ 0 },
		};
		return QISearch(this, qit, riid, ppv);
	}

	IFACEMETHODIMP_(ULONG) AddRef()
	{
		return InterlockedIncrement(&_cRef);
	}

	IFACEMETHODIMP_(ULONG) Release()
	{
		LONG cRef = InterlockedDecrement(&_cRef);
		if (!cRef)
			delete this;
		return cRef;
	}

	// IClassFactory
	IFACEMETHODIMP CreateInstance(__in IUnknown* pUnkOuter, __in REFIID riid, __deref_out void** ppv)
	{
		HRESULT hr;
		if (!pUnkOuter)
		{
			hr = TrasaCPFilterCreateInstance(riid, ppv);
		}
		else
		{
			*ppv = NULL;
			hr = CLASS_E_NOAGGREGATION;
		}
		return hr;
	}

	IFACEMETHODIMP LockServer(__in BOOL bLock)
	{
		if (bLock)
		{
			DllAddRef();
		}
		else
		{
			DllRelease();
		}
		return S_OK;
	}

private:
	~CClassFactory()
	{
	}
	long _cRef;
};

HRESULT CClassFactory_CreateInstance(__in REFCLSID rclsid, __in REFIID riid, __deref_out void** ppv)
{
	*ppv = NULL;

	HRESULT hr;

	if (CLSID_CSample == rclsid)
	{
		CClassFactory* pcf = new CClassFactory();
		if (pcf)
		{
			hr = pcf->QueryInterface(riid, ppv);
			pcf->Release();
		}
		else
		{
			hr = E_OUTOFMEMORY;
		}
	}
	else
	{
		hr = CLASS_E_CLASSNOTAVAILABLE;
	}
	return hr;
}

void DllAddRef()
{
	InterlockedIncrement(&g_cRef);
}

void DllRelease()
{
	InterlockedDecrement(&g_cRef);
}

STDAPI DllCanUnloadNow()
{
	return (g_cRef > 0) ? S_FALSE : S_OK;
}

STDAPI DllGetClassObject(__in REFCLSID rclsid, __in REFIID riid, __deref_out void** ppv)
{
	return CClassFactory_CreateInstance(rclsid, riid, ppv);
}

STDAPI_(BOOL) DllMain(__in HINSTANCE hinstDll, __in DWORD dwReason, __in void*)
{
	switch (dwReason)
	{
	case DLL_PROCESS_ATTACH:
		DisableThreadLibraryCalls(hinstDll);
		break;
	case DLL_PROCESS_DETACH:
	case DLL_THREAD_ATTACH:
	case DLL_THREAD_DETACH:
		break;
	}

	g_hinst = hinstDll;
	return TRUE;
}

