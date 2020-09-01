
#if ( _MSC_VER >= 800 )
#pragma warning ( 3 : 4100 ) // enable "Unreferenced formal parameter"
#pragma warning ( 3 : 4219 ) // enable "trailing ',' used for variable argument list"
#endif

#pragma warning(disable : 4146)

#ifndef WIN32_NO_STATUS
#include <ntstatus.h>
#define WIN32_NO_STATUS
#endif
#include <unknwn.h>

#include "trasautils.h"

//#include "curl/curl.h"

#pragma comment(lib, "winhttp.lib")

using namespace std;

//for json
#include <nlohmann\json.hpp>
using json = nlohmann::json;


//WriteLogFile(ptr);
#define LOGFILE L"C:\\Program Files\\TrasaWin\\TrasaLogs.txt"

BOOL WriteLogFile(LPWSTR Stringval)
{
	HANDLE hFile;
	DWORD dwBytesWritten;

	LPCSTR filer = "C:\\Program Files\\TrasaWin\\TrasaLogs.txt";
	hFile = CreateFile(filer,
		GENERIC_WRITE,
		0,
		NULL,
		OPEN_ALWAYS,
		FILE_ATTRIBUTE_NORMAL | FILE_FLAG_SEQUENTIAL_SCAN,
		NULL);
	if (hFile == INVALID_HANDLE_VALUE)
		return FALSE;

	SetFilePointer(hFile, 0, NULL, FILE_END);
	WriteFile(hFile, Stringval, (wcslen(Stringval) * sizeof(WCHAR)), &dwBytesWritten, NULL);
	// WriteFile(hFile, Stringval, (lstrlen(Stringval) * sizeof(WCHAR)), &dwBytesWritten, NULL);
	//WriteFile(hFile, L"exam", 100, &dwBytesWritten, NULL);

	CloseHandle(hFile);

	return TRUE;
}

// Reading File from local config store
#define BUFFERSIZE 200
DWORD g_BytesTransferred = 0;

VOID CALLBACK FileIOCompletionRoutine(
	__in  DWORD dwErrorCode,
	__in  DWORD dwNumberOfBytesTransfered,
	__in  LPOVERLAPPED lpOverlapped
);

VOID CALLBACK FileIOCompletionRoutine(
	__in  DWORD dwErrorCode,
	__in  DWORD dwNumberOfBytesTransfered,
	__in  LPOVERLAPPED lpOverlapped)
{
	printf("Error code:\t%x\n", dwErrorCode);
	printf("Number of bytes:\t%x\n", dwNumberOfBytesTransfered);
	g_BytesTransferred = dwNumberOfBytesTransfered;
}

#define CONFIGFILE L"C:\\Program Files\\TrasaCP\\config.dat"
// readFileFunc reads config file and returns json value as char array.
// TO-DO below function currently snaps plain text values as config file since we are not storing
// config values encrypted. When encrypted, this function should be able to deecrypt and return plain text values.
std::string readFileFunc() {

	// create file dECLERATIONS
	std::string filePath = "C:\\Program Files\\TrasaCP\\config.dat";

	//std::wstring sstemp = std::wstring(filePath.begin(), filePath.end());
	LPCSTR fileName = filePath.c_str();


	HANDLE hFile;
	DWORD  dwBytesRead = 0;
	char   ReadBuffer[BUFFERSIZE] = { 0 };
	OVERLAPPED ol = { 0 };

	hFile = CreateFile(fileName,               // file to open
		GENERIC_READ,          // open for reading
		FILE_SHARE_READ,       // share for reading
		NULL,                  // default security
		OPEN_EXISTING,         // existing file only
		FILE_ATTRIBUTE_NORMAL | FILE_FLAG_OVERLAPPED, // normal file
		NULL);                 // no attr. template

	if (hFile == INVALID_HANDLE_VALUE)
	{
		//DisplayError(TEXT("CreateFile"));
		printf("Terminal failure: unable to open file \"%s\" for read.\n", fileName);
	}

	// Read one character less than the buffer size to save room for
	// the terminating NULL character. 

	if (FALSE == ReadFileEx(hFile, ReadBuffer, BUFFERSIZE - 1, &ol, FileIOCompletionRoutine))
	{
		//DisplayError(TEXT("ReadFile"));
		printf("Terminal failure: Unable to read from file.\n GetLastError=%08x\n", GetLastError());
		CloseHandle(hFile);

	}

	// Read one character less than the buffer size to save room for
	// the terminating NULL character. 
	SleepEx(5000, TRUE);
	dwBytesRead = g_BytesTransferred;
	// This is the section of code that assumes the file is ANSI text. 
	// Modify this block for other data types if needed.

	if (dwBytesRead > 0 && dwBytesRead <= BUFFERSIZE - 1)
	{
		ReadBuffer[dwBytesRead] = '\0'; // NULL character

		printf("Data read from %s (%d bytes): \n", fileName, dwBytesRead);
		//printf("buffer: %s\n", ReadBuffer);
	}
	else if (dwBytesRead == 0)
	{
		printf("No data read from file %s\n", fileName);
	}
	else
	{
		printf("\n ** Unexpected value for dwBytesRead ** \n");
	}

	// It is always good practice to close the open file handles even though
	// the app will exit here and clean up open handles anyway.

	CloseHandle(hFile);

	printf("outside close handle \n");
	//	std::cout << ReadBuffer << std::endl;


	return ReadBuffer;
}

