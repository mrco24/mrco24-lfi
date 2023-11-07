<h1 align="center">
<a href="https://cooltext.com"><img src="https://images.cooltext.com/5678562.png" width="459" height="110" alt="MRCO24 - LFI" /></a>
</h1>
<h4 align="center">LFI Detection</h4>
<p align="center">
  <a href="https://github.com/mrco24/mrco24-lfi">
    <img src="https://img.shields.io/badge/Mrco24-Lfi_Detection-green">
  </a>
   <a href="https://github.com/mrco24/mrco24-lf">
    <img src="https://img.shields.io/static/v1?label=Update&message=V1.0&color=green">
  </a>
  <a href="https://twitter.com/mrco24">
      <img src="https://img.shields.io/twitter/follow/mrco24?style=social">
  </a>
</p>

# Current Features:
<img width="550" alt="2023-11-04 22_40_48-kali linux  Running  - Oracle VM VirtualBox" src="https://github.com/mrco24/mrco24-lfi/assets/95030367/5584d1d7-29da-4442-9ee4-36980cadb09d">

- **High-Speed Scanning:** Lightning-fast scanning of web applications for Local File Inclusion (LFI) vulnerabilities.

- **Multi-URL Support:** Simultaneously scan multiple URLs for LFI issues.

- **Comprehensive Payload Scanning:** Extensive payload library for thorough testing.

- **Rapid Vulnerability Detection:** Quickly identifies LFI vulnerabilities with minimal false positives.

- **Multi-Threaded:** Utilizes multi-threading for efficient scanning.

- **Vulnerability URL Output:** Provides a list of vulnerable URLs for further analysis or action.

`mrco24-lfi` requires **go1.19** to install successfully. Run the following command to install the latest version: 
```
go get -u github.com/mrco24/mrco24-lfi
```
# Usage:
```
mrco24-lfi -u live_url.txt -p payloads.txt -o output.txt -v
```
# Remove =
```
sed 's/=.*$/=/' url.txt | anew | tee -a live_url.txt
```






