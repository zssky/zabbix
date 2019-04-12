package zabbix

import "github.com/AlekSi/reflector"

// History history data
type History struct {
	ID         string      `json:"id,omitempty"`
	Clock      string      `json:"clock,omitempty"`
	ItemID     string      `json:"itemid,omitempty"`
	LogeventID string      `json:"logeventid,omitempty"`
	NS         string      `json:"ns,omitempty"`
	Severity   string      `json:"severity,omitempty"`
	Source     string      `json:"source,omitempty"`
	Timestamp  string      `json:"timestamp,omitempty"`
	Value      interface{} `json:"value,omitempty"`
}

// Historys - History array
type Historys []History

// HistorysGet Wrapper for History.get: https://www.zabbix.com/documentation/3.0/manual/api/reference/history/get
func (api *API) HistorysGet(params Params) (res Historys, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("history.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}
