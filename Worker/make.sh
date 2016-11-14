sudo apt-get -y install python-pip
sudo pip install -r ./requirements.txt
sudo /etc/init.d/ssh restart
sudo python itertest.py &
