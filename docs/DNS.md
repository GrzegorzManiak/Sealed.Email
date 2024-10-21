# Sample DNS Records

## `A` Records

These are the root **MX** records for the **NOISE** email service.
These are the servers that other domains will point to in order to use the **NOISE** email service.

```dns
doom.mx.noise.email.	1	IN	A	188.141.60.114 ; NOISE-MX
madlib.mx.noise.email.	1	IN	A	188.141.60.114 ; NOISE-MX
west.mx.noise.email.	1	IN	A	188.141.60.114 ; NOISE-MX
```

## `MX` Records

This is a sample setup for a subdomain of the **NOISE** email service, **beta.noise.email**.
I am currently using this for testing, a customer would be asked to point their **MX** 
records to the **NOISE** servers if they wanted to use their own domain with the 
**NOISE** email service.

```dns
beta.noise.email.	1	IN	MX	99 west.mx.noise.email.   ; NOISE-BETA-DOMAIN
beta.noise.email.	1	IN	MX	90 madlib.mx.noise.email. ; NOISE-BETA-DOMAIN
beta.noise.email.	1	IN	MX	50 west.mx.noise.email.   ; NOISE-BETA-DOMAIN
```

## `TXT` Records

These are the **TXT** records for the **NOISE** email service.

The first record is the **SPF** record for the **beta.noise.email** domain, it points
to the **_spf.mx.noise.email** record, which is the **root SPF** record for the **NOISE** email service.

The second record is the **DMARC** record for the **beta.noise.email** domain, it just explains
what to do with the **DMARC** reports.

The third record is the **SPF** record for the **_spf.mx.noise.email** domain, it is the **root SPF** record

```dns
beta.noise.email.	1	IN	TXT	"v=spf1 include:_spf.mx.noise.email -all" ; NOISE-BETA-DOMAIN
_dmarc.beta.noise.email.	1	IN	TXT	"v=DMARC1; p=none; rua=mailto:dmarc-reports@beta.noise.email; ruf=mailto:dmarc-failures@beta.noise.email; sp=none; aspf=r;"
_spf.mx.noise.email.	1	IN	TXT	"v=spf1 ip4:188.141.60.114 -all" ; NOISE-SPF
```