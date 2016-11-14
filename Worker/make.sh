sudo apt-get -y install python-pip
sudo pip install -r ./requirements.txt
python itertest.py &
sudo /etc/init.d/ssh restart
