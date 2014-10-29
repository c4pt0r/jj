# jj

##Overview

Redis alike NoSQL Memory DB for JSON document, support JSON Path

```
jdocset [key] [val]
jdocget [key]

jset [key] [json path] [val]
jget [key] [json path]

jincr [key] [json path] [integer]

jpush [key] [json path] [val]
jpop [key] [json path]

save
bgsave
```

Example:

```
127.0.0.1:9999> jdocset a {}
OK
127.0.0.1:9999> jdocget a
"{}"
127.0.0.1:9999> jset a a [1,2,3,4,5]
OK
127.0.0.1:9999> jdocget a
"{\"a\":[1,2,3,4,5]}"
127.0.0.1:9999> jget a a
"[1,2,3,4,5]"
127.0.0.1:9999> jget a a[0]
"1"
127.0.0.1:9999> jincr a a[0] 100
OK
127.0.0.1:9999> jget a a[0]
"101"
127.0.0.1:9999> jdocget a
"{\"a\":[101,2,3,4,5]}"
127.0.0.1:9999> jpush a a {}
OK
127.0.0.1:9999> jget a a
"[101,2,3,4,5,{}]"
127.0.0.1:9999> jpop a a
"101"
127.0.0.1:9999> jpop a a
"2"
127.0.0.1:9999> jpop a a
"3"
127.0.0.1:9999> jpop a a
"4"
127.0.0.1:9999> jpop a a
"5"
127.0.0.1:9999> jpop a a
"{}"
127.0.0.1:9999> jdocget a
"{\"a\":[]}"

```

* Python APIs
* Golang APIs

