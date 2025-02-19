# Changelog

All notable changes to this project will be documented in this file.  

## Version 3.10.15 - 2023-03-29

[FIX] (Windows) Prevent installer from inadvertently overwriting PATH environment variable in rare cases  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.10.15.exe)  
SHA256: a146106203baf6d5122d06a563d30d85517d8ff8d792c3b0cd4136b115858f6e  

## Version 3.10.15 - 2023-02-28

[FIX] (Linux) Fixed application shortcut registration issue on some distributions  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.10.15_amd64.deb)  
SHA256: 75c3f1867f9841b127d6936038bf10b6edbbe0bb27be7c2ee91b32b2def95bc3  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.10.15-1.x86_64.rpm)  
SHA256: ea9b5719c21bd0ad2dc662acc15fa82f42beb357c31cc3e030f45b1d27467687  

## Version 3.10.14 - 2023-02-24

[NEW] Ability to open Firewall and AntiTracker settings with one click  
[NEW] Option to prevent usage of the same provider in Multi-Hop chain  
[NEW] (Linux) IVPN can be installed on Fedora Silverblue (using rpm-ostree)  
[IMPROVED] The favorite servers list is common for all VPN protocols  
[IMPROVED] Changing protocol type or MultiHop does not require disconnecting the current VPN connection  
[IMPROVED] Various UI fixes/improvements  
[IMPROVED] (Windows) Improved boot-time firewall rules to prevent potential leaks on system boot  
[IMPROVED] (Linux) Installed files locations are corrected to fit the Filesystem Hierarchy Standard  
[IMPROVED] (Linux) IVPN Firewall now also controls the FORWARD chain (no leaks anymore when using IVPN on Qubes OS as "ProxyVM")  
[FIX] The Antitracker toggle state in the UI is consistent with actions from CLI  
[FIX] (macOS) The wrong DNS configuration may stay after VPN is disconnected in some corner cases  
[FIX] (Linux) "Force management of DNS using resolve.conf" does not require a reboot anymore  
[FIX] (Linux) Bad DNS configuration after reinitiation of the main network interface  
[FIX] (Linux) Split Tunneling stops working after reinitiation of the main network interface  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.10.14.exe)  
SHA256: 720e31271aee95ce68907c238b88271533bc67239d5bc08871809952f8224619  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.10.14.dmg)  
SHA256: 7736e92d6c04a41040664cf37b6b2f48ada0f46786dcd0fe8316994390b04b1e  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.10.14-arm64.dmg)  
SHA256:  4422d251e3b5b5c97bc5b5a7b790f8e9f595a98032d3f6165c82ce18c9aefae2  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.10.14_amd64.deb)  
SHA256: 1b0ba710ca8f1b5f369dc0d91475dc9056127bdef4fdd9f0076c3ce0ea442764  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.10.14-1.x86_64.rpm)  
SHA256: d54ba666a1297e167b6f96b02079d5d59c14468cd7320c0cd81dd958b411ad14  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.10.14_amd64.deb)  
SHA256: 37585671d941b946ed65b394c88e82d63cd7e0c28016ea1063db0bb5e17f95bf  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.10.14-1.x86_64.rpm)  
SHA256: d6c6e16a16308e0fa49d30c14817fd1cfa4d06acfe6cc8ab5d769ce2210e86f1  

## Version 3.10.0 - 2022-12-05

[NEW] (Windows/Linux) UI: light/dark options for system tray icon  
[NEW] Ability to manage ‘Autoconnect on launch’ settings from CLI  
[NEW] Ability to manage ‘WiFi control’ settings from CLI  
[NEW] (Linux) ‘WiFi control’ option: ‘Autoconnect on joining WiFi networks without encryption’  
[IMPROVED] Diagnostic logs extended with additional data  
[FIX] (Linux) Firewall fails on some versions of Raspberry Pi OS  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.10.0.exe)  
SHA256: 13022a189ddc8608e052a60ff800cfd21ed854bdcefc1dde212a810cedadda7b  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.10.0.dmg)  
SHA256: d3a99cd45604e800108ecc8138cc106c4cf553743babc06a3cbb01f93fefbb68  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.10.0-arm64.dmg)  
SHA256:  5acf2e20531fb3de0644bdd9fb57a27fcc7fa9628349f7b77ae095614e6b161d  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.10.0_amd64.deb)  
SHA256: cae1e13e04dfc8c721d2adcb9712d99db443701cede02cc9b44ba349c1559cbf  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.10.0-1.x86_64.rpm)  
SHA256: f775da2fa3b8feff97857b23d3db5e9224dd1d7007d95a1e2c73e722ff4d28b1  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.10.0_amd64.deb)  
SHA256: 7acf5f2bf32d0fdf0d0bfdd92b700deb4e20af54d9caab571bef86172ae7a77e  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.10.0-1.x86_64.rpm)  
SHA256: 899aae46263b8dec0e0681e1974c737a887afe9ffed0e5632600e9dfe287d54d  

