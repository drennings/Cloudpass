export PEM_PATH="/home/roberto/.ssh/cc_assignment0.pem"
export PEM_NAME="cc_assignment0"
export SEC_GROUP="sg-e26fb89b"
export AWS_REGION="us-west-2"
export INST_TYPE="t2.micro"
export IMG_ID="ami-d732f0b7"
cd src/main
go build -o ../../deploy/main
cd ../..
sudo -E ./deploy/main
