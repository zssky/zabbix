package zabbix

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func CreateApplication(host *Host, t *testing.T) *Application {
	apps := Applications{{HostID: host.ID, Name: fmt.Sprintf("App %d for %s", rand.Int(), host.Host)}}
	err := getAPI(t).ApplicationsCreate(apps)
	if err != nil {
		t.Fatal(err)
	}
	return &apps[0]
}

func DeleteApplication(app *Application, t *testing.T) {
	err := getAPI(t).ApplicationsDelete(Applications{*app})
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplications(t *testing.T) {
	api := getAPI(t)

	group := CreateHostGroup(t)
	defer DeleteHostGroup(group, t)

	host := CreateHost(group, t)
	defer DeleteHost(host, t)

	app := CreateApplication(host, t)
	if app.ID == "" {
		t.Errorf("Id is empty: %#v", app)
	}

	app2 := CreateApplication(host, t)
	if app2.ID == "" {
		t.Errorf("Id is empty: %#v", app2)
	}
	if reflect.DeepEqual(app, app2) {
		t.Errorf("Apps are equal:\n%#v\n%#v", app, app2)
	}

	apps, err := api.ApplicationsGet(Params{"hostids": host.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Errorf("Failed to create apps: %#v", apps)
	}

	app2, err = api.ApplicationGetByID(app.ID)
	if err != nil {
		t.Fatal(err)
	}
	app2.TemplateID = ""
	if !reflect.DeepEqual(app, app2) {
		t.Errorf("Apps are not equal:\n%#v\n%#v", app, app2)
	}

	app2, err = api.ApplicationGetByHostIDAndName(host.ID, app.Name)
	if err != nil {
		t.Fatal(err)
	}
	app2.TemplateID = ""
	if !reflect.DeepEqual(app, app2) {
		t.Errorf("Apps are not equal:\n%#v\n%#v", app, app2)
	}

	DeleteApplication(app, t)
}
