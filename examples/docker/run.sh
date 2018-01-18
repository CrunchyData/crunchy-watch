DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export CONTAINER_NAME=watch
$DIR/cleanup.sh

docker run \
	--privileged=true \
	--link primary:primary \
	--link replica:replica \
	-e PLATFORM="docker" \
	-e PG_PRIMARY="primary:5432" \
	-e PG_REPLICA="replica:5432" \
	-e PG_USERNAME="primaryuser" \
	-e PG_PASSWORD="password" \
	-e PG_DATABASE="postgres" \
	-e PG_HEALTHCHECK_INTERVAL="30s" \
	-e PG_FAILOVER_WAIT="10s" \
	--name=$CONTAINER_NAME \
	--hostname=$CONTAINER_NAME \
	-d crunchydata/crunchy-watch:$CCP_IMAGE_TAG

# -e PG_FAILOVER_WAIT="10s" \