## Version 3.9.45 - 2022-11-09

[FIX] Ability to send diagnostic logs  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.9.45.exe)  
SHA256: e6e3df093c44c83a87cbf538d39983945f9955fadaec3f57a1a23fdfc8fcd678  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.9.45.dmg)  
SHA256: 6768f6c79563ab0205604ba8b3513d033b4d4b5d456f9100556cecc8cf79cc8f  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.9.45-arm64.dmg)  
SHA256: 3b91774e18b4e288d7794c5342002f7a90015d8cc2c637b1f13303dd3b06c5fd   

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.9.45_amd64.deb)  
SHA256: d9ba7d9ceeed05c2e8913a5c36ad67194319246c8de932c30e0623e50fdb804d  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.9.45-1.x86_64.rpm)  
SHA256: 9330ca5f4dc6e8375410b4005368d910b95aedfe91bc47e030d0f8876e3882c3  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.9.45_amd64.deb)  
SHA256: 114abc02f0c0b62965ac6904962ba1fd9703585cebb8b495e34bdd1de2c124c0  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.9.45-1.x86_64.rpm)  
SHA256: 580e3bca37c35984ce7f96448195fa02a26c4a9073c40fcd537495fa71856c4f  

## Version 3.9.43 - 2022-10-28

[NEW] obfs4 support for OpenVPN connections  
[IMPROVED] UI: Eliminated the delay which sometimes occurred before showing dialogs in app settings  
[IMPROVED] UI: Migrated to latest frameworks  
[IMPROVED] (Linux) Removed dependency from the "which" command  
[FIX] UI: Mismatch of servers when changing entry/servers immediately one after the other  
[FIX] UI: The application did not connect to the last port selected  
[FIX] UI: The obfsproxy settings were disabled after changing the protocol  
[FIX] (macOS) The application was unusable when installed from Homebrew Cask  
[FIX] (macOS) Restoring DNS configuration after changing WiFi networks  
[FIX] (macOS) Automatic updates installer issue  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.9.43.exe)  
SHA256: 5bae4107305c33aed8c6c657965317a0a1ba6bf026d244d5da412bc5cccf98ad  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.9.43.dmg)  
SHA256: 0e167ddc6418338b4f163008ee418eb0a3d5d6c8ff487b065fd7429b3d286b19  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.9.43-arm64.dmg)  
SHA256: ad0d9f1eae0a42759df1f1115f8edd1c9176f8af46516d6d2ada2958409e1a1e

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.9.43_amd64.deb)  
SHA256: 033564a1d2cb4c45c064036628b70f0fbefba45796f8cdbd1ab091c89af430c4  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.9.43-1.x86_64.rpm)  
SHA256: d58ff2255535d1a543824a24c6723d48da7c5af79509ac375952ee5ceee88fe0  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.9.43_amd64.deb)  
SHA256: 56b52ced5fb494d11abc50dc3ea61375348bf6c8d247c70e5be702aa03e71787  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.9.43-1.x86_64.rpm)  
SHA256: 64ec17d679355107c312c27c74591601d696738df6a2211d8b60966474de88b8  

## Version 3.9.32 - 2022-09-15