std::string sendRequest(std::string user, std::string totp) {

	//std::string to store respvalue;
	DWORD dwSize = 0;
	LPSTR respValue;
	respValue = new char[dwSize + 1];

	// Get config values from local file
	std::string jsonValue = readFileFunc().c_str();
	auto configVals = json::parse(jsonValue);

	// IMP!! code below to CONVERT string to PWSTR
	//char username[128];
	//wcstombs(username, user, 128);
	//char password[128];
	//wcstombs(password, pass, 128);

	//converting unicode_string to string
	//std::wstring username(user.Buffer, user.Length);

	// get workstation name
	TCHAR  infoBuf[INFO_BUFFER_SIZE];
	DWORD  bufCharCount = INFO_BUFFER_SIZE;

	// Get and display the name of the computer. 
	bufCharCount = INFO_BUFFER_SIZE;
	if (!GetComputerName(infoBuf, &bufCharCount))
		_tprintf(TEXT("GetComputerName"));
	//_tprintf(TEXT("\nComputer name:      %s"), infoBuf);

	// create http request structure (post data)
	json reqConfig;
	reqConfig["appID"] = configVals.value("appID", "name");
	reqConfig["appSecret"] = configVals.value("appSecret", "name");
	reqConfig["user"] = user;
	reqConfig["totp"] = totp;
	reqConfig["workstation"] = infoBuf;
	//reqConfig["password"] = password;

	std::string _data = reqConfig.dump();


	std::string mockstatus = "false";

	// start http request.
	//DWORD dwSize = 0;
	DWORD dwDownloaded = 0;
	LPSTR pszOutBuffer;
	BOOL  bResults = FALSE;
	HINTERNET  hSession = NULL,
		hConnect = NULL,
		hRequest = NULL;

	//std::string _data = "tester=tester ";

	// Use WinHttpOpen to obtain a session handle.
	hSession = WinHttpOpen(L"TrasaAuth Client/1.0",
		WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
		WINHTTP_NO_PROXY_NAME,
		WINHTTP_NO_PROXY_BYPASS, 0);

	// we get hostname (endpoint) from config file.
	std::string hostname = configVals["trasaHost"];
	// coverting std::string to LPCWSTR
	std::wstring stemp = std::wstring(hostname.begin(), hostname.end());
	LPCWSTR endpoint = stemp.c_str();

	// Specify an HTTP server.
	if (hSession)
	{


		//WinHttpSetTimeouts(hSession, 30000, 30000, 30000, 30000);

		hConnect = WinHttpConnect(hSession, endpoint, 3339, 0); //3331, 0); //INTERNET_DEFAULT_HTTPS_PORT, 0);

																					   // Create an HTTP request handle.
		if (hConnect)
		{
			hRequest = WinHttpOpenRequest(hConnect, L"POST", L"/api/v1/remote/auth/passuserstotpval",
				NULL, WINHTTP_NO_REFERER,
				WINHTTP_DEFAULT_ACCEPT_TYPES,
				NULL); //WINHTTP_FLAG_SECURE);
		}
		else {
			return mockstatus;
		}


		// Send a request.
		if (hRequest)
		{
			bResults = WinHttpSendRequest(hRequest,
				WINHTTP_NO_ADDITIONAL_HEADERS, 0,
				(LPVOID)_data.c_str(), _data.length(), _data.length(), 0);
		}
		else {
			return mockstatus;
		}



		// End the request.
		if (bResults)
		{
			bResults = WinHttpReceiveResponse(hRequest, NULL);
		}
		else {
			return mockstatus;
		}


		// Keep checking for data until there is nothing left.
		// Keep checking for data until there is nothing left.
		if (bResults)
		{
			do
			{
				// Check for available data.
				dwSize = 0;
				if (!WinHttpQueryDataAvailable(hRequest, &dwSize))
					return mockstatus;

				// Allocate space for the buffer.
				pszOutBuffer = new char[dwSize + 1];
				if (!pszOutBuffer)
				{
					//printf("Out of memory\n");
					dwSize = 0;
					return mockstatus;
				}
				else
				{
					// Read the data.
					// ZeroMemory(pszOutBuffer, dwSize + 1);
					if (!ZeroMemory(pszOutBuffer, dwSize + 1))
					{
						return mockstatus;
					}
					//printf("first data was: %s\n", pszOutBuffer);
					if (!WinHttpReadData(hRequest, (LPVOID)pszOutBuffer, dwSize, &dwDownloaded))
					{
						return mockstatus;
					}

					//else
					//printf(" response data is %s\n", pszOutBuffer);

					//const char *s = pszOutBuffer;
					//std::string str(s);
					//val.append(str);
					respValue = pszOutBuffer;

					//printf("first val is: %s\n", respValue);
					break;
					//lstrcpy (val, pszOutBuffer);

					//std::string s = CT2A(pszOutBuffer);
					// Free the memory allocated to the buffer.
					//delete[] pszOutBuffer;
				}
			} while (dwSize > 0);
		}
		else
		{
			return mockstatus;
		}

		// Report any errors.
		if (!bResults)
		{
			return mockstatus;
		}

		// Close any open handles.
		if (hRequest) WinHttpCloseHandle(hRequest);
		if (hConnect) WinHttpCloseHandle(hConnect);
		if (hSession) WinHttpCloseHandle(hSession);

		if (strlen(respValue) < 1)
		{
			//std::string mockstatus = "failed";
			return mockstatus;
		}

		json resp = json::parse(respValue);


		// a simple struct to model a person

		std::string status = resp.value("status", "name");
		std::string reason = resp.value("reason", "name");

		return status;

		//////////////////////////////////////

		if (!WinHttpCloseHandle(hSession))
		{
			return mockstatus;
		}
	}
	else {
		return mockstatus;
	}


	return mockstatus;
}


