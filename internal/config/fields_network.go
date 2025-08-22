package config

const GroupNetwork = "Network"

func init() {
	Fields.Add(
		FieldNetworkProxyAll,
		FieldNetworkProxyHttp,
		FieldNetworkProxyHttps,
	)
}

var FieldNetworkProxyAll = &Field{
	Name:        "proxy.all",
	Group:       GroupNetwork,
	Type:        FieldTypeString,
	Default:     nil,
	Description: "Set a proxy server for all network traffic",
	ValidateTag: "url",
	Example:     "http://user:password@host:port",
}

var FieldNetworkProxyHttp = &Field{
	Name:        "proxy.http",
	Group:       GroupNetwork,
	Type:        FieldTypeString,
	Default:     nil,
	Description: "Set a proxy server for HTTP traffic",
	ValidateTag: "url",
	Example:     "http://user:password@host:port",
}

var FieldNetworkProxyHttps = &Field{
	Name:        "proxy.https",
	Group:       GroupNetwork,
	Type:        FieldTypeString,
	Default:     nil,
	Description: "Set a proxy server for HTTPS traffic",
	ValidateTag: "url",
	Example:     "http://user:password@host:port",
}