[NEW] Support for custom ports  
[NEW] Ability to adjust MTU value for WireGuard connections  
[NEW] (Linux) Ability to change DNS management method to directly modify the '/etc/resolv.conf' file  
[IMPROVED] CLI: `-any` option now connects to a random server  
[IMPROVED] Use the same fastest server configuration for OpenVPN and WireGuard  
[FIX] Uninstalling the app does not remove custom CLI settings  
[FIX] UI: App reconnects when clicking on the protocol section of the main panel  
[FIX] Minor UI fixes  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.9.32.exe)  
SHA256: 2b2df9a1e560b186e333d1020b0dc32d6879dbf00492cdb7b92603b458c4fcb9  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.9.32.dmg)  
SHA256: f2c5ff37ec33c427694ff8680f4e3567ae207ee24038eb7fe41faf7e05e1a417  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.9.32-arm64.dmg)  
SHA256: 74a640f30c22a5197b26de7b2fbc4f2b21cdfe848ed7869ec13d6d94253bee35  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.9.32_amd64.deb)  
SHA256: 77997488e180cbd2a7f770581cd0218bf237e4c4928f557dae14b9548ce1d8c3  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.9.32-1.x86_64.rpm)  
SHA256: 88c364d4366c9bf75db376d3c511421eef493e3458c7f7499502f11a1078b8c0  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.9.32_amd64.deb)  
SHA256: 24cc3909242ce883a7b70997afa49cf14b0865d0fbe10498e0903b4c95cf0ba5  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.9.32-1.x86_64.rpm)  
SHA256: 5a74aedf7df9485083c3398d780fc59cdd8b7c2c63abfb67c8fcb681f55babc7  

## Version 3.9.14 - 2022-08-16

[FIX] (Linux) OpenVPN connection issue  
[FIX] (Linux) Sometimes installer gets stuck during the update  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.9.14_amd64.deb)  
SHA256: 3ae5632b8efc2359ec1fe7964db6e05a6ec841c2d511e848b13add11bc4dae29  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.9.14-1.x86_64.rpm)  
SHA256: 6876a3a3037c31ee5597dcd2fc7f66a55982464a036cb5e8c8b6d1ed79b34e61  

## Version 3.9.9 - 2022-08-10

[FIX] (Windows) OpenVPN connection issue  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.9.9.exe)  
SHA256: c433289eb2b7d0661e5d68d0776b536670f7cc7524b40de17130af73fb7c1fa5  

## Version 3.9.8 - 2022-08-09

[NEW] (Windows/macOS) Option to take part in Beta testing  
[NEW] (macOS) Dock menu in macOS app  
[IMPROVED] Show specific hosts in the favorite menu of the system tray  
[IMPROVED] Updating account status  
[IMPROVED] Reworked UI for sending diagnostic logs  
[IMPROVED] CLI: Show specific host and obfsproxy status in the connection info  
[IMPROVED] (Windows) Updated: WireGuard v0.5.3; OpenVPN v2.5.7 (OpenSSL v1.1.1o)  
[IMPROVED] (macOS) Updated: WireGuard v0.0.20220316; OpenVPN v2.5.7 (OpenSSL v1.1.1o)  
[IMPROVED] (macOS/Linux) Removed unnecessary popup when shutting down  
[IMPROVED] (Linux) Implemented the DNS-change protection mechanism  
[FIX] (Linux) Applications in Split Tunnel were blocked by the firewall  
[FIX] (Linux) VPN's DNS server was missing after waking the system from suspend state on some Linux distributions  
[FIX] (macOS) Application won't close with Cmd+Q  
[FIX] Sometimes UI shows the wrong connection status after connecting to Fastest Server  
[FIX] UI: Update Port and obfsproxy status used in the CLI  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.9.8.exe)  
SHA256: 0b6ab6256a142070d0fb0244234ccd1ec1ddbe0fc631f6ab3c1cc04b035c30a2  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.9.8.dmg)  
SHA256: b64e2fe7ea8296a63cb3e75274c73379731d64dc492778214e6f3ae728e330d1  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.9.8-arm64.dmg)  
SHA256: f915d4954c3f8be86f523b057bca2d976c221628458a8d4ffccfe520417e73f5  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.9.8_amd64.deb)  
SHA256: a196f136a2271913767ab441e0a5e61be9eca81b15c0dfce13d422fdac4ccd5c  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.9.8-1.x86_64.rpm)  
SHA256: 181945f51690fe6fcf17e3ebd50ce0b253e76e4633e22e26d0522550356cac29  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.9.8_amd64.deb)  
SHA256: c818ff2dfb97dcb360d193adfa8aaa261758d1066383e5bed4369f08721e5559  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.9.8-1.x86_64.rpm)  
SHA256: 3abc4235a6f905b7cdf76ef6e4dc49b29f9a715161e69351ed53f2e5b0d80c5a  

