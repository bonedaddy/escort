# Escort

`escort` is a shitty atempt at AV bypassing using a combination of powershell, "Invoke-Executable" (IEX) cmdlets, and DNS lookups. Essentially the idea is to have a small on-disk powershell script that does a DNS lookup for one or more TXT records. Each TXT record is an entire powershell script or part of a powershell script which consists of DEFLATE compressed data that is base64 encoded.


# Caveats

Due to the usage of DEFLATE the powershell script you are compressing and base64 encoding must be a minimum of 45 characters in length, otherwise you should skip the deflate process.