# trasaWrkstnAgent

trasaWrkstnAgent collects device hygiene of workstations. trasa browser extensions and cli can use this agent to retreive device hygiene of users.

## Building

Debian: `$ go build`

## Usage

### Debian

- The binary build `trasaWrkstnAgent` should be copied in `/usr/local/bin/`

- The file `trasaWrkstnAgent.nix.json` should be copied to `/usr/lib/mozilla/native-messaging-hosts/trasaWrkstnAgent.json`. Create directory if not exist !

- Create log file in as `/var/log/trasaWrkstnAgent.log`. This log file must be chown with current user permission to be able to writtin by `trasaWrkstnAgent` binary.

### Windows

- The binary build `trasaWrkstnAgent.exe` should be copied in `C:\Program Files\trasaWrkstnAgent`

- Dependency of bitlocker (contents of device/bitlocer-status) should be copied to `C:\Program Files\trasaWrkstnAgent`

- `trasaWrkstnAgent.windows.json` should be copied to `C:\Program Files\Trasa\trasaWrkstnAgent\trasaWrkstnAgent.json`

- Create a registry key `LOCAL_MACHINE\SOFTWARE\Mozilla\NativeMessagingHosts\trasaWrkstnAgent` with a default value `C:\Program Files\Trasa\trasaWrkstnAgent\trasaWrkstnAgent.json` (path of trasaWrkstnAgent.json).

- Create log file/folder in as `C:\Program Files\trasaWrkstnAgent\trasaWrkstnAgent.log`. This logs file must be writable with current user.
