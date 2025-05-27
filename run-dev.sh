export GRPC_PORT="9090"
export REST_PORT="80"
export POSTGRES_HOST="localhost"
export POSTGRES_USERNAME="postgres"
export POSTGRES_PASSWORD="admin"
export POSTGRES_DATABASE="postgres"
export POSTGRES_PORT="5432"
export POSTGRES_SSLMODE="disable"
export POSTGRES_TIMEZONE="Asia/Jakarta"
export POSTGRES_MAX_CONNECTIONS="100"
export POSTGRES_MAX_IDLE_CONNECTIONS="10"
export POSTGRES_CONNECTIONS_MAX_IDLE_TIME="3600"

sh contracts/gen-proto.sh
go run .