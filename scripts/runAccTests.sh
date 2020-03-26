read -p 'HarborURL: ' urlvar
read -p 'Username: ' uservar
read -sp 'Password: ' passvar
echo ''


HARBOR_URL=$urlvar HARBOR_USERNAME=$uservar HARBOR_PASSWORD=$passvar TF_ACC=1 go test -timeout 20m $(go list ./... | grep -v 'vendor') -v
