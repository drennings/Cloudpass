from itertools import product
import time
import hashlib
import requests
import json
import sys
import threading
from flask import Flask
from flask import request
app = Flask(__name__)

alphabet = '0123456789'
length = 8


class Worker:
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
        """ Start bruteforcing the password while sending heartbeats
        to the master every self.heartbeat_interval seconds to show we
        are alive.
        """
        set_interval(self.notify_master, self.heartbeat_interval)

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
                    self.notify_master()
                    sys.exit(0)

    def notify_master(self):
        """ Sends our state to the master, this is done every
        'heartbeat_interval' seconds as well as when the solution is found.
        """
        requests.post(self.master_addr, self.toJSON())

    def toJSON(self):
        """ Dumps our state to JSON such that we can send it to the master
        over HTTP.
        """
        return json.dumps(self, default=lambda o: o.__dict__,
                          sort_keys=True, indent=4)


def set_interval(func, sec):
    """ Repeats a function every 'sec' seconds.
    """
    def func_wrapper():
        set_interval(func, sec)
        func()
    t = threading.Timer(sec, func_wrapper)
    t.start()
    return t


# Called by master to start the worker
@app.route('/start', methods=['POST'])
def start():
    """ Called by the master when the worker is created, starts a worker.
    """
    worker_id = request.args.get('workerId')
    master_addr = request.args('masterAddr')
    hash_str = request.args.get('hash')
    hash_type = request.args.get('type')
    share = request.args('share')
    cap = request.args('cap')

    worker = Worker(worker_id, master_addr, hash_str, hash_type, share, cap)
    worker.start()

    return 'Starting worker.'

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
