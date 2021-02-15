helm status terraform-provider-harbor-acctest 2>/dev/null
if [ $? == 1 ]; then
  helm install terraform-provider-harbor-acctest harbor/harbor --set expose.type=loadBalancer,notary.enabled=false,expose.tls.enabled=false,externalURL=http://localhost:80 --version "1.5.3"
fi

echo
