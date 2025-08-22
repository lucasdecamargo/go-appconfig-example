package config

// GroupNetwork is the logical group name for network-related configuration fields
const GroupNetwork = "Network"

func init() {
	// Register all network configuration fields
	Fields.Add(
		FieldNetworkProxyAll,
		FieldNetworkProxyHttp,
		FieldNetworkProxyHttps,
	)
}

// Proxy configuration fields

// FieldNetworkProxyAll defines a proxy server for all network traffic
var FieldNetworkProxyAll = &Field{
	Name:        "proxy.all",
	Group:       GroupNetwork,
	Type:        FieldTypeString,
	Default:     nil,
	Description: "Set a proxy server for all network traffic",
	ValidateTag: "url",
	Example:     "http://user:password@host:port",
}

// FieldNetworkProxyHttp defines a proxy server for HTTP traffic only
var FieldNetworkProxyHttp = &Field{
	Name:        "proxy.http",
	Group:       GroupNetwork,
	Type:        FieldTypeString,
	Default:     nil,
	Description: "Set a proxy server for HTTP traffic",
	ValidateTag: "url",
	Example:     "http://user:password@host:port",
}

// FieldNetworkProxyHttps defines a proxy server for HTTPS traffic only
var FieldNetworkProxyHttps = &Field{
	Name:        "proxy.https",
	Group:       GroupNetwork,
	Type:        FieldTypeString,
	Default:     nil,
	Description: "Set a proxy server for HTTPS traffic",
	ValidateTag: "url",
	Example:     "http://user:password@host:port",
}
