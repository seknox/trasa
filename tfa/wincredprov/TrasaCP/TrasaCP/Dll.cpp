//
// THIS CODE AND INFORMATION IS PROVIDED "AS IS" WITHOUT WARRANTY OF
// ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A
// PARTICULAR PURPOSE.
//
// Copyright (c) 2006 Microsoft Corporation. All rights reserved.
//
// Standard dll required functions and class factory implementation.

#include <windows.h>
#include <unknwn.h>
#include "Dll.h"
#include "guid.h"

static LONG g_cRef = 0;   // global dll reference count

// IClassFactory ///////////////////////////////////////////////////////////////////////

extern HRESULT TrasaProvider_CreateInstance(REFIID riid, void** ppv);

HINSTANCE g_hinst = NULL;   // global dll hinstance


class CClassFactory : public IClassFactory
{
  public:
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

    STDMETHOD (QueryInterface)(REFIID riid, void** ppv) 
    {
        HRESULT hr;
        if (ppv != NULL)
        {
            if (IID_IClassFactory == riid || IID_IUnknown == riid)
            {
                *ppv = static_cast<IUnknown*>(this);
                reinterpret_cast<IUnknown*>(*ppv)->AddRef();
                hr = S_OK;
            }
            else
            {
                *ppv = NULL;
                hr = E_NOINTERFACE;
            }
        }
        else
        {
            hr = E_INVALIDARG;
        }
        return hr;
    }

    // IClassFactory
    STDMETHOD (CreateInstance)(IUnknown* pUnkOuter, REFIID riid, void** ppv)
    {
        HRESULT hr;
        if (!pUnkOuter)
        {
            hr = TrasaProvider_CreateInstance(riid, ppv);
        }
        else
        {
            hr = CLASS_E_NOAGGREGATION;
        }
        return hr;
    }

    STDMETHOD (LockServer)(BOOL bLock)
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
     CClassFactory() : _cRef(1) {}
    ~CClassFactory(){}

  private:
    LONG _cRef;

    friend HRESULT CClassFactory_CreateInstance(REFCLSID rclsid, REFIID riid, void** ppv);
};

HRESULT CClassFactory_CreateInstance(REFCLSID rclsid, REFIID riid, void** ppv)
{
    HRESULT hr;
    if (CLSID_TrasaProvider == rclsid)
    {
        CClassFactory* pcf = new CClassFactory;
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

// DLL Functions ///////////////////////////////////////////////////////////////////////

BOOL WINAPI DllMain(
    HINSTANCE hinstDll,
    DWORD dwReason,
    LPVOID pReserved
    )
{
    UNREFERENCED_PARAMETER(pReserved);

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

void DllAddRef()
{
    InterlockedIncrement(&g_cRef);
}

void DllRelease()
{
    InterlockedDecrement(&g_cRef);
}

// DLL entry point.
STDAPI DllCanUnloadNow()
{
    HRESULT hr;

    if (g_cRef > 0)
    {
        hr = S_FALSE;   // cocreated objects still exist, don't unload
    }
    else
    {
        hr = S_OK;      // refcount is zero, ok to unload
    }

    return hr;
}

// DLL entry point.
STDAPI DllGetClassObject(REFCLSID rclsid, REFIID riid, void** ppv)
{
    return CClassFactory_CreateInstance(rclsid, riid, ppv);
}

