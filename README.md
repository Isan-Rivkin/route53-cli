# Route53 Query 

Command line utility for Route53. 

# Use Cases

### Where does <some dns record> points to? 

<b> Input </b>

```bash
r53 *.foo.goo.website.com
``` 

<b> Output </b>

```bash
+---+--------+-----------------------+------------------------+
| # | ZoneId | Record                | Target                 |
+---+--------+-----------------------+------------------------+
| 1 | ABC    | *.foo.goo.website.com | lb1.elb.amazonaws.com  |
+---+--------+-----------------------|------------------------+
```

