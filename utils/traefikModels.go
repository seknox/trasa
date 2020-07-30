package utils

import (
	"net"
	"reflect"
)

////////////////////////////////////////////////////////////////////////////////////////////
//////////// Below models is borrowed from traefik. Traefik source has some issue //////////
//////////// with go mod so we have copied pasted model here. In future, if possible /////////
//////////// import these directly from traefik package  ///////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

// +k8s:deepcopy-gen=true
type Configuration struct {
	HTTP HTTPConfiguration `json:"http,omitempty" toml:"http,omitempty" yaml:"http,omitempty"`
	TCP  TCPConfiguration  `json:"tcp,omitempty" toml:"tcp,omitempty" yaml:"tcp,omitempty"`
	TLS  TLSConfiguration  `json:"tls,omitempty" toml:"tls,omitempty" yaml:"tls,omitempty"`
}

//HTTPConfiguration contains all the HTTP configuration parameters.
type HTTPConfiguration struct {
	Routers     map[string]Router     `json:"routers,omitempty" toml:"routers,omitempty" yaml:"routers,omitempty"`
	Middlewares map[string]Middleware `json:"middlewares,omitempty" toml:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Services    map[string]Service    `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
}

// type HTTPConfiguration struct {
// 	Routers     *Router     `json:"routers,omitempty" toml:"routers,omitempty" yaml:"routers,omitempty"`
// 	Middlewares *Middleware `json:"middlewares,omitempty" toml:"middlewares,omitempty" yaml:"middlewares,omitempty"`
// 	Services    *Service    `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
// }

// +k8s:deepcopy-gen=true

// Service holds a service configuration (can only be of one type at the same time).
type Service struct {
	LoadBalancer ServersLoadBalancer `json:"loadBalancer,omitempty" toml:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty"`
	// Weighted     WeightedRoundRobin  `json:"weighted,omitempty" toml:"weighted,omitempty" yaml:"weighted,omitempty" label:"-"`
	// Mirroring    Mirroring           `json:"mirroring,omitempty" toml:"mirroring,omitempty" yaml:"mirroring,omitempty" label:"-"`
}

// +k8s:deepcopy-gen=true

// Router holds the router configuration.
type Router struct {
	EntryPoints []string         `json:"entryPoints,omitempty" toml:"entryPoints,omitempty" yaml:"entryPoints,omitempty"`
	Middlewares []string         `json:"middlewares,omitempty" toml:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Service     string           `json:"service,omitempty" toml:"service,omitempty" yaml:"service,omitempty"`
	Rule        string           `json:"rule,omitempty" toml:"rule,omitempty" yaml:"rule,omitempty"`
	Priority    int              `json:"priority,omitempty" toml:"priority,omitempty,omitzero" yaml:"priority,omitempty"`
	TLS         *RouterTLSConfig `json:"tls,omitempty" toml:"tls,omitempty" yaml:"tls,omitempty" label:"allowEmpty"`
}

// +k8s:deepcopy-gen=true

// RouterTLSConfig holds the TLS configuration for a router
type RouterTLSConfig struct {
	Options      string   `json:"options,omitempty" toml:"options,omitempty" yaml:"options,omitempty"`
	CertResolver string   `json:"certResolver,omitempty" toml:"certResolver,omitempty" yaml:"certResolver,omitempty"`
	Domains      []Domain `json:"domains,omitempty" toml:"domains,omitempty" yaml:"domains,omitempty"`
}

// Mirroring holds the Mirroring configuration.
type Mirroring struct {
	Service string          `json:"service,omitempty" toml:"service,omitempty" yaml:"service,omitempty"`
	Mirrors []MirrorService `json:"mirrors,omitempty" toml:"mirrors,omitempty" yaml:"mirrors,omitempty"`
}

// +k8s:deepcopy-gen=true

// MirrorService holds the MirrorService configuration.
type MirrorService struct {
	Name    string `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`
	Percent int    `json:"percent,omitempty" toml:"percent,omitempty" yaml:"percent,omitempty"`
}

// +k8s:deepcopy-gen=true

// WeightedRoundRobin is a weighted round robin load-balancer of services.
type WeightedRoundRobin struct {
	Services []WRRService `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
	Sticky   *Sticky      `json:"sticky,omitempty" toml:"sticky,omitempty" yaml:"sticky,omitempty"`
}

// +k8s:deepcopy-gen=true

// WRRService is a reference to a service load-balanced with weighted round robin.
type WRRService struct {
	Name   string `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`
	Weight *int   `json:"weight,omitempty" toml:"weight,omitempty" yaml:"weight,omitempty"`
}

// SetDefaults Default values for a WRRService.
func (w *WRRService) SetDefaults() {
	defaultWeight := 1
	w.Weight = &defaultWeight
}

// +k8s:deepcopy-gen=true

// Sticky holds the sticky configuration.
type Sticky struct {
	Cookie *Cookie `json:"cookie,omitempty" toml:"cookie,omitempty" yaml:"cookie,omitempty"`
}

// +k8s:deepcopy-gen=true

