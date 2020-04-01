if [ -z "$HARBOR_URL" ]; then
  read -p 'HarborURL: ' HARBOR_URL
fi
if [ -z "$HARBOR_USERNAME" ]; then
  read -p 'Username: ' HARBOR_USERNAME
fi
if [ -z "$HARBOR_PASSWORD" ]; then
  read -sp 'Password: ' HARBOR_PASSWORD
fi
echo ''

export HARBOR_URL
export HARBOR_USERNAME
export HARBOR_PASSWORD

TF_ACC=1 go test -timeout 20m $(go list ./... | grep -v 'vendor') -v
