sudo apt-get -y install python-pip
sudo pip install -r ./requirements.txt
sudo nohup python itertest.py > worker.log &
sleep 2
echo \n
