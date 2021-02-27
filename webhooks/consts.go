package webhooks

const (
	AllowedIPWhitelistAnnotation   = "openshift.astrokube.io/route-allowed-ip-whitelist"
	ForbiddenIPWhitelistAnnotation = "openshift.astrokube.io/route-forbidden-ip-whitelist"
	RequiredIPWhitelistAnnotation  = "openshift.astrokube.io/route-required-ip-whitelist"
	RouteIPWhitelistAnnotation     = "haproxy.router.openshift.io/ip_whitelist"
)