// Cookie holds the sticky configuration based on cookie.
type Cookie struct {
	Name     string `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`
	Secure   bool   `json:"secure,omitempty" toml:"secure,omitempty" yaml:"secure,omitempty"`
	HTTPOnly bool   `json:"httpOnly,omitempty" toml:"httpOnly,omitempty" yaml:"httpOnly,omitempty"`
}

// +k8s:deepcopy-gen=true

// ServersLoadBalancer holds the ServersLoadBalancer configuration.
type ServersLoadBalancer struct {
	Sticky             *Sticky             `json:"sticky,omitempty" toml:"sticky,omitempty" yaml:"sticky,omitempty" label:"allowEmpty"`
	Servers            []Server            `json:"servers,omitempty" toml:"servers,omitempty" yaml:"servers,omitempty" label-slice-as-struct:"server"`
	HealthCheck        *HealthCheck        `json:"healthCheck,omitempty" toml:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	PassHostHeader     *bool               `json:"passHostHeader" toml:"passHostHeader" yaml:"passHostHeader"`
	ResponseForwarding *ResponseForwarding `json:"responseForwarding,omitempty" toml:"responseForwarding,omitempty" yaml:"responseForwarding,omitempty"`
}

// Mergeable tells if the given service is mergeable.
func (l *ServersLoadBalancer) Mergeable(loadBalancer *ServersLoadBalancer) bool {
	savedServers := l.Servers
	defer func() {
		l.Servers = savedServers
	}()
	l.Servers = nil

	savedServersLB := loadBalancer.Servers
	defer func() {
		loadBalancer.Servers = savedServersLB
	}()
	loadBalancer.Servers = nil

	return reflect.DeepEqual(l, loadBalancer)
}

// SetDefaults Default values for a ServersLoadBalancer.
func (l *ServersLoadBalancer) SetDefaults() {
	defaultPassHostHeader := true
	l.PassHostHeader = &defaultPassHostHeader
}

// +k8s:deepcopy-gen=true

// ResponseForwarding holds configuration for the forward of the response.
type ResponseForwarding struct {
	FlushInterval string `json:"flushInterval,omitempty" toml:"flushInterval,omitempty" yaml:"flushInterval,omitempty"`
}

// +k8s:deepcopy-gen=true

// Server holds the server configuration.
type Server struct {
	URL    string `json:"url,omitempty" toml:"url,omitempty" yaml:"url,omitempty" label:"-"`
	Scheme string `toml:"-" json:"-" yaml:"-"`
	Port   string `toml:"-" json:"-" yaml:"-"`
}

// SetDefaults Default values for a Server.
func (s *Server) SetDefaults() {
	s.Scheme = "http"
}

// +k8s:deepcopy-gen=true

// HealthCheck holds the HealthCheck configuration.
type HealthCheck struct {
	Scheme string `json:"scheme,omitempty" toml:"scheme,omitempty" yaml:"scheme,omitempty"`
	Path   string `json:"path,omitempty" toml:"path,omitempty" yaml:"path,omitempty"`
	Port   int    `json:"port,omitempty" toml:"port,omitempty,omitzero" yaml:"port,omitempty"`
	// FIXME change string to types.Duration
	Interval string `json:"interval,omitempty" toml:"interval,omitempty" yaml:"interval,omitempty"`
	// FIXME change string to types.Duration
	Timeout  string            `json:"timeout,omitempty" toml:"timeout,omitempty" yaml:"timeout,omitempty"`
	Hostname string            `json:"hostname,omitempty" toml:"hostname,omitempty" yaml:"hostname,omitempty"`
	Headers  map[string]string `json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty"`
}

