helm status terraform-provider-harbor-acctest 2>/dev/null
if [ $? == 1 ]; then
  helm install terraform-provider-harbor-acctest harbor/harbor --set expose.type=nodePort,expose.tls.enabled=false
fi

echo
