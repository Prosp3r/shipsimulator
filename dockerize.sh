az acr login --name raaedge
export VERSION=1.0.0
docker build --no-cache -f ./nmea/cmd/nmeagenerator/Dockerfile -t raaedge.azurecr.io/shipsimulator-nmea:$VERSION .
docker push raaedge.azurecr.io/shipsimulator-nmea:$VERSION
