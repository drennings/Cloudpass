
alphabet = 'abcdefghijklmnopqrstuvwxyz'
complete_list = []
for current in xrange(10):
    a = [i for i in alphabet]
    for y in xrange(current):
        a = [x+i for i in alphabet for x in a]
    complete_list = complete_list+a
	print a


