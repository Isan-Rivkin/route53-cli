# Route53 Query 

Move fast with custom CLI for Route53. 

Get info about your Route53 records from the Command line - quickly!

* Direct link to resource on AWS Console.
* Automatic verification of Nameservers against the real world.
* Recursive search for all results via `-R` flag.
* Supports any `~/.aws/credentials` profile.
* Works directly against AWS API.


## The Problem

Ever wondered "what records are behind this R53 record I have?" 

i.e  Where does `app.foo.goo.website.com` points to in R53?  

## Without `route53-cli`: 

Go to browser -> AWS console -> login -> Route53 -> find the hosted zone -> find the record 

**Even if you find the record in the AWS console you can't be sure this is the source of truth since the record could be defined in multiple route 53 hosted zones of the organization. The only way to verify is by comparing the Nameservers in the real world.**

## The solution 

- The CLI will use the default AWS login methods and query the Route53 API in a smart way to find the record. 
- By default the nameservers will be verified against the real world via Dig implementation (turn off via `--skip-ns`).
- By using `-R` the query will continue recursivly and expand all the records until the "leaf".
- If the record value is an AWS resource it will output the URL to AWS console for quick access.

## With `route53-cli`:

**Input**

```bash
r53 -R -r 'app.foo.goo.website.com'
``` 

**Output**

```bash

┌─────────────────┬────────────────────────────┬───────────────┬─────────┐
│ HOSTED ZONE     │ ID                         │ TOTAL RECORDS │ PRIVATE │
├─────────────────┼────────────────────────────┼───────────────┼─────────┤
│ website.com.    │ /hostedzone/ABFDCEWQSFDSFD │           127 │ false   │
└─────────────────┴────────────────────────────┴───────────────┴─────────┘

+---+-----------------------|------------------------|--------------------------------------|---------+
| # | Record                | Target                 | Console Link                         | TYPE    |
+---+-----------------------|------------------------|--------------------------------------|---------+
| 1 | *.foo.goo.website.com | r-re1.website.com.     |                                      | A       |
+---+-----------------------|------------------------|--------------------------------------|---------+
| 2 | *.foo.goo.website.com | r-re2.website.com.     |                                      | A       |
+---+-----------------------|------------------------|--------------------------------------|---------+
| 3 | r-re1.website.com.    | lb1.elb.amazonaws.com  | https://console.aws.amazon.com/alb-1 | A       |
+---+-----------------------|------------------------|--------------------------------------|---------+
| 4 | r-re2.website.com.    | lb2.elb.amazonaws.com  | https://console.aws.amazon.com/alb-2 | A       |
+---+-----------------------|------------------------|------------------------------------------------+

```

# Install 

### Brew 

MacOS (and ubuntu supported) installation via Brew:

```bash
brew tap isan-rivkin/toolbox
brew install r53
```

### Download Binary

1. [from releases](https://github.com/Isan-Rivkin/route53-cli/releases)

2. Move the binary to global dir and change name to `r53`:

```bash
cd <downloaded zip dir>
mv r53 /usr/local/bin
```

### Install from Source

```bash
git clone 
cd route53-cli
make install BIN_DIR='/path/to/target/bin/dir'
```

### Version Check 

The CLI will query anonymously a remote version server to check if the current version of the CLI is updated.
If the current client version indeed outdated the server will return instructions for update. 

The server will add the request to a hit counter stored internaly for usage metrics. 

**None of the user query are passed to the server, only OS type and version.**

**The route53 querys themselves are done directly via the AWS Api.**

This behaviour is on by default and can be optouted out via setting the envrionment variable `R53_VERSION_CHECK=false`. 

### How it works 

Example pseudocode: 

```python
# i.e https://example.com/p/a?ok=11&not=23
# into example.com 
record = '*.<a>.<b><c>.<d>'
record = strip_non_domain(record)
record = split(record)
 
hasWildCard = record[0] == '*'
len == record.length
# a 
if len == 1 and not hasWildCard: 
    lookup(1)
# *.a
if len == 2 and hasWildCard:  
    lookup(1)
# a.b
if len == 2 and not hasWildCard:  
    lookup(2)
# *.a.b
if len == 3 and hasWildCard: 
    lookup(2)
# a.b.c
if len == 3 and hasWildCard: 
    if not lookup(2)
        lookup(3)
# *.a.b.c.d.e
if len > 3:
    for chunksNum in (2, len):
        if (chunksNum == len and hasWildCard):
            return None
        if res = lookup(chunksNum)
            return res 

def lookup(dnsChunk):
    zoneIds = getZoneIds(dnsChunk)
    for zoneId in zoneIds:
        if zoneId name == dnsChunk: 
            aliasTargets = getAliasTargets(record, zoneId)
            return aliasTargetes
```