## Version 3.9.0 - 2022-07-21

[NEW] Ability to connect to a specific host in a location  
[NEW] (Linux) Snap package support  
[IMPROVED] Updating account status on every application start  
[IMPROVED] Fetch VPN connection ports info from the backend  
[IMPROVED] Warning about existing VPN connection when downgrading subscription plan   
[IMPROVED] (Linux) Using systemd-resolved for controlling DNS settings (when possible)  
[IMPROVED] (Linux) Upgrade retains user settings (applicable only for future updates)  
[IMPROVED] (Linux) Retain firewall configuration on upgrade (applicable only for future updates)  
[FIX] Communication with IVPN servers blocked by ISP’s proxy  
[FIX] (Linux) Intermittent issue with users logged out after upgrade (applicable only for future updates)  
[FIX] (Linux) Firewall rules: potential DNS leak via the link-local IPv6 addresses     

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.9.0.exe)  
SHA256: 45194bc1c45a0a71919ff8a65d873c116011d7c639f4ba51634d7ca99871755c  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.9.0.dmg)  
SHA256: 1401aebc95034a9f4ec78fec99cac37ea8af9f0723e2ad40848f02414c0749c8  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.9.0-arm64.dmg)  
SHA256: 68f44e4bdc734cf659268ab4905d09384b7c7d7899f3a71fba6ae71732579bb4  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.9.0_amd64.deb)  
SHA256: 9278a40f9afc8d0bf92a03c8fc00a216162e5ce1a7858e19f176806f0c581496  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.9.0-1.x86_64.rpm)  
SHA256: 317a7ec94f3e473f2b9b034d0e807f27e38a9e0d29464809158a979eccab8c9b  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.9.0_amd64.deb)  
SHA256: 05152ab69c8388ce75182e48c2b63fc48874a0467cd96952477f804d4b3b6488  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.9.0-1.x86_64.rpm)  
SHA256: 1fcf3cd1d86a9dfd5a1e4181f5cbfe62aad3308ddb69aec6347382a7dc4ec94a  

## Version 3.8.20 - 2022-06-01

[NEW] UI: Ability to see detailed info about the application version  
[IMPROVED] UI: Displaying server info in system tray instead of ‘Random Server’ for Multi-Hop connections  
[IMPROVED] (Linux) ‘Allow LAN’: Monitoring changes in local interfaces configuration   
[IMPROVED] Minor UI improvements  
[FIX] (Linux) Enhanced firewall rules to avoid potential IP leaks  
[FIX] (Linux|macOS) Enhanced firewall rules to avoid potential DNS leaks  
[FIX] ‘Fastest Server’ detection issue when auto-connect on application launch  
[FIX] Keep paused after regeneration of the WireGuard keys   
[FIX] UI: Tray menu now shows only favorite servers for the current protocol  
[FIX] (Linux) DNS configuration issue after VPN resume  
[FIX] (Windows) Service crash when using specific configuration for custom DNS   
[FIX] (Windows) Split Tunnel configuration missing from Settings     

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.8.20.exe)  
SHA256: ee3eaa5dfc4de5ef3644e40c9a0920aa461f5f1288fc8717e01840ac941d4d99  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.8.20.dmg)  
SHA256: 3532012d627699e2c3027a4d1778547c6bcbe64c97bd69015be25fbb93b982cf  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.8.20-arm64.dmg)  
SHA256: 2de8c7f384973982484c6656edd3b9fc01091d1c2c282ee979aceea3b6b88413  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.8.20_amd64.deb)  
SHA256: 7b1037aa224b785c84b44531cc0d5454a328e21082f3011c6d4308de231a007a  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.8.20-1.x86_64.rpm)  
SHA256: 831e7845361574ff5a7529fc47b54405002a209c66e6da7ea3d8bba2d902cc79  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.8.20_amd64.deb)  
SHA256: 4949949d219bb746fc0d5b9d9eb3e7e6a0f26aa01d32d9cda687b5461d82119c  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.8.20-1.x86_64.rpm)  
SHA256: e01180b5325ccb2c21eb714bae7b74203d1e3947c9b38c88d08a039fcf1fed64  

## Version 3.8.7 - 2022-05-04

