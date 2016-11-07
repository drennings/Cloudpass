from itertools import product
import time
import hashlib
import requests
import json
import threading
from flask import Flask
from flask import request
app = Flask(__name__)

alphabet = '0123456789'
length = 8


class Job:
    def __init__(self, worker_id, master_addr,
                 hash_str, hash_type, share, capacity):
        self.worker_id = worker_id
        self.master_addr = master_addr
        self.hash_str = hash_str
        self.hash_type = hash_type
        self.share = share
        self.capacity = capacity
        self.progress = 0
        self.heartbeat_interval = 10

    def start(self):
        self.schedule_heartbeats()

        # Bruteforce our share of the solution space
        self.time_start = time.time()
        for i in range(length):
            for p in product(alphabet, repeat=i):
                sol = ''.join(p)
                md5 = hashlib.md5()
                sol_ascii_bytes = bytes(sol, 'utf-8')
                md5.update(sol_ascii_bytes)

                if self.hash_str == md5.hexdigest():
                    self.solution = sol
                    self.time_stop = time.time()

    def schedule_heartbeats(self):
        def send_heartbeat():
            requests.post(self.master_addr, self.toJSON())
        set_interval(send_heartbeat, self.heartbeat_interval)

    def notify_master(self):
        requests.post(self.master_addr, self.solution)

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__,
                          sort_keys=True, indent=4)


def set_interval(func, sec):
    def func_wrapper():
        set_interval(func, sec)
        func()
    t = threading.Timer(sec, func_wrapper)
    t.start()
    return t


# Called by master to start the worker
@app.route('/start', methods=['POST'])
def start():
    worker_id = request.args.get('workerId')
    master_addr = request.args('masterAddr')
    hash_str = request.args.get('hash')
    hash_type = request.args.get('type')
    share = request.args('share')
    capacity = request.args('cap')

    job = Job(worker_id, master_addr, hash_str, hash_type, share, capacity)
    job.start()

    return 'Starting job.'

if __name__ == "__main__":
    app.run()
