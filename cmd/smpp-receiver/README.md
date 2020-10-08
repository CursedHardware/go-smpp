# SMPP Receiver

This simple program receive `deliver_sm` message and call a local hook function

## Configuration

Write to `configure.json` file:

```json
{
  "hook": "./send-email.py",
  "default_account": {
    "smsc": "target ip:target port",
    "password": "your password"
  },
  "devices": [
    { "system_id": "tenant-1" },
    { "system_id": "tenant-2" }
  ]
}
```

Sample hook script:

```python
#!/usr/bin/env python3
import json
import requests

payload = json.load(sys.stdin)
"""
{
    "smsc": "[login smsc address]",
    "system_id": "[login system id]",
    "system_type": "[login system type]",
    "source": "[source phone number]",
    "target": "[target phone number]",
    "message": "[merged message content]",
    "deliver_time": "[iso8601 formatted]"
}
"""

to_addresses = {
    "tenant-1": "tenant one addresses",
    "tenant-2": "tenant two addresses",
}

data = {
    "to": to_addresses[payload["system_id"]],
    "from": "%(target)s <[your from address]>" % payload,
    "subject": payload["source"],
    "text": "%(message)s\n\nDate: %(deliver_time)s" % payload
}

requests.post(
    "https://api.mailgun.net/v3/[your api domain]/messages",
    auth=("api", "your api token"),
    data=data,
)
```