[NEW] Enhanced App Authentication  
[NEW] Custom Firewall exceptions  
[FIX] (Windows) Firewall was blocking connections on port 53  
[FIX] (Windows) Service crash when using specific custom DNS configuration  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.8.7.exe)  
SHA256: d9b16b2b3eff87e32a09b74558d3570d1b1ca82b77e9483ceb7aac179a568661  

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.8.7.dmg)  
SHA256: e1e9d166b55af4103a85198b4d195622a1a1820b8d39f4201b9ff012a6b3bab1  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.8.7-arm64.dmg)  
SHA256: 7ea309b19f344d291395b362feaa190271d5ac4c37f4480978cb6b8e442323b2  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.8.7_amd64.deb)  
SHA256: 0ada7f81f7525ce6352535754a12f74b32437e5725c95c4fedd696fa2620a051  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.8.7-1.x86_64.rpm)  
SHA256: 9661b837e3873bb4ea8ba24d4f54057086186531f23cec4a7c0ff7f7e178bb84  

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.8.7_amd64.deb)  
SHA256: 27b9b5365569ae177878eea17a390facc522333839a5ae91b564829ce0e67fe3  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.8.7-1.x86_64.rpm)  
SHA256: be63d92659d1f5268e7edb53b4635886705589a5440ba19b2b7a6c588a87f4b8  

## Version 3.7.0 - 2022-04-04

[NEW] (macOS/Linux) Configure custom DNS over HTTPS  
[NEW] (Windows) Configure custom DNS over HTTPS (Windows 8/10)  
[NEW] (Windows/Linux) Last used apps for Split Tunneling now show on the top of the list  
[NEW] Added new connection ports (UDP 80, UDP 443)  
[FIX] Vulnerability fixes  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.7.0.exe)  
SHA256: 41dcf6f286e9d4c29e7d857af270b794e07e5dbce5767e57f6e7f27d64ce56e2   

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.7.0.dmg)  
SHA256: 8870487f9ca1e24f2c93ab1b926380f1ec7fe5abb7252a31411f90f234db20ff  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.7.0-arm64.dmg)  
SHA256: ee691f502f6614c3824b58a52ecdcd79039d27b7df08cdc94258a8dd8c290b57  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.7.0_amd64.deb)  
SHA256: 7603e2c2edd2355068f34afceaf7d3ef6f0bfff2310b341faa8bd16933f077f6  
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.7.0-1.x86_64.rpm)  
SHA256: bf99eaac31708978d98412fdbeeb5be544f91737f2423804b40851d98be59676

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.7.0_amd64.deb)  
SHA256: 012f7ff53c0c42e94f53963cdbe1e5354222fedc1611c380b3e8f06fab2302a9  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.7.0-1.x86_64.rpm)  
SHA256: 94351142da24e75ed97fca4721a5d871b638a9f648cd89394089d5bb0af2827a

## Version 3.6.4 - 2022-03-15

[NEW] Configure custom DNS over HTTPS (Windows 11)  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.6.4.exe)  
SHA256: 022de7299f6578822f2225d39b070fe8551c06e4a164980716be867f967a969b  

## Version 3.5.2 - 2022-01-31

[FIX] (Windows) Split Tunneling issues  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.5.2.exe)  
SHA256: 19c64882a2975c6c0db538c7698ae51ccf18fdbc32c983fec6ef8b621b3cfaa5  

## Version 3.5.1 - 2022-01-20

[NEW] (Linux) Split Tunneling  
[IMPROVED] Ping marker coloring logic

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.5.1_amd64.deb)  
SHA256: 34937d312a6e69bc2f72bf92d3b791dff222157941e0b616578be2cb5410c5eb   
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.5.1-1.x86_64.rpm)  
SHA256: 7fb0c81db801abe1e3e751750fed671ffce508159f2de665ba69b755ec37a86a

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.5.1_amd64.deb)  
SHA256: c9c0c1ed554bd44a3d3ba7f25b5720f40528e40b281a4872264b24555c78c8a5  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.5.1-1.x86_64.rpm)  
SHA256: 9ab6d4507c455d957dfe027c008b8644520ea1a82936a38718b7facc25120be6

## Version 3.4.5 - 2021-11-23

