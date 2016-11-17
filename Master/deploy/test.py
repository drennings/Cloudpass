import requests
import sys
res = requests.post("http://10.0.0.1", {"a": "b"})
print(res.status_code, res.text)
sys.exit()
