import re
import requests

url = 'https://data.rcsb.org/rest/v1/holdings/current/entry_ids'

r = requests.get(url)

structures = re.findall(r'\w{4}', r.text)
structures.sort()

print('\n'.join(structures))
