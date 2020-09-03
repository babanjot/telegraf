package opcua_client

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/influxdata/telegraf/config"
)

type OPCTags struct {
	Name           string
	Namespace      string
	IdentifierType string
	Identifier     string
	DataType       string
	Want           string
}

func TestClient1(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var testopctags = []OPCTags{
		{"ProductName", "0", "i", "2261", "string", "open62541 OPC UA Server"},
		{"ProductUri", "0", "i", "2262", "string", "http://open62541.org"},
		{"ManufacturerName", "0", "i", "2263", "string", "open62541"},
	}

	var o OpcUA
	var err error

	o.Name = "testing"
	o.Endpoint = "opc.tcp://opcua.rocks:4840"
	o.AuthMethod = "Anonymous"
	o.Interval = config.Duration(10 * time.Millisecond)
	o.ConnectTimeout = config.Duration(10 * time.Second)
	o.RequestTimeout = config.Duration(1 * time.Second)
	o.SecurityPolicy = "None"
	o.SecurityMode = "None"
	for _, tags := range testopctags {
		o.NodeList = append(o.NodeList, MapOPCTag(tags))
	}
	err = o.Init()
	if err != nil {
		t.Errorf("Initialize Error: %s", err)
	}
	err = Connect(&o)
	if err != nil {
		t.Fatalf("Connect Error: %s", err)
	}

	for i, v := range o.NodeData {
		if v.Value != nil {
			types := reflect.TypeOf(v.Value)
			value := reflect.ValueOf(v.Value)
			compare := fmt.Sprintf("%v", value.Interface())
			if compare != testopctags[i].Want {
				t.Errorf("Tag %s: Values %v for type %s  does not match record", o.NodeList[i].Name, value.Interface(), types)
			}
		} else {
			t.Errorf("Tag: %s has value: %v", o.NodeList[i].Name, v.Value)
		}
	}
}

func MapOPCTag(tags OPCTags) (out OPCTag) {
	out.Name = tags.Name
	out.Namespace = tags.Namespace
	out.IdentifierType = tags.IdentifierType
	out.Identifier = tags.Identifier
	out.DataType = tags.DataType
	return out
}
