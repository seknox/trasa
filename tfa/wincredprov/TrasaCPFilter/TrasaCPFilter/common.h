 
// This file contains some global variables that describe what our
// sample tile looks like.  For example, it defines what fields a tile has 
// and which fields show in which states of LogonUI.

#pragma once
#include <credentialprovider.h>
#include <ntsecapi.h>
#define SECURITY_WIN32
#include <security.h>
#include <intsafe.h>

#define MAX_ULONG  ((ULONG)(-1))

// The indexes of each of the fields in our credential provider's tiles.
enum SAMPLE_FIELD_ID 
{
    SFI_TILEIMAGE       = 0,
	SFI_LARGE_TEXT		= 1,
    SFI_USERNAME       =  2,
    SFI_PASSWORD        = 3,
    SFI_TOTP            = 4,
    SFI_SUBMIT_BUTTON   = 5, 
    SFI_NUM_FIELDS      = 6,  // Note: if new fields are added, keep NUM_FIELDS last.  This is used as a count of the number of fields
};

// The first value indicates when the tile is displayed (selected, not selected)
// the second indicates things like whether the field is enabled, whether it has key focus, etc.
struct FIELD_STATE_PAIR
{
    CREDENTIAL_PROVIDER_FIELD_STATE cpfs;
    CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE cpfis;
};

// These two arrays are seperate because a credential provider might
// want to set up a credential with various combinations of field state pairs 
// and field descriptors.

// The field state value indicates whether the field is displayed
// in the selected tile, the deselected tile, or both.
// The Field interactive state indicates when 
static const FIELD_STATE_PAIR s_rgFieldStatePairs[] = 
{
    { CPFS_DISPLAY_IN_BOTH, CPFIS_NONE },                   // SFI_TILEIMAGE
	{ CPFS_DISPLAY_IN_BOTH, CPFIS_NONE },                   // SFI_LARGE_TEXT
    { CPFS_DISPLAY_IN_BOTH, CPFIS_FOCUSED },                   // SFI_USERNAME
    { CPFS_DISPLAY_IN_SELECTED_TILE, CPFIS_NONE },       // SFI_PASSWORD
    { CPFS_DISPLAY_IN_SELECTED_TILE, CPFIS_NONE },       // SFI_TOTP
    { CPFS_DISPLAY_IN_SELECTED_TILE, CPFIS_NONE    },       // SFI_SUBMIT_BUTTON   
};

// Field descriptors for unlock and logon.
// The first field is the index of the field.
// The second is the type of the field.
// The third is the name of the field, NOT the value which will appear in the field.
static const CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR s_rgCredProvFieldDescriptors[] =
{
    { SFI_TILEIMAGE, CPFT_TILE_IMAGE, L"Image" },
	{ SFI_LARGE_TEXT, CPFT_LARGE_TEXT, L"TRASA-PROVIDER" },
    { SFI_USERNAME, CPFT_EDIT_TEXT, L"Username" },
    { SFI_PASSWORD, CPFT_PASSWORD_TEXT, L"Password" },
    { SFI_TOTP, CPFT_EDIT_TEXT, L"TOTP" },
    { SFI_SUBMIT_BUTTON, CPFT_SUBMIT_BUTTON, L"Submit" },
};