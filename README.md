# Route53 Query 

Command line utility for Route53. 

Get info about your Route53 records from the Command line - quickly!

Instead of: 
Go to browser -> aws console -> login -> route53 -> find the hosted zone -> find the record 

Just do 1 command from the cli :) 

# Use Cases

### Where does `app.foo.goo.website.com` points to in R53? 

<b> Input </b>

```bash
r53 -r 'app.foo.goo.website.com'
``` 

<b> Output </b>

```bash
+---+--------+-----------------------+------------------------+
| # | ZoneId | Record                | Target                 |
+---+--------+-----------------------+------------------------+
| 1 | ABC    | *.foo.goo.website.com | lb1.elb.amazonaws.com  |
+---+--------+-----------------------|------------------------+

```

# Install 

### Download Binary

1. [from releases](https://github.com/Isan-Rivkin/route53-cli/releases)

2. move the binary to global dir and change name to r53:

```bash
cd <downloaded zip>
mv route53-cli r53
mv r53 /usr/local/bin
```

### Install from Source

```bash
git clone 
cd route53-cli
make install BIN_DIR='/path/to/target/bin/dir'
```

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