# DeFi Pilot scanner

A Ethereum transaction scanner that populates the UTU Trust Engine with data from DeFi sources

# Build

Create `private/config.yaml` and `private/protocols.json` files. You might want to use `example.config.yaml` (add your keys/URLs where indicated) and `example.protocols.json` (can be used as-is).

Then 
```make build```

# Build and push docker

Assuming AWS CLI and access is configured correctly, 

```
make docker
make docker-push

```

# Deploy the K8S pod

Delete the existing one in the `defi-portal` namespace:
```kubectl -n defi-portal delete pod/defi-portal-scanner-<suffix>```