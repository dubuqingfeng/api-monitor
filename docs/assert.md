## HTTP API Assert Type

### HTTP Status

http.Status:

200,
202

### response json

1. match
2. jsonpath 

### Cookie

cookie

### Header

header

## Content

```json
{
    "status": [{
        "type": "equals",
        "value": "200"
    }, {
        "type": "equals",
        "value": "200"
    }],
    "jsonpath": [
    {
        "type": "equals",
        "key": "$.b[? @.key==\"c\"].value",
        "value": "result"
    }, 
    {
        "type": "contains",
        "key": "$.b",
        "value": "result"
    }],
    "cookie": {
    },
    "header": {
    }
}
```