[FIX] (Windows|Linux) Fixed the problem of opening a minimized application  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.4.5.exe)  
SHA256: 6569254e6368a20e306ff36d7bd49d90197f93a4dc71a98de56a5f23b48bdfda  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.4.5_amd64.deb)  
SHA256: 5956cb9e99a4133b29680c7e29e081850b0f8e06c3dfd5b966fbaa17d466bedf   
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.4.5-1.x86_64.rpm)  
SHA256: d04dad1c4c0d197fdac5a7555aa0823ce02012e1e416100e78184474b67f123d

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.4.5_amd64.deb)  
SHA256: 49f3be89d3758b41b7057ea411be3fda16b785b1af9ce46b87539602ca8246ea  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.4.5-1.x86_64.rpm)  
SHA256: 6773d51a14cf5f9ecc469a29889c7b8abc1f3408db8fdaa70b3fb9c54a074f6b

## Version 3.4.4 - 2021-11-17

[IMPROVED] (Windows) Updated WireGuard: v0.5.2  
[IMPROVED] (Windows) Firewall rules  
[FIX] (Windows) The Wireguard binary is signed to avoid false positive detections by antiviruses  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.4.4.exe)  
SHA256: e8ec9a4a5c15fadf8343e3852d0cbb6dbb7ba8b9f5893a28fa9e37aa7c631b64  

## Version 3.4.1 - 2021-11-09

[IMPROVED] (Windows) The TAP driver is signed by new certificate  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.4.1.exe)  
SHA256: 58f084098edbaeaadc3e4ee13a4bbbd9af9e6ce74cc4e5fba0cc1be73f2d3cfd  

## Version 3.4.0 - 2021-11-08

[NEW] Multi-Hop for WireGuard protocol  
[NEW] Option to reset app settings on logout  
[NEW] Option to keep Firewall state on logout  
[NEW] CLI option to show all servers and to connect to specific server  
[NEW] (Linux) Obfsproxy now works on Linux  
[IMPROVED] Speed up the response timeout to API server  
[IMPROVED] Force automatic WireGuard key regeneration if the rotation interval has passed  
[IMPROVED] (Windows) Updated WireGuard: v0.4.9  
[IMPROVED] (Windows) Updated: OpenVPN: v2.5.3; OpenSSL: 1.1.1k  
[IMPROVED] (macOS) Updated: OpenVPN: v2.5.3; OpenSSL: 1.1.1k  
[IMPROVED] (macOS) Updated WireGuard binaries: wireguard-go: v0.0.20210424; wireguard-tools v1.0.20210914  
[IMPROVED] (Linux) WireGuard-tools integrated into a package (for a kernel since 5.6, no dependencies required to use WireGuard)  
[FIX] Fastest server settings were ignored in some cases  
[FIX] Option to run multiple UI instances in some cases  
[FIX] Server selection issues  
[FIX] Other minor issues and improvements  
[FIX] (Windows) Compatibility with Windows Server  
[FIX] (Windows) IVPN Firewall rules overlap blocking rules from Windows Firewall  
[FIX] (Windows) Icons created in %temp% each time app is launched  
[FIX] (macOS) Unable to start WireGuard connection if more than 10 utunX devices configured  
[FIX] (Linux) "Allow LAN traffic" does not persist after a restart  
[FIX] (Linux) UI crash after some Linux distribution updates  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.4.0.exe)  
SHA256: 01d876ad506ccf9def6c8ded2c104b740bb3093d728ad52168aecf597113f7d4   

[Download IVPN Client for macOS (Intel)](https://repo.ivpn.net/macos/bin/IVPN-3.4.0.dmg)  
SHA256: ca9d45f7df2eb95fa5f57ada9012d6add95113635b74f21df36c40725687b3f2  
[Download IVPN Client for macOS (M1)](https://repo.ivpn.net/macos/bin/IVPN-3.4.0-arm64.dmg)  
SHA256: 8a1f4bb2c01f289b2ca241b86c0b5eec4b1225de06777d076d2ef534e20e7481  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.4.0_amd64.deb)  
SHA256: fad328c95679c983d162d117e909c4c0b5eacd7b5dd54b8de7e1a1c4dbeca64c   
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.4.0-1.x86_64.rpm)  
SHA256: 933c397078be24eba87cce63c3d49b507e62efb623a34f9349725461de719130

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.4.0_amd64.deb)  
SHA256: 7e50c58ed16c5817e79b253e7b198a76c4660218a1e236598a59a288eaaf89e3  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.4.0-1.x86_64.rpm)  
SHA256: cf95c4e07912aa03c7596d56b31d323664efbf44469cc9fee54771800d96d1db

