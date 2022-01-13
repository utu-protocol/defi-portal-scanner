# DeFi Portal scanner

A Ethereum transaction scanner that populates the UTU Trust Engine with data from DeFi sources

# How to Use

## OCEAN Aquarius
`defi-portal-scanner ocean scanpush`
This asks OCEAN Aquarius for information about addresses, pools, builds up an internal state, and pushes the state to UTU Trust API, transforming a few things along the way for convenience's sake (of the people using the UTU Trust API).

It is very simple to run: no need for a config file. Just set the `APIKEY` and `APIURL` environment variables so that it can authenticate with Trust API. If `APIURL` is not set it will default to `https://stage-api.ututrust.com/core-api`.

for example

`APIKEY="fdsafdsafdsa" defi-portal-scanner ocean scanpush`

This should be run periodically via a cron. No more than once per hour, out of courtesy to OCEAN.

## Defi Portal Scanner
This runs a webserver. When a `POST /subscribe/<address>` comes in, it will query Etherscan and try to make sense of the answer. Then it will POST something back to the UTU Trust API.

` defi-portal-scanner listen --scan -c private/config.yaml -p private/protocols.json --http`


# Build and Deploy

## Build

Create `private/config.yaml` and `private/protocols.json` files. You might want to use `example.config.yaml` (add your keys/URLs where indicated) and `example.protocols.json` (can be used as-is).

Then 
```make build```

## Build and push docker

Assuming AWS CLI and access is configured correctly, 

```
make docker
make docker-push

```

## Deploy the K8S pod

Delete the existing one in the `defi-portal` namespace:
```kubectl -n defi-portal delete pod/defi-portal-scanner-<suffix>```