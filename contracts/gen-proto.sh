# go to the contracts directory
cd ./contracts

# clone the googleapis repo in this project we just use the api/anootations.proto for the HTTP gateway
git clone https://github.com/googleapis/googleapis.git

# prepare directory for the pb files
mkdir -p pb

# generate the pb files
mkdir -p pb/loan
protoc -I . -I googleapis\
  --go_out ./pb/loan --go_opt paths=source_relative \
  --go-grpc_out ./pb/loan --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./pb/loan --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt generate_unbound_methods=true \
  ./loan.proto;

  mkdir -p pb/payment
  protoc -I . -I googleapis\
    --go_out ./pb/payment --go_opt paths=source_relative \
    --go-grpc_out ./pb/payment --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./pb/payment --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    ./payment.proto;

# go back to root of project
cd ./..