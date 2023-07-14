# Traefik Proxy POC

This is a POC demonstrating how Traefik could work as the AAP Gateway Proxy

The docker-compose.yaml file in this project contains 3 services which are intended to mimic a full gateway setup using Traefik. These are:

- gateway-proxy: the traefik proxy.
- gateway-api: an apache container that mimics how the Gateway API would function serving by serving a static JSON file `/api/proxy-config` and `/api/gateway-token` which tell the proxy which services are configured and provide a Gateway API token.
- galaxy: galaxy_ng instance

How it works:

Traefik is configured to load the proxy configuration from `http://localhost:7080/api/proxy-config` (see `/traefik/traefik.yaml`). This URL mocks out how the actual Gateway API would behave in this scenario by returning a traefik json configuration. Traefik is configured to poll this API every few seconds and will reload the proxy configurations whenever they are changed. In an actual deployment, this proxy config would be generated on the fly by the AAP Gateway and contain all of the services that are configured in the Gateway Database.

The proxy configuration served by our fake Gateway API (see `/api/api/proxy-config`) does a couple of things:

- proxies `/proxy/hub/` on the proxy server to `/api/galaxy/` on the galaxy ng server
- rewrites redirect headers so that redirects map to the proxy base url and not `/api/galaxy/`
- rewrites the response body api paths so that they are returned as the proxy base url
- configures a custom middleware to perform a mocked gateway authentication

The AAP authentication middleware is the most interesting part of this. It's a custom middleware (see `/plugins-local/src/github.com/traefik/plugindemo/demo.go`) that is loaded into traefik. For each request, this middleware calls `/api/gateway-token` on the Gateway API, which returns a token that gets injected into the `Authorization` header. For this demo the header is just the basic auth header for the `admin:admin` user in GalaxyNG. For the actual application, this middleware would be configured to call the Gateway API using the request's cookies, retrieve the gateway token and then perform the rest of the Gateway Authentication flow.

## Setup

To run the demo:

```
# Start the demo
docker-compose up -d

# Wait about a minute for galaxy to start up
# Create the admin user
# NOTE: For the demo to work, the username/password MUST be admin:admin

dnewswan-mac:traefik-gateway-demo dnewswan$ docker-compose exec galaxy pulpcore-manager createsuperuser
Username: admin
Email address: a@a.com
Password: admin
Password (again): admin
The password is too similar to the username.
This password is too short. It must contain at least 9 characters.
This password is too common.
Bypass password validation and create user anyway? [y/N]: y
Superuser created successfully.

# Navigate to the `/proxy/hub/` path on the proxy.
curl http://localhost:9080/proxy/hub/ | jq
{
  "available_versions": {
    "v3": "v3/",
    "pulp-v3": "pulp/api/v3/"
  },
  "server_version": "4.7.1",
  "galaxy_ng_version": "4.7.1",
  "galaxy_ng_commit": "",
  "galaxy_importer_version": "0.4.11",
  "pulp_core_version": "3.23.8",
  "pulp_ansible_version": "0.17.3",
  "pulp_container_version": "2.14.6"
}
```

Normally this API endpoint would require authentication, but our middleware is retrieving the Gateway "token" from the Gateway API and using that to authenticate all of our requests.
