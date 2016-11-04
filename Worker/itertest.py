from itertools import product
import time
import hashlib
import sys
import math

solution = 'ef775988943825d2871e1cfa75473ec0'
alphabet = '0123456789'
length = 8

upper_bound = 0
for i in range(length): # cartesian product of lengths 0 - length
    upper_bound += math.factorial(i)

print("Calculated upper bound: " + str(upper_bound))

start = time.time()
tries = 1
for i in range(length):
    for p in product(alphabet, repeat=i):
        if (tries % (upper_bound / 100)) == 0: # every 1%
            progress = tries / upper_bound * 100
            sys.stdout.write("Processing pins of length {0}, {1}%\r".format(i, math.floor(progress)))
            sys.stdout.flush()

        sol = ''.join(p)
        md5 = hashlib.md5()
        sol_ascii_bytes = bytes(sol, 'utf-8')
        md5.update(sol_ascii_bytes)
        if solution == md5.hexdigest():
            print("Found solution: " + sol)
            print(str(int((time.time() - start)) * 1000) + "ms")
            sys.exit(0)
        tries += 1
