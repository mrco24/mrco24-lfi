# MRCO24-LFI

go get -u github.com/mrco24/mrco24-lfi

#Tools - Use 

mrco24-lfi -u urls.txt -p payloads.txt -o w.txt
#Remove =
sed 's/=.*$/=/' url.txt | anew | tee -a live_url.txt
#wthout = run tools
mrco24-lfi -u good_url.txt -p payloads.txt -o w.txt -v
