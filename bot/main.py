import os
import random
import requests
import subprocess
import time
import traceback
import twitter
import xml.etree.ElementTree as ET

RATE = 60 * 30

TWITTER_CONSUMER_KEY = None
TWITTER_CONSUMER_SECRET = None
TWITTER_ACCESS_TOKEN_KEY = None
TWITTER_ACCESS_TOKEN_SECRET = None

try:
    from config import *
except ImportError:
    print 'no config found!'

def rcsb(structure_id):
    cmd = 'rcsb %s' % structure_id
    subprocess.call(cmd, shell=True)

def twitter_api():
    return twitter.Api(
        consumer_key=TWITTER_CONSUMER_KEY,
        consumer_secret=TWITTER_CONSUMER_SECRET,
        access_token_key=TWITTER_ACCESS_TOKEN_KEY,
        access_token_secret=TWITTER_ACCESS_TOKEN_SECRET)

def tweet(status, media):
    api = twitter_api()
    api.PostUpdate(status, media)

def random_structure_id():
    with open('structures.txt') as fp:
        lines = fp.read().split()
        return random.choice(lines).strip().upper()

def structure_title(structure_id):
    url = 'http://www.rcsb.org/pdb/rest/describePDB?structureId=%s' % structure_id
    r = requests.get(url)
    root = ET.fromstring(r.content)
    for pdb in root.findall('PDB'):
        return pdb.attrib.get('title')

def structure_url(structure_id):
    return 'http://www.rcsb.org/pdb/explore.do?structureId=%s' % structure_id

def trunctate_text(text, max_length):
    if len(text) > max_length:
        return text[:max_length-3] + '...'
    return text

def tweet_text(structure_id):
    title = structure_title(structure_id)
    title = trunctate_text(title, 110)
    url = structure_url(structure_id)
    return '%s: %s %s' % (structure_id, title, url)

def generate():
    structure_id = random_structure_id()
    status_text = tweet_text(structure_id)
    print status_text
    out_path = '%s.png' % structure_id
    print 'rendering image'
    rcsb(structure_id)
    if os.path.exists(out_path):
        print 'uploading to twitter'
        tweet(status_text, out_path)
        print 'done'
    else:
        print 'failed'

def main():
    previous = 0
    while True:
        now = time.time()
        if now - previous > RATE:
            previous = now
            try:
                generate()
            except Exception:
                traceback.print_exc()
        time.sleep(5)

if __name__ == '__main__':
    main()
