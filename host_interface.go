package zabbix

type (
	// InterfaceType - interface type
	InterfaceType int
)

const (
	// Agent - type agent
	Agent InterfaceType = 1
	// SNMP - type snmp
	SNMP InterfaceType = 2
	// IPMI - type IPMI
	IPMI InterfaceType = 3
	// JMX - type JMX
	JMX InterfaceType = 4
)

// HostInterface - https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostinterface/definitions
type HostInterface struct {
	DNS   string        `json:"dns"`
	IP    string        `json:"ip"`
	Main  int           `json:"main"`
	Port  string        `json:"port"`
	Type  InterfaceType `json:"type"`
	UseIP int           `json:"useip"`
}

// HostInterfaces - host interface
type HostInterfaces []HostInterface