## Version 3.3.40 - 2021-09-14

[NEW] (Windows) Split Tunneling  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.3.40.exe)  
SHA256: 9875bc8ee2124464b66fa70555270865caf03c827e4323fdf6fb2a7a83589606  

## Version 3.3.30 - 2021-08-31

[NEW] Preparation for 2FA in CLI  
[IMPROVED] (Linux) The installation path changed from '/usr/local/bin' to '/usr/bin'  
[FIXED] Server selection order incorrect when sorted by country  
[FIXED] (Linux) Removed unnecessary paths from package which may lead to conflict with other software  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.3.30.exe)  
SHA256: 981bce29c543df2485687edcc9383e1fe5acc343cba0d8b8ea8beada8c57a3e6   

[Download IVPN Client for macOS](https://repo.ivpn.net/macos/bin/IVPN-3.3.30.dmg)  
SHA256: 7155967dda8f53580ab2d158fa57b447efe0c40a29f28b884bf33fc0f8fcb12d  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.3.30_amd64.deb)  
SHA256: 89d20099b8e36b704106074c60a89ff189ff6e99e999a3ae748801b3ba76bd07   
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.3.30-1.x86_64.rpm)  
SHA256: 7b432c77c85bee2267bbbb218ee761b8c036208b14350476afa7179b133ad0a3

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.3.30_amd64.deb)  
SHA256: 229d70cfcb7bee5a7a888b5864797a5fec09cbd320f4d1a0c374cd30b17b2452  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.3.30-1.x86_64.rpm)  
SHA256: f7a77300bcc261af44e0d146970a89a4598d54be3161b2913516051d57f13a52  

## Version 3.3.20 - 2021-06-29

[NEW] IPv6 inside WireGuard tunnel  
[NEW] IPv6 connection info  
[NEW] New option in settings ‘Allow access to IVPN server when Firewall is enabled’  
[NEW] (Windows) Contrast tray icon (black or white; depends on Windows color theme)  
[FIXED] VPN was active after reboot when connected to Trusted WIFI  
[FIXED] Sometimes application was failing to connect to IVPN daemon  
[FIXED] (Windows) The daemon service was not starting when the 'Windows Events Logs' service is not running  
[FIXED] (macOS) WireGuard compatibility with old macOS versions  

[Download IVPN Client for Windows](https://repo.ivpn.net/windows/bin/IVPN-Client-v3.3.20.exe)  
SHA256: 02b312a0edf21765229c1e8f12e48bec2539b241e05afcda65b90c4f9730d950   

[Download IVPN Client for macOS](https://repo.ivpn.net/macos/bin/IVPN-3.3.20.dmg)  
SHA256: 14d4f51e2167a9c07e56d35de7632570168c69ed93beb4711128e5ddd04ca67f  

[Download IVPN Client for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn_3.3.20_amd64.deb)  
SHA256: 9972a0b55e1383df67357d046db238b29c70b7865dcba959037da17b7439f20a   
[Download IVPN Client for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-3.3.20-1.x86_64.rpm)  
SHA256: 469d5b41b2092612cf82a9b269e4205ca4ebbcc651b5fbd196be8e138f005487

[Download IVPN Client UI for Linux (DEB)](https://repo.ivpn.net/stable/pool/ivpn-ui_3.3.20_amd64.deb)  
SHA256: a9cd6f2e2e1c7f981d0788b0f6e381c56e8b44f44daad217b66b652d5eead947  
[Download IVPN Client UI for Linux (RPM)](https://repo.ivpn.net/stable/pool/ivpn-ui-3.3.20-1.x86_64.rpm)  
SHA256: 80e37b4c2c00fa89411e6bf403b72c60b66c09c2bd0ec0f0cdf264e76de00492  

## For old versions of IVPN for Desktop please refer to:

[Windows/macOS and Linux UI](https://github.com/ivpn/desktop-app-ui2/blob/master/CHANGELOG.md)  
[Linux (cli)](https://github.com/ivpn/desktop-app-cli/blob/master/CHANGELOG.md)
