# OpenShift Multitenancy Controller

OpenShift Multitenancy Controller improves the multitenancy features for OpenShift clusters.

## Compatibility

| Distribution | Min version | Max version |
| --- | --- | --- |
| OpenShift | 4.1 | 4.7 |

## Usage

### OpenShift Routes

Those annotations can be applied to Namespaces in order to restrict the usage of Routes:

| Annotation | Description | Type | Default value |
| --- | --- | --- | --- |
| `openshift.astrokube.io/route-allowed-ip-whitelist` | List of allowed values for the Route annotation `haproxy.router.openshift.io/ip_whitelist`. If not set, this rule will be ignored. | array(string) | null |
| `openshift.astrokube.io/route-forbidden-ip-whitelist` | List of required values for the Route annotation `haproxy.router.openshift.io/ip_whitelist`. If not set, this rule will be ignored. | array(string) | null |
| `openshift.astrokube.io/route-required-ip-whitelist` | List of forbidden values for the Route annotation `haproxy.router.openshift.io/ip_whitelist`. If not set, this rule will be ignored. | array(string) | null |

## Install

```bash
curl https://raw.githubusercontent.com/astrokube/openshift-multitenancy-controller/main/deploy/install.yaml | oc apply -f -
```

## Uninstall

```bash
curl https://raw.githubusercontent.com/astrokube/openshift-multitenancy-controller/main/deploy/install.yaml | oc delete -f -
```

## Examples

### Force all Routes to have the same IP Whitelist

We want to force all Routes inside the `example` Namespace to have a fixed IP Whitelist: 10.0.0.0/24 and 10.0.1.100

```bash
oc annotate namespace example openshift.astrokube.io/route-required-ip-whitelist=10.0.0.0/24,10.0.1.100
oc annotate namespace example openshift.astrokube.io/route-allowed-ip-whitelist=10.0.0.0/24,10.0.1.100
```

### Restrict all Routes to use only some values for IP Whitelist

We want to allow Routes inside the `example` Namespace to be allowed to use only some IP/CIDR on IP Whitelist: 10.0.0.0/24 and 10.0.1.100

```bash
oc annotate namespace example openshift.astrokube.io/route-allowed-ip-whitelist=10.0.0.0/24,10.0.1.100
```

### Forbid all Routes to use some values for IP Whitelist

We want to forbid all Routes inside the `example` Namespace to use some IP/CIDR on IP Whitelist: 10.0.0.0/24 and 10.0.1.100

```bash
oc annotate namespace example openshift.astrokube.io/route-forbidden-ip-whitelist=10.0.0.0/24,10.0.1.100
```
