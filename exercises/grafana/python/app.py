from flask import Flask, request, redirect
from os import environ
import redis
import json
from itertools import zip_longest


app = Flask(__name__)


if environ.get('REDIS_SERVER') is not None:
   redis_server = environ.get('REDIS_SERVER')
else:
   redis_server = 'localhost'

if environ.get('REDIS_PORT') is not None:
   redis_port = int(environ.get('REDIS_PORT'))
else:
   redis_port = 6379

if environ.get('REDIS_PASSWORD') is not None:
   redis_password = environ.get('REDIS_PASSWORD')
else:
   redis_password = ''

if environ.get('REDIS_LEADERBOARDS') is not None:
   redis_leaderboards = environ.get('REDIS_LEADERBOARDS').split(",")
else:
   redis_leaderboards = []

if environ.get('REDIS_LEADERBOARD_SET') is not None:
   redis_leaderboard_set = environ.get('REDIS_LEADERBOARD_SET')
else:
   redis_leaderboard_set = ""

if environ.get('REDIS_SCANPREFIX') is not None:
   redis_scanprefix = environ.get('REDIS_SCANPREFIX')
else:
   redis_scanprefix = "*"

if environ.get('REDIS_TOPK') is not None:
   redis_topk = environ.get('REDIS_TOPK')
else:
   redis_topk = 10




client = redis.Redis(
   host=redis_server,
   port=redis_port,
   password=redis_password)

def batcher(iterable, n):
   args = [iter(iterable)] * n
   return zip_longest(*args)

@app.route('/', methods = ['GET'])
def ping():
   return "OK"

@app.route('/search', methods = ['POST'])
def search():
   if redis_leaderboard_set != "":
      b = []
      for x in client.smembers(redis_leaderboard_set):
         b.append(x.decode("UTF-8"))
      return json.dumps(b)
   elif len(redis_leaderboards) > 0:
      return json.dumps(redis_leaderboards)
   else:
      l = []
      for keybatch in batcher(client.scan_iter(redis_scanprefix),5000):
         for x in keybatch:
            if x != None:
               if client.type(x.decode("utf-8")).decode("utf-8") == "zset":
                  l.append(x.decode("utf-8"))
      return json.dumps(l)

@app.route('/query', methods = ['POST'])
def query():
   req = request.get_json()
   k = req['targets'][0]['target']
   table_def = [
      { "type": "table",
      "columns": [ { "text": "Value", "type": "string" }, { "text": "Score", "type": "string" } ],
      "rows": [] } ]
   x = client.zrevrange(k, 0, redis_topk, withscores=True)
   for j in x:
      table_def[0]["rows"].append([j[0].decode('utf-8'), j[1]])
   return json.dumps(table_def)

if __name__ == '__main__':
   app.debug = True
   app.run(port=5000, host="0.0.0.0")
