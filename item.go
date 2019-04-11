package zabbix

import (
	"fmt"

	"github.com/AlekSi/reflector"
)

type (
	// ItemType - item type
	ItemType int
	// ValueType - value type
	ValueType int
	// DataType - data type
	DataType int
	// DeltaType - delta type
	DeltaType int
)

const (
	// ZabbixAgent - Zabbix agent
	ZabbixAgent ItemType = 0
	// SNMPv1Agent - SNMPv1 agent
	SNMPv1Agent ItemType = 1
	// ZabbixTrapper - Zabbix trapper
	ZabbixTrapper ItemType = 2
	// SimpleCheck - simple check
	SimpleCheck ItemType = 3
	// SNMPv2Agent - SNMPv2 agent
	SNMPv2Agent ItemType = 4
	// ZabbixInternal - Zabbix internal
	ZabbixInternal ItemType = 5
	// SNMPv3Agent - SNMPv3 agent
	SNMPv3Agent ItemType = 6
	// ZabbixAgentActive - Zabbix agent (active)
	ZabbixAgentActive ItemType = 7
	// ZabbixAggregate - Zabbix aggregate
	ZabbixAggregate ItemType = 8
	// WebItem - web item
	WebItem ItemType = 9
	// ExternalCheck - external check
	ExternalCheck ItemType = 10
	// DatabaseMonitor - database monitor
	DatabaseMonitor ItemType = 11
	// IPMIAgent - IPMI agent
	IPMIAgent ItemType = 12
	// SSHAgent - SSH agent
	SSHAgent ItemType = 13
	// TELNETAgent - TELNET agent
	TELNETAgent ItemType = 14
	// Calculated - calculated
	Calculated ItemType = 15
	// JMXAgent - JMX agent
	JMXAgent ItemType = 16
	// SNMPTrap - SNMP trap
	SNMPTrap ItemType = 17

	// Float -  numeric float
	Float ValueType = 0
	// Character - character
	Character ValueType = 1
	// Log - log
	Log ValueType = 2
	// Unsigned - numeric unsigned
	Unsigned ValueType = 3
	// Text - text
	Text ValueType = 4

	// Decimal - (default) decimal
	Decimal DataType = 0
	// Octal - octal
	Octal DataType = 1
	// Hexadecimal - hexadecimal
	Hexadecimal DataType = 2
	// Boolean - boolean
	Boolean DataType = 3

	// AsIs - (default) as is
	AsIs DeltaType = 0
	// Speed - Delta, speed per second
	Speed DeltaType = 1
	// Delta - Delta, simple change
	Delta DeltaType = 2
)

// Item - https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/definitions
type Item struct {
	ID          string    `json:"itemid,omitempty"`
	Delay       int       `json:"delay"`
	HostID      string    `json:"hostid"`
	InterfaceID string    `json:"interfaceid,omitempty"`
	Key         string    `json:"key_"`
	Name        string    `json:"name"`
	Type        ItemType  `json:"type"`
	ValueType   ValueType `json:"value_type"`
	DataType    DataType  `json:"data_type"`
	Delta       DeltaType `json:"delta"`
	Description string    `json:"description"`
	Error       string    `json:"error"`
	History     int       `json:"history,omitempty"`
	Trends      int       `json:"trends,omitempty"`

	// Fields below used only when creating applications
	ApplicationIds []string `json:"applications,omitempty"`
}

// Items - the array of item
type Items []Item

// ByKey - Converts slice to map by key. Panics if there are duplicate keys.
func (items Items) ByKey() (res map[string]Item) {
	res = make(map[string]Item, len(items))
	for _, i := range items {
		_, present := res[i.Key]
		if present {
			panic(fmt.Errorf("Duplicate key %s", i.Key))
		}
		res[i.Key] = i
	}
	return
}

// ItemsGet - Wrapper for item.get https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/get
func (api *API) ItemsGet(params Params) (res Items, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("item.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// ItemsGetByApplicationID - Gets items by application Id.
func (api *API) ItemsGetByApplicationID(id string) (res Items, err error) {
	return api.ItemsGet(Params{"applicationids": id})
}

// ItemsCreate - Wrapper for item.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/create
func (api *API) ItemsCreate(items Items) (err error) {
	response, err := api.CallWithError("item.create", items)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids := result["itemids"].([]interface{})
	for i, id := range itemids {
		items[i].ID = id.(string)
	}
	return
}

// ItemsDelete - Wrapper for item.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/delete
// Cleans ItemId in all items elements if call succeed.
func (api *API) ItemsDelete(items Items) (err error) {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	err = api.ItemsDeleteByIds(ids)
	if err == nil {
		for i := range items {
			items[i].ID = ""
		}
	}
	return
}

// ItemsDeleteByIds - Wrapper for item.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/delete
func (api *API) ItemsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("item.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids1, ok := result["itemids"].([]interface{})
	l := len(itemids1)
	if !ok {
		// some versions actually return map there
		itemids2 := result["itemids"].(map[string]interface{})
		l = len(itemids2)
	}
	if len(ids) != l {
		err = &ExpectedMore{len(ids), l}
	}
	return
}
