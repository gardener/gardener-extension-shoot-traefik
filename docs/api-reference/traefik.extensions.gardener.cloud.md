# API Reference

## Packages
- [traefik.extensions.gardener.cloud/v1alpha2](#traefikextensionsgardenercloudv1alpha2)


## traefik.extensions.gardener.cloud/v1alpha2

Package v1alpha2 provides the v1alpha2 version of the external API types.
Unlike v1alpha1, the configuration fields live at the top level with no spec wrapper,
following the convention of other Gardener extension providerConfig APIs.



#### HTTPEntrypointType

_Underlying type:_ _string_

HTTPEntrypointType defines how the Traefik LoadBalancer handles plain HTTP (port 80).



_Appears in:_
- [TraefikConfig](#traefikconfig)

| Field | Description |
| --- | --- |
| `Enabled` | HTTPEntrypointEnabled exposes the Service port 80 and serves plain HTTP (default).<br /> |
| `Redirect` | HTTPEntrypointRedirect exposes the Service port 80 and redirects all HTTP requests to HTTPS (301).<br /> |
| `Disabled` | HTTPEntrypointDisabled does not expose the Service port 80. The container web entrypoint (:8000)<br />is kept regardless, because the /ping health probes depend on it.<br /> |


#### IngressProviderType

_Underlying type:_ _string_

IngressProviderType defines the type of Kubernetes Ingress provider to use.



_Appears in:_
- [TraefikConfig](#traefikconfig)

| Field | Description |
| --- | --- |
| `KubernetesIngress` | IngressProviderKubernetesIngress is the standard Kubernetes Ingress provider.<br /> |
| `KubernetesIngressNGINX` | IngressProviderKubernetesIngressNGINX is the NGINX-compatible Kubernetes Ingress provider.<br />This provider supports NGINX Ingress Controller annotations, making it easier to migrate<br />from NGINX Ingress Controller to Traefik.<br /> |


#### ProxyProtocolConfig



ProxyProtocolConfig configures PROXY protocol handling for the Traefik LoadBalancer.



_Appears in:_
- [TraefikConfig](#traefikconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled explicitly turns PROXY protocol on or off. When nil, PROXY protocol defaults<br />to enabled on AWS shoots and disabled on all other providers. |  |  |
| `trustedIPs` _string array_ | TrustedIPs is the list of CIDRs (or bare IPs) Traefik trusts to send a PROXY protocol<br />header on the web and websecure entrypoints. When empty, it defaults to the shoot's<br />node and pod network CIDRs, since the LoadBalancer lives inside the shoot VPC. |  |  |




