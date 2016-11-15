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
length = 7
heartbeat = None

class Worker:
    def __init__(self, worker_id, master_addr,
                 hash_str, hash_type, share, capacity):
        self.worker_id = worker_id
        self.master_addr = master_addr
        self.hash_str = hash_str
        self.hash_type = hash_type
        self.share = share
        self.capacity = capacity
        self.heartbeat_interval = 10
        self.solution = None
        self.solutions_tried = 0
        self.solutions_total = 0
        for i in range(length+1):
            self.solutions_total += len(alphabet)**i

    def start(self):
        """ Start bruteforcing the password while sending heartbeats
        to the master every self.heartbeat_interval seconds to show we
        are alive.
        """
        print("Starting bruteforce for hash: " + self.hash_str)

        set_interval(self.notify_master, self.heartbeat_interval)

        # Bruteforce our share of the solution space
        self.time_start = time.time()

        for i in range(length+1):
            for p in product(alphabet, repeat=i):
                self.solutions_tried += 1

                # Only do our share of the solution space
                if (self.solutions_tried % self.capacity) != self.share:
                    continue
                sol = ''.join(p)
                hash_func = None

                # Determine hash function
                if self.hash_type == "md5":
                    hash_func = hashlib.md5()
                elif self.hash_type == "sha1":
                    hash_func = hashlib.sha1()

                sol_ascii_bytes = bytes(sol)
                hash_func.update(sol_ascii_bytes)

                # Solution not found
                if self.solutions_tried == self.solutions_total:
                    print("Did not find solution for hash " + self.hash_str)
                    self.stop()

                # Solution found
                if self.hash_str == hash_func.hexdigest():
                    print("Found solution for hash " + self.hash_str + ": " +
                          str(sol))
                    self.solution = sol
                    self.stop()
                    return

    def stop(self):
        self.time_stop = time.time()
        self.notify_master()
        heartbeat.cancel()

    def notify_master(self):
        """ Sends our state to the master, this is done every
        'heartbeat_interval' seconds as well as when the solution is found.
        """
        print('Pinging master' + str(self.solutions_tried))
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
    global heartbeat
    heartbeat = threading.Timer(sec, func_wrapper)
    heartbeat.start()


# Called by master to start the worker
@app.route('/start', methods=['POST'])
def start():
    """ Called by the master when the worker is created, starts a worker.
    """
    json = request.get_json()
    worker_id = json['workerId']
    master_addr = json['masterAddr']
    hash_str = json['hashStr']
    hash_type = json['hashType']
    share = int(json['share'])
    cap = int(json['cap'])

    worker = Worker(worker_id, master_addr, hash_str, hash_type, share, cap)

    print('Received request!')
    print(json)

    t = threading.Thread(target=worker.start)
    t.start()

    return 'Starting worker.'

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
