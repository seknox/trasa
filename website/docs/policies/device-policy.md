---
id: device-policy
title: Device Policy
sidebar_label: Device Policy
---
Device Policy supports enables controlling access based on the security hygiene of user devices. 

All of these device policies are blocking i.e. user access is denied if any one of them matches. 

### Untrusted Devices:
Block all the devices which are not manually marked as "trusted" by admin. 

### Autologin enabled:
Access is blocked if the user device can be logged in without a password.

### Idle screen lock disabled:
Access is blocked if screen-lock is disabled.

### Remote login enabled (Workstation):
Access is blocked if remote access (RDP, SSH) is enabled in the user device

### Jailbroken device (Mobile device):
Access is blocked if the mobile device is jailbroken or rooted.

### Debugging enabled (Mobile device):
Access is blocked if the mobile device has debugging enabled.

### Emulated device (Mobile device):
Access is blocked if the mobile device is not a real device.

### Disk not encrypted (Workstation):
Access is blocked if disk encryption is not set.

### Firewall disabled (Workstation):
Access is blocked if the firewall is disabled.

### Antivirus disabled (Windows only)
Access is blocked if antivirus is disabled.
