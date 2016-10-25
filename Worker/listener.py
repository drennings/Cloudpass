from flask import Flask
from flask import request
app = Flask(__name__)

# When a job is sent by the master
@app.route('/job', methods=['POST'])
def job():
	#initiate workers for a new job
	hash = request.arges.get('hash')
	hashtype = request.args.get('type')
	share = request.args('share')
	capacity = request.args('cap')
	return 'Running /job'

# When a worker is initiated by the master
@app.route('/worker')
def worker():
	# get code to be executed from github
    return 'Getting worker code from github'

# When this worker has found the solution
@app.route('/done')
def done():
	# store solution in the db
	# and communicate to master that we are done 
    return 'Running /done'

# When another worker has found the solution
@app.route('/stop')
def stop():
	# kill the current workers
    return 'Running /stop'
