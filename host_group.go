package zabbix

import (
	"github.com/AlekSi/reflector"
)

type (
	// InternalType define
	InternalType int
)

const (
	// NotInternal - no
	NotInternal InternalType = 0
	// Internal - yes
	Internal InternalType = 1
)

// HostGroup - https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/definitions
type HostGroup struct {
	ID       string       `json:"groupid,omitempty"`
	Name     string       `json:"name"`
	Internal InternalType `json:"internal,omitempty"`
}

// HostGroups - the array of HostGroup
type HostGroups []HostGroup

// HostGroupID - host group id
type HostGroupID struct {
	GroupID string `json:"groupid"`
}

// HostGroupIds -  host groupd ids
type HostGroupIds []HostGroupID

// HostGroupsGet - Wrapper for hostgroup.get: https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/get
func (api *API) HostGroupsGet(params Params) (res HostGroups, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("hostgroup.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// HostGroupGetByID - Gets host group by Id only if there is exactly 1 matching host group.
func (api *API) HostGroupGetByID(id string) (res *HostGroup, err error) {
	groups, err := api.HostGroupsGet(Params{"groupids": id})
	if err != nil {
		return
	}

	if len(groups) == 1 {
		res = &groups[0]
	} else {
		e := ExpectedOneResult(len(groups))
		err = &e
	}
	return
}

// HostGroupsCreate - Wrapper for hostgroup.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/create
func (api *API) HostGroupsCreate(hostGroups HostGroups) (err error) {
	response, err := api.CallWithError("hostgroup.create", hostGroups)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	groupids := result["groupids"].([]interface{})
	for i, id := range groupids {
		hostGroups[i].ID = id.(string)
	}
	return
}

// HostGroupsDelete - Wrapper for hostgroup.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/delete
// Cleans GroupId in all hostGroups elements if call succeed.
func (api *API) HostGroupsDelete(hostGroups HostGroups) (err error) {
	ids := make([]string, len(hostGroups))
	for i, group := range hostGroups {
		ids[i] = group.ID
	}

	err = api.HostGroupsDeleteByIds(ids)
	if err == nil {
		for i := range hostGroups {
			hostGroups[i].ID = ""
		}
	}
	return
}

// HostGroupsDeleteByIds - Wrapper for hostgroup.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/delete
func (api *API) HostGroupsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("hostgroup.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	groupids := result["groupids"].([]interface{})
	if len(ids) != len(groupids) {
		err = &ExpectedMore{len(ids), len(groupids)}
	}
	return
}
