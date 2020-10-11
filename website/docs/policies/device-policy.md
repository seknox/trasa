---
id: device-policy
title: Device Policy
sidebar_label: Device Policy
---
Device Policy supports enables controlling access based on the security hygiene of user devices. 

All of these device policies are blocking i.e. user access is denied if any one of them matches. 

1. Untrusted Devices:
    + Block all the devices which are not manually marked as "trusted" by admin. 

2. Autologin enabled:
    + Block if the user can log in without a password.

3. Idle screen lock disabled:
    + Block if screen-lock is disabled.

4. Remote login enabled (Workstation):
    + Block if remote access (RDP, SSH) is enabled in the device

5. Jailbroken device (Mobile device):
    + Block if the mobile device is jailbroken or rooted.

6. Debugging enabled (Mobile device):
    + Block if the mobile device has debugging enabled.

7. Emulated device (Mobile device):
    + Block if the mobile device is not a real device.

8. Disk not encrypted (Workstation):
    + Block if disk encryption is not set.

9. Firewall disabled (Workstation):
    + Block if the firewall is disabled.

10. Antivirus disabled (Windows only)
    + Block if antivirus is disabled.