type Middleware struct {
	AddPrefix         *AddPrefix         `json:"addPrefix,omitempty" toml:"addPrefix,omitempty" yaml:"addPrefix,omitempty"`
	StripPrefix       *StripPrefix       `json:"stripPrefix,omitempty" toml:"stripPrefix,omitempty" yaml:"stripPrefix,omitempty"`
	StripPrefixRegex  *StripPrefixRegex  `json:"stripPrefixRegex,omitempty" toml:"stripPrefixRegex,omitempty" yaml:"stripPrefixRegex,omitempty"`
	ReplacePath       *ReplacePath       `json:"replacePath,omitempty" toml:"replacePath,omitempty" yaml:"replacePath,omitempty"`
	ReplacePathRegex  *ReplacePathRegex  `json:"replacePathRegex,omitempty" toml:"replacePathRegex,omitempty" yaml:"replacePathRegex,omitempty"`
	Chain             *Chain             `json:"chain,omitempty" toml:"chain,omitempty" yaml:"chain,omitempty"`
	IPWhiteList       *IPWhiteList       `json:"ipWhiteList,omitempty" toml:"ipWhiteList,omitempty" yaml:"ipWhiteList,omitempty"`
	Headers           *Headers           `json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty"`
	Errors            *ErrorPage         `json:"errors,omitempty" toml:"errors,omitempty" yaml:"errors,omitempty"`
	RateLimit         *RateLimit         `json:"rateLimit,omitempty" toml:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	RedirectRegex     *RedirectRegex     `json:"redirectRegex,omitempty" toml:"redirectRegex,omitempty" yaml:"redirectRegex,omitempty"`
	RedirectScheme    *RedirectScheme    `json:"redirectScheme,omitempty" toml:"redirectScheme,omitempty" yaml:"redirectScheme,omitempty"`
	BasicAuth         *BasicAuth         `json:"basicAuth,omitempty" toml:"basicAuth,omitempty" yaml:"basicAuth,omitempty"`
	DigestAuth        *DigestAuth        `json:"digestAuth,omitempty" toml:"digestAuth,omitempty" yaml:"digestAuth,omitempty"`
	ForwardAuth       *ForwardAuth       `json:"forwardAuth,omitempty" toml:"forwardAuth,omitempty" yaml:"forwardAuth,omitempty"`
	InFlightReq       *InFlightReq       `json:"inFlightReq,omitempty" toml:"inFlightReq,omitempty" yaml:"inFlightReq,omitempty"`
	Buffering         *Buffering         `json:"buffering,omitempty" toml:"buffering,omitempty" yaml:"buffering,omitempty"`
	CircuitBreaker    *CircuitBreaker    `json:"circuitBreaker,omitempty" toml:"circuitBreaker,omitempty" yaml:"circuitBreaker,omitempty"`
	Compress          *Compress          `json:"compress,omitempty" toml:"compress,omitempty" yaml:"compress,omitempty" label:"allowEmpty"`
	PassTLSClientCert *PassTLSClientCert `json:"passTLSClientCert,omitempty" toml:"passTLSClientCert,omitempty" yaml:"passTLSClientCert,omitempty"`
	Retry             *Retry             `json:"retry,omitempty" toml:"retry,omitempty" yaml:"retry,omitempty"`
	Trasa             Trasa              `json:"trasa,omitempty" toml:"trasa,omitempty" yaml:"trasa,omitempty"`
}

// TRASA
type Trasa struct {
	TrasaEmtptyStruct string `json:"trasaEmtptyStruct" toml:"trasaEmtptyStruct" yaml:"trasaEmtptyStruct"`
}

// Domain holds a domain name with SANs.
type Domain struct {
	Main string   `description:"Default subject name." json:"main,omitempty" toml:"main,omitempty" yaml:"main,omitempty"`
	SANs []string `description:"Subject alternative names." json:"sans,omitempty" toml:"sans,omitempty" yaml:"sans,omitempty"`
}

// AddPrefix holds the AddPrefix configuration.
type AddPrefix struct {
	Prefix string `json:"prefix,omitempty" toml:"prefix,omitempty" yaml:"prefix,omitempty"`
}

// +k8s:deepcopy-gen=true

// Auth holds the authentication configuration (BASIC, DIGEST, users).
type Auth struct {
	Basic   *BasicAuth   `json:"basic,omitempty" toml:"basic,omitempty" yaml:"basic,omitempty" export:"true"`
	Digest  *DigestAuth  `json:"digest,omitempty" toml:"digest,omitempty" yaml:"digest,omitempty" export:"true"`
	Forward *ForwardAuth `json:"forward,omitempty" toml:"forward,omitempty" yaml:"forward,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// BasicAuth holds the HTTP basic authentication configuration.
type BasicAuth struct {
	Users        Users  `json:"users,omitempty" toml:"users,omitempty" yaml:"users,omitempty"`
	UsersFile    string `json:"usersFile,omitempty" toml:"usersFile,omitempty" yaml:"usersFile,omitempty"`
	Realm        string `json:"realm,omitempty" toml:"realm,omitempty" yaml:"realm,omitempty"`
	RemoveHeader bool   `json:"removeHeader,omitempty" toml:"removeHeader,omitempty" yaml:"removeHeader,omitempty"`
	HeaderField  string `json:"headerField,omitempty" toml:"headerField,omitempty" yaml:"headerField,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// Buffering holds the request/response buffering configuration.
type Buffering struct {
	MaxRequestBodyBytes  int64  `json:"maxRequestBodyBytes,omitempty" toml:"maxRequestBodyBytes,omitempty" yaml:"maxRequestBodyBytes,omitempty"`
	MemRequestBodyBytes  int64  `json:"memRequestBodyBytes,omitempty" toml:"memRequestBodyBytes,omitempty" yaml:"memRequestBodyBytes,omitempty"`
	MaxResponseBodyBytes int64  `json:"maxResponseBodyBytes,omitempty" toml:"maxResponseBodyBytes,omitempty" yaml:"maxResponseBodyBytes,omitempty"`
	MemResponseBodyBytes int64  `json:"memResponseBodyBytes,omitempty" toml:"memResponseBodyBytes,omitempty" yaml:"memResponseBodyBytes,omitempty"`
	RetryExpression      string `json:"retryExpression,omitempty" toml:"retryExpression,omitempty" yaml:"retryExpression,omitempty"`
}

// +k8s:deepcopy-gen=true

// Chain holds a chain of middlewares
type Chain struct {
	Middlewares []string `json:"middlewares,omitempty" toml:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

// +k8s:deepcopy-gen=true

// CircuitBreaker holds the circuit breaker configuration.
type CircuitBreaker struct {
	Expression string `json:"expression,omitempty" toml:"expression,omitempty" yaml:"expression,omitempty"`
}

// +k8s:deepcopy-gen=true

// Compress holds the compress configuration.
type Compress struct {
	ExcludedContentTypes []string `json:"excludedContentTypes,omitempty" toml:"excludedContentTypes,omitempty" yaml:"excludedContentTypes,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// DigestAuth holds the Digest HTTP authentication configuration.
type DigestAuth struct {
	Users        Users  `json:"users,omitempty" toml:"users,omitempty" yaml:"users,omitempty"`
	UsersFile    string `json:"usersFile,omitempty" toml:"usersFile,omitempty" yaml:"usersFile,omitempty"`
	RemoveHeader bool   `json:"removeHeader,omitempty" toml:"removeHeader,omitempty" yaml:"removeHeader,omitempty"`
	Realm        string `json:"realm,omitempty" toml:"realm,omitempty" yaml:"realm,omitempty"`
	HeaderField  string `json:"headerField,omitempty" toml:"headerField,omitempty" yaml:"headerField,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// ErrorPage holds the custom error page configuration.
type ErrorPage struct {
	Status  []string `json:"status,omitempty" toml:"status,omitempty" yaml:"status,omitempty"`
	Service string   `json:"service,omitempty" toml:"service,omitempty" yaml:"service,omitempty"`
	Query   string   `json:"query,omitempty" toml:"query,omitempty" yaml:"query,omitempty"`
}

// +k8s:deepcopy-gen=true

// ForwardAuth holds the http forward authentication configuration.
type ForwardAuth struct {
	Address             string     `json:"address,omitempty" toml:"address,omitempty" yaml:"address,omitempty"`
	TLS                 *ClientTLS `json:"tls,omitempty" toml:"tls,omitempty" yaml:"tls,omitempty"`
	TrustForwardHeader  bool       `json:"trustForwardHeader,omitempty" toml:"trustForwardHeader,omitempty" yaml:"trustForwardHeader,omitempty" export:"true"`
	AuthResponseHeaders []string   `json:"authResponseHeaders,omitempty" toml:"authResponseHeaders,omitempty" yaml:"authResponseHeaders,omitempty"`
}

// +k8s:deepcopy-gen=true

// Headers holds the custom header configuration.
type Headers struct {
	CustomRequestHeaders  map[string]string `json:"customRequestHeaders,omitempty" toml:"customRequestHeaders,omitempty" yaml:"customRequestHeaders,omitempty"`
	CustomResponseHeaders map[string]string `json:"customResponseHeaders,omitempty" toml:"customResponseHeaders,omitempty" yaml:"customResponseHeaders,omitempty"`

	// AccessControlAllowCredentials is only valid if true. false is ignored.
	AccessControlAllowCredentials bool `json:"accessControlAllowCredentials,omitempty" toml:"accessControlAllowCredentials,omitempty" yaml:"accessControlAllowCredentials,omitempty"`
	// AccessControlAllowHeaders must be used in response to a preflight request with Access-Control-Request-Headers set.
	AccessControlAllowHeaders []string `json:"accessControlAllowHeaders,omitempty" toml:"accessControlAllowHeaders,omitempty" yaml:"accessControlAllowHeaders,omitempty"`
	// AccessControlAllowMethods must be used in response to a preflight request with Access-Control-Request-Method set.
	AccessControlAllowMethods []string `json:"accessControlAllowMethods,omitempty" toml:"accessControlAllowMethods,omitempty" yaml:"accessControlAllowMethods,omitempty"`
	// AccessControlAllowOrigin Can be "origin-list-or-null" or "*". From (https://www.w3.org/TR/cors/#access-control-allow-origin-response-header)
	AccessControlAllowOrigin string `json:"accessControlAllowOrigin,omitempty" toml:"accessControlAllowOrigin,omitempty" yaml:"accessControlAllowOrigin,omitempty"`
	// AccessControlExposeHeaders sets valid headers for the response.
	AccessControlExposeHeaders []string `json:"accessControlExposeHeaders,omitempty" toml:"accessControlExposeHeaders,omitempty" yaml:"accessControlExposeHeaders,omitempty"`
	// AccessControlMaxAge sets the time that a preflight request may be cached.
	AccessControlMaxAge int64 `json:"accessControlMaxAge,omitempty" toml:"accessControlMaxAge,omitempty" yaml:"accessControlMaxAge,omitempty"`
	// AddVaryHeader controls if the Vary header is automatically added/updated when the AccessControlAllowOrigin is set.
	AddVaryHeader bool `json:"addVaryHeader,omitempty" toml:"addVaryHeader,omitempty" yaml:"addVaryHeader,omitempty"`

	AllowedHosts            []string          `json:"allowedHosts,omitempty" toml:"allowedHosts,omitempty" yaml:"allowedHosts,omitempty"`
	HostsProxyHeaders       []string          `json:"hostsProxyHeaders,omitempty" toml:"hostsProxyHeaders,omitempty" yaml:"hostsProxyHeaders,omitempty"`
	SSLRedirect             bool              `json:"sslRedirect,omitempty" toml:"sslRedirect,omitempty" yaml:"sslRedirect,omitempty"`
	SSLTemporaryRedirect    bool              `json:"sslTemporaryRedirect,omitempty" toml:"sslTemporaryRedirect,omitempty" yaml:"sslTemporaryRedirect,omitempty"`
	SSLHost                 string            `json:"sslHost,omitempty" toml:"sslHost,omitempty" yaml:"sslHost,omitempty"`
	SSLProxyHeaders         map[string]string `json:"sslProxyHeaders,omitempty" toml:"sslProxyHeaders,omitempty" yaml:"sslProxyHeaders,omitempty"`
	SSLForceHost            bool              `json:"sslForceHost,omitempty" toml:"sslForceHost,omitempty" yaml:"sslForceHost,omitempty"`
	STSSeconds              int64             `json:"stsSeconds,omitempty" toml:"stsSeconds,omitempty" yaml:"stsSeconds,omitempty"`
	STSIncludeSubdomains    bool              `json:"stsIncludeSubdomains,omitempty" toml:"stsIncludeSubdomains,omitempty" yaml:"stsIncludeSubdomains,omitempty"`
	STSPreload              bool              `json:"stsPreload,omitempty" toml:"stsPreload,omitempty" yaml:"stsPreload,omitempty"`
	ForceSTSHeader          bool              `json:"forceSTSHeader,omitempty" toml:"forceSTSHeader,omitempty" yaml:"forceSTSHeader,omitempty"`
	FrameDeny               bool              `json:"frameDeny,omitempty" toml:"frameDeny,omitempty" yaml:"frameDeny,omitempty"`
	CustomFrameOptionsValue string            `json:"customFrameOptionsValue,omitempty" toml:"customFrameOptionsValue,omitempty" yaml:"customFrameOptionsValue,omitempty"`
	ContentTypeNosniff      bool              `json:"contentTypeNosniff,omitempty" toml:"contentTypeNosniff,omitempty" yaml:"contentTypeNosniff,omitempty"`
	BrowserXSSFilter        bool              `json:"browserXssFilter,omitempty" toml:"browserXssFilter,omitempty" yaml:"browserXssFilter,omitempty"`
	CustomBrowserXSSValue   string            `json:"customBrowserXSSValue,omitempty" toml:"customBrowserXSSValue,omitempty" yaml:"customBrowserXSSValue,omitempty"`
	ContentSecurityPolicy   string            `json:"contentSecurityPolicy,omitempty" toml:"contentSecurityPolicy,omitempty" yaml:"contentSecurityPolicy,omitempty"`
	PublicKey               string            `json:"publicKey,omitempty" toml:"publicKey,omitempty" yaml:"publicKey,omitempty"`
	ReferrerPolicy          string            `json:"referrerPolicy,omitempty" toml:"referrerPolicy,omitempty" yaml:"referrerPolicy,omitempty"`
	FeaturePolicy           string            `json:"featurePolicy,omitempty" toml:"featurePolicy,omitempty" yaml:"featurePolicy,omitempty"`
	IsDevelopment           bool              `json:"isDevelopment,omitempty" toml:"isDevelopment,omitempty" yaml:"isDevelopment,omitempty"`
}

// +k8s:deepcopy-gen=true

// IPStrategy holds the ip strategy configuration.
type IPStrategy struct {
	Depth       int      `json:"depth,omitempty" toml:"depth,omitempty" yaml:"depth,omitempty" export:"true"`
	ExcludedIPs []string `json:"excludedIPs,omitempty" toml:"excludedIPs,omitempty" yaml:"excludedIPs,omitempty"`
	// TODO(mpl): I think we should make RemoteAddr an explicit field. For one thing, it would yield better documentation.
}

// +k8s:deepcopy-gen=true

// IPWhiteList holds the ip white list configuration.
type IPWhiteList struct {
	SourceRange []string    `json:"sourceRange,omitempty" toml:"sourceRange,omitempty" yaml:"sourceRange,omitempty"`
	IPStrategy  *IPStrategy `json:"ipStrategy,omitempty" toml:"ipStrategy,omitempty" yaml:"ipStrategy,omitempty"  label:"allowEmpty"`
}

// +k8s:deepcopy-gen=true

// InFlightReq limits the number of requests being processed and served concurrently.
type InFlightReq struct {
	Amount          int64            `json:"amount,omitempty" toml:"amount,omitempty" yaml:"amount,omitempty"`
	SourceCriterion *SourceCriterion `json:"sourceCriterion,omitempty" toml:"sourceCriterion,omitempty" yaml:"sourceCriterion,omitempty"`
}

// SetDefaults Default values for a InFlightReq.
func (i *InFlightReq) SetDefaults() {
	i.SourceCriterion = &SourceCriterion{
		RequestHost: true,
	}
}

// +k8s:deepcopy-gen=true

// PassTLSClientCert holds the TLS client cert headers configuration.
type PassTLSClientCert struct {
	PEM  bool                      `json:"pem,omitempty" toml:"pem,omitempty" yaml:"pem,omitempty"`
	Info *TLSClientCertificateInfo `json:"info,omitempty" toml:"info,omitempty" yaml:"info,omitempty"`
}

// +k8s:deepcopy-gen=true

// SourceCriterion defines what criterion is used to group requests as originating from a common source.
// The precedence order is IPStrategy, then RequestHeaderName.
// If none are set, the default is to use the request's remote address field.
type SourceCriterion struct {
	IPStrategy        *IPStrategy `json:"ipStrategy" toml:"ipStrategy, omitempty"`
	RequestHeaderName string      `json:"requestHeaderName,omitempty" toml:"requestHeaderName,omitempty" yaml:"requestHeaderName,omitempty"`
	RequestHost       bool        `json:"requestHost,omitempty" toml:"requestHost,omitempty" yaml:"requestHost,omitempty"`
}

// +k8s:deepcopy-gen=true

// RateLimit holds the rate limiting configuration for a given router.
type RateLimit struct {
	// Average is the maximum rate, in requests/s, allowed for the given source.
	// It defaults to 0, which means no rate limiting.
	Average int64 `json:"average,omitempty" toml:"average,omitempty" yaml:"average,omitempty"`
	// Burst is the maximum number of requests allowed to arrive in the same arbitrarily small period of time.
	// It defaults to 1.
	Burst           int64            `json:"burst,omitempty" toml:"burst,omitempty" yaml:"burst,omitempty"`
	SourceCriterion *SourceCriterion `json:"sourceCriterion,omitempty" toml:"sourceCriterion,omitempty" yaml:"sourceCriterion,omitempty"`
}

// SetDefaults sets the default values on a RateLimit.
func (r *RateLimit) SetDefaults() {
	r.Burst = 1
	r.SourceCriterion = &SourceCriterion{
		IPStrategy: &IPStrategy{},
	}
}

// +k8s:deepcopy-gen=true

// RedirectRegex holds the redirection configuration.
type RedirectRegex struct {
	Regex       string `json:"regex,omitempty" toml:"regex,omitempty" yaml:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty" toml:"replacement,omitempty" yaml:"replacement,omitempty"`
	Permanent   bool   `json:"permanent,omitempty" toml:"permanent,omitempty" yaml:"permanent,omitempty"`
}

// +k8s:deepcopy-gen=true

// RedirectScheme holds the scheme redirection configuration.
type RedirectScheme struct {
	Scheme    string `json:"scheme,omitempty" toml:"scheme,omitempty" yaml:"scheme,omitempty"`
	Port      string `json:"port,omitempty" toml:"port,omitempty" yaml:"port,omitempty"`
	Permanent bool   `json:"permanent,omitempty" toml:"permanent,omitempty" yaml:"permanent,omitempty"`
}

// +k8s:deepcopy-gen=true

// ReplacePath holds the ReplacePath configuration.
type ReplacePath struct {
	Path string `json:"path,omitempty" toml:"path,omitempty" yaml:"path,omitempty"`
}

// +k8s:deepcopy-gen=true

// ReplacePathRegex holds the ReplacePathRegex configuration.
type ReplacePathRegex struct {
	Regex       string `json:"regex,omitempty" toml:"regex,omitempty" yaml:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty" toml:"replacement,omitempty" yaml:"replacement,omitempty"`
}

// +k8s:deepcopy-gen=true

// Retry holds the retry configuration.
type Retry struct {
	Attempts int `json:"attempts,omitempty" toml:"attempts,omitempty" yaml:"attempts,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// StripPrefix holds the StripPrefix configuration.
type StripPrefix struct {
	Prefixes   []string `json:"prefixes,omitempty" toml:"prefixes,omitempty" yaml:"prefixes,omitempty"`
	ForceSlash bool     `json:"forceSlash,omitempty" toml:"forceSlash,omitempty" yaml:"forceSlash,omitempty"` // Deprecated
}

// SetDefaults Default values for a StripPrefix.
func (s *StripPrefix) SetDefaults() {
	s.ForceSlash = true
}

// +k8s:deepcopy-gen=true

// StripPrefixRegex holds the StripPrefixRegex configuration.
type StripPrefixRegex struct {
	Regex []string `json:"regex,omitempty" toml:"regex,omitempty" yaml:"regex,omitempty"`
}

// +k8s:deepcopy-gen=true

// TLSClientCertificateInfo holds the client TLS certificate info configuration.
type TLSClientCertificateInfo struct {
	NotAfter  bool                        `json:"notAfter,omitempty" toml:"notAfter,omitempty" yaml:"notAfter,omitempty"`
	NotBefore bool                        `json:"notBefore,omitempty" toml:"notBefore,omitempty" yaml:"notBefore,omitempty"`
	Sans      bool                        `json:"sans,omitempty" toml:"sans,omitempty" yaml:"sans,omitempty"`
	Subject   *TLSCLientCertificateDNInfo `json:"subject,omitempty" toml:"subject,omitempty" yaml:"subject,omitempty"`
	Issuer    *TLSCLientCertificateDNInfo `json:"issuer,omitempty" toml:"issuer,omitempty" yaml:"issuer,omitempty"`
}

// +k8s:deepcopy-gen=true

// TLSCLientCertificateDNInfo holds the client TLS certificate distinguished name info configuration
// cf https://tools.ietf.org/html/rfc3739
type TLSCLientCertificateDNInfo struct {
	Country         bool `json:"country,omitempty" toml:"country,omitempty" yaml:"country,omitempty"`
	Province        bool `json:"province,omitempty" toml:"province,omitempty" yaml:"province,omitempty"`
	Locality        bool `json:"locality,omitempty" toml:"locality,omitempty" yaml:"locality,omitempty"`
	Organization    bool `json:"organization,omitempty" toml:"organization,omitempty" yaml:"organization,omitempty"`
	CommonName      bool `json:"commonName,omitempty" toml:"commonName,omitempty" yaml:"commonName,omitempty"`
	SerialNumber    bool `json:"serialNumber,omitempty" toml:"serialNumber,omitempty" yaml:"serialNumber,omitempty"`
	DomainComponent bool `json:"domainComponent,omitempty" toml:"domainComponent,omitempty" yaml:"domainComponent,omitempty"`
}

// +k8s:deepcopy-gen=true

// Users holds a list of users
type Users []string

// +k8s:deepcopy-gen=true

// ClientTLS holds the TLS specific configurations as client
// CA, Cert and Key can be either path or file contents.
type ClientTLS struct {
	CA                 string `json:"ca,omitempty" toml:"ca,omitempty" yaml:"ca,omitempty"`
	CAOptional         bool   `json:"caOptional,omitempty" toml:"caOptional,omitempty" yaml:"caOptional,omitempty"`
	Cert               string `json:"cert,omitempty" toml:"cert,omitempty" yaml:"cert,omitempty"`
	Key                string `json:"key,omitempty" toml:"key,omitempty" yaml:"key,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify,omitempty" toml:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty"`
}

// Checker allows to check that addresses are in a trusted IPs
type Checker struct {
	authorizedIPs    []*net.IP
	authorizedIPsNet []*net.IPNet
}

// DepthStrategy a strategy based on the depth inside the X-Forwarded-For from right to left
type DepthStrategy struct {
	Depth int
}

/////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
/////////////			TCP Configuration 		////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// TCPConfiguration contains all the TCP configuration parameters.
type TCPConfiguration struct {
	Routers  map[string]TCPRouter  `json:"routers,omitempty" toml:"routers,omitempty" yaml:"routers,omitempty"`
	Services map[string]TCPService `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
}

// +k8s:deepcopy-gen=true

// TCPService holds a tcp service configuration (can only be of one type at the same time).
type TCPService struct {
	LoadBalancer *TCPServersLoadBalancer `json:"loadBalancer,omitempty" toml:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty"`
	Weighted     *TCPWeightedRoundRobin  `json:"weighted,omitempty" toml:"weighted,omitempty" yaml:"weighted,omitempty" label:"-"`
}

// +k8s:deepcopy-gen=true

// TCPWeightedRoundRobin is a weighted round robin tcp load-balancer of services.
type TCPWeightedRoundRobin struct {
	Services []TCPWRRService `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
}

// +k8s:deepcopy-gen=true

// TCPWRRService is a reference to a tcp service load-balanced with weighted round robin.
type TCPWRRService struct {
	Name   string `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`
	Weight *int   `json:"weight,omitempty" toml:"weight,omitempty" yaml:"weight,omitempty"`
}

// SetDefaults Default values for a TCPWRRService.
func (w *TCPWRRService) SetDefaults() {
	defaultWeight := 1
	w.Weight = &defaultWeight
}

// +k8s:deepcopy-gen=true

// TCPRouter holds the router configuration.
type TCPRouter struct {
	EntryPoints []string            `json:"entryPoints,omitempty" toml:"entryPoints,omitempty" yaml:"entryPoints,omitempty"`
	Service     string              `json:"service,omitempty" toml:"service,omitempty" yaml:"service,omitempty"`
	Rule        string              `json:"rule,omitempty" toml:"rule,omitempty" yaml:"rule,omitempty"`
	TLS         *RouterTCPTLSConfig `json:"tls,omitempty" toml:"tls,omitempty" yaml:"tls,omitempty" label:"allowEmpty"`
}

// +k8s:deepcopy-gen=true

// RouterTCPTLSConfig holds the TLS configuration for a router
type RouterTCPTLSConfig struct {
	Passthrough  bool     `json:"passthrough" toml:"passthrough" yaml:"passthrough"`
	Options      string   `json:"options,omitempty" toml:"options,omitempty" yaml:"options,omitempty"`
	CertResolver string   `json:"certResolver,omitempty" toml:"certResolver,omitempty" yaml:"certResolver,omitempty"`
	Domains      []Domain `json:"domains,omitempty" toml:"domains,omitempty" yaml:"domains,omitempty"`
}

// +k8s:deepcopy-gen=true

// TCPServersLoadBalancer holds the LoadBalancerService configuration.
type TCPServersLoadBalancer struct {
	// TerminationDelay, corresponds to the deadline that the proxy sets, after one
	// of its connected peers indicates it has closed the writing capability of its
	// connection, to close the reading capability as well, hence fully terminating the
	// connection. It is a duration in milliseconds, defaulting to 100. A negative value
	// means an infinite deadline (i.e. the reading capability is never closed).
	TerminationDelay *int        `json:"terminationDelay,omitempty" toml:"terminationDelay,omitempty" yaml:"terminationDelay,omitempty"`
	Servers          []TCPServer `json:"servers,omitempty" toml:"servers,omitempty" yaml:"servers,omitempty" label-slice-as-struct:"server"`
}

// SetDefaults Default values for a TCPServersLoadBalancer
func (l *TCPServersLoadBalancer) SetDefaults() {
	defaultTerminationDelay := 100 // in milliseconds
	l.TerminationDelay = &defaultTerminationDelay
}

// Mergeable tells if the given service is mergeable.
func (l *TCPServersLoadBalancer) Mergeable(loadBalancer *TCPServersLoadBalancer) bool {
	savedServers := l.Servers
	defer func() {
		l.Servers = savedServers
	}()
	l.Servers = nil

	savedServersLB := loadBalancer.Servers
	defer func() {
		loadBalancer.Servers = savedServersLB
	}()
	loadBalancer.Servers = nil

	return reflect.DeepEqual(l, loadBalancer)
}

// +k8s:deepcopy-gen=true

// TCPServer holds a TCP Server configuration
type TCPServer struct {
	Address string `json:"address,omitempty" toml:"address,omitempty" yaml:"address,omitempty" label:"-"`
	Port    string `toml:"-" json:"-" yaml:"-"`
}

/////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
/////////////			TLS Configuration 		////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// TLSConfiguration contains all the configuration parameters of a TLS connection.
type TLSConfiguration struct {
	Certificates []CertAndStores    `json:"certificates,omitempty"  toml:"certificates,omitempty" yaml:"certificates,omitempty" label:"-"`
	Options      map[string]Options `json:"options,omitempty" toml:"options,omitempty" yaml:"options,omitempty"`
	Stores       map[string]Store   `json:"stores,omitempty" toml:"stores,omitempty" yaml:"stores,omitempty"`
}

// ClientAuth defines the parameters of the client authentication part of the TLS connection, if any.
type ClientAuth struct {
	CAFiles []FileOrContent `json:"caFiles,omitempty" toml:"caFiles,omitempty" yaml:"caFiles,omitempty"`
	// ClientAuthType defines the client authentication type to apply.
	// The available values are: "NoClientCert", "RequestClientCert", "VerifyClientCertIfGiven" and "RequireAndVerifyClientCert".
	ClientAuthType string `json:"clientAuthType,omitempty" toml:"clientAuthType,omitempty" yaml:"clientAuthType,omitempty"`
}

// FileOrContent hold a file path or content
type FileOrContent string

// +k8s:deepcopy-gen=true

// Options configures TLS for an entry point
type Options struct {
	MinVersion       string     `json:"minVersion,omitempty" toml:"minVersion,omitempty" yaml:"minVersion,omitempty" export:"true"`
	MaxVersion       string     `json:"maxVersion,omitempty" toml:"maxVersion,omitempty" yaml:"maxVersion,omitempty" export:"true"`
	CipherSuites     []string   `json:"cipherSuites,omitempty" toml:"cipherSuites,omitempty" yaml:"cipherSuites,omitempty"`
	CurvePreferences []string   `json:"curvePreferences,omitempty" toml:"curvePreferences,omitempty" yaml:"curvePreferences,omitempty"`
	ClientAuth       ClientAuth `json:"clientAuth,omitempty" toml:"clientAuth,omitempty" yaml:"clientAuth,omitempty"`
	SniStrict        bool       `json:"sniStrict,omitempty" toml:"sniStrict,omitempty" yaml:"sniStrict,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// Store holds the options for a given Store
type Store struct {
	DefaultCertificate *Certificate `json:"defaultCertificate,omitempty" toml:"defaultCertificate,omitempty" yaml:"defaultCertificate,omitempty"`
}

// CertAndStores allows mapping a TLS certificate to a list of entry points.
type CertAndStores struct {
	Certificate `yaml:",inline"`
	Stores      []string `json:"stores,omitempty" toml:"stores,omitempty" yaml:"stores,omitempty"`
}

// Certificate holds a SSL cert/key pair
// Certs and Key could be either a file path, or the file content itself
type Certificate struct {
	CertFile FileOrContent `json:"certFile,omitempty" toml:"certFile,omitempty" yaml:"certFile,omitempty"`
	KeyFile  FileOrContent `json:"keyFile,omitempty" toml:"keyFile,omitempty" yaml:"keyFile,omitempty"`
}
