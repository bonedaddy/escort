param (
    [string]$encoded_data = "ykjNyclXKM8vykkBBAAA//8="
)
# encoding from https://gist.github.com/vortexau/13de5b6f9e46cf419f1540753c573206
$data = [System.Convert]::FromBase64String($encoded_data)
$ms = New-Object System.IO.MemoryStream
$ms.Write($data, 0, $data.Length)
$ms.Seek(0,0) | Out-Null

$sr = New-Object System.IO.StreamReader(New-Object System.IO.Compression.DeflateStream($ms, [System.IO.Compression.CompressionMode]::Decompress))

while ($line = $sr.ReadLine()) {  
    $line
}