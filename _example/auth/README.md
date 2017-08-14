## Auth

Digota's Authentication is based on private ClientCAs, basically we are creating CA and signing any certificate we want to approve with same CA. 

> How about revoking certificate? The CRL approch here is whitelist instead of blacklist, just remove client serial from your config.

The easiest way to generate certificates is using [certstrap](https://github.com/square/certstrap).

##### Create CA

```bash
$ certstrap init --common-name "ca.company.com"
output:
Created out/ca.company.com.key
Created out/ca.company.com.crt
Created out/ca.company.com.crl
```

##### Create Server Certificate

```bash
$ certstrap request-cert --domain server.company.com
output:
Created out/server.company.com.key
Created out/server.company.com.csr
```

##### Create Client Certificate

```bash
$ certstrap request-cert --domain client.company.com
output:
Created out/client.company.com.key
Created out/client.company.com.csr
```

##### Sign Certificate

```bash
$ certstrap sign --CA "ca.company.com" client.company.com
output:
Created out/client.company.com.crt from out/client.company.com.csr signed by out/ca.company.com.key
```

##### Approve Certificate

Take the certificate serial

```bash
$ openssl x509 -in out/client.com.crt -serial | grep -Po '(?<=serial=)\w+'
output:
serial=A2FF9503829A3A0DDE9CB87191A472D4
```
Append the serial and scopes(`WRITE`,`READ`,`WILDCARD`) to your config

```yml
...
...
...
clients:
- serial: "A2FF9503829A3A0DDE9CB87191A472D4"
  scopes:
  - READ
  - WRITE
```