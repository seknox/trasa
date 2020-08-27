//
// THIS CODE AND INFORMATION IS PROVIDED "AS IS" WITHOUT WARRANTY OF
// ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A
// PARTICULAR PURPOSE.
//
// Copyright (c) 2006 Microsoft Corporation. All rights reserved.
//
// Helper functions for copying parameters and packaging the buffer
// for GetSerialization.

#pragma once
#include "common.h"
#include <windows.h>
#include <strsafe.h>

#pragma warning(push)
#pragma warning(disable : 4995)
#include <shlwapi.h>
#pragma warning(pop)



//makes a copy of a field descriptor using CoTaskMemAlloc
HRESULT FieldDescriptorCoAllocCopy(
    const CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR& rcpfd,
    CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR** ppcpfd
    );

//makes a copy of a field descriptor on the normal heap
HRESULT FieldDescriptorCopy(
    const CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR& rcpfd,
    CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR* pcpfd
    );

//creates a UNICODE_STRING from a NULL-terminated string
HRESULT UnicodeStringInitWithString(
    PWSTR pwz, 
    UNICODE_STRING* pus
    );

//initializes a KERB_INTERACTIVE_UNLOCK_LOGON with weak references to the provided credentials
HRESULT KerbInteractiveUnlockLogonInit(
    PWSTR pwzDomain,
    PWSTR pwzUsername,
    PWSTR pwzPassword,
    CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus,
    KERB_INTERACTIVE_UNLOCK_LOGON* pkiul
    );

//packages the credentials into the buffer that the system expects
HRESULT KerbInteractiveUnlockLogonPack(
    const KERB_INTERACTIVE_UNLOCK_LOGON& rkiulIn,
    BYTE** prgb,
    DWORD* pcb
    );

//get the authentication package that will be used for our logon attempt
HRESULT RetrieveNegotiateAuthPackage(
    ULONG * pulAuthPackage
    );


//encrypt a password (if necessary) and copy it; if not, just copy it
HRESULT ProtectIfNecessaryAndCopyPassword(
    PWSTR pwzPassword,
    CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus,
    PWSTR* ppwzProtectedPassword);

HRESULT KerbInteractiveUnlockLogonRepackNative(
    __in_bcount(cbWow) BYTE* rgbWow,
    DWORD cbWow,
    __out_bcount(*pcbNative) BYTE** prgbNative,
    __out DWORD* pcbNative);

void KerbInteractiveUnlockLogonUnpackInPlace(
    __inout_bcount(cb) KERB_INTERACTIVE_UNLOCK_LOGON* pkiul,
    DWORD cb
    );

HRESULT DomainUsernameStringAlloc(
    const PWSTR pwszDomain,
    const PWSTR pwszUsername,
    __deref_out PWSTR* ppwszDomainUsername);
