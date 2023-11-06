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

# Installation:
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
# Current Features:
- This script will collec





