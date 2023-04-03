//////////////////////////////////////////////////////////////////////
//
// Copyright 2021 Dell, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//////////////////////////////////////////////////////////////////////////

//go:build !campus_pkg

package transformer_test

import (
	"fmt"
	"testing"
	"time"
)

func Test_OC_Telemetry_Post_Sensor_Group(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces/interface[name=Ethernet0]": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Post Sensor Group  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/sensor-groups"
	payload := "{\"openconfig-telemetry:sensor-group\":[{\"sensor-group-id\":\"sen2\",\"config\":{\"sensor-group-id\":\"sen2\"},\"sensor-paths\":{\"sensor-path\":[{\"path\":\"/openconfig-interfaces:interfaces/interface[name=Ethernet0]\",\"config\":{\"path\":\"/openconfig-interfaces:interfaces/interface[name=Ethernet0]\"}}]}}]}"
	t.Run("Post Sensor Group", processSetRequest(url, payload, "POST", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": map[string]interface{}{"NULL": "NULL"}}}
	expected2 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces/interface[name=Ethernet0]": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Post Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_GROUP|sen2", expected1, false))
	t.Run("Verify Post Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_PATH|sen2|/openconfig-interfaces:interfaces/interface[name=Ethernet0]", expected2, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}

func Test_OC_Telemetry_Patch_Sensor_Group(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Patch Sensor Group  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/sensor-groups"
	payload := "{\"openconfig-telemetry:sensor-groups\":{\"sensor-group\":[{\"sensor-group-id\":\"sen2\",\"config\":{\"sensor-group-id\":\"sen2\"},\"sensor-paths\":{\"sensor-path\":[{\"path\":\"/openconfig-interfaces:interfaces\",\"config\":{\"path\":\"/openconfig-interfaces:interfaces\"}}]}}]}}"
	t.Run("Patch Sensor Group", processSetRequest(url, payload, "PATCH", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": map[string]interface{}{"NULL": "NULL"}}}
	expected2 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Patch Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_GROUP|sen2", expected1, false))
	t.Run("Verify Patch Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_PATH|sen2|/openconfig-interfaces:interfaces", expected2, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}

func Test_OC_Telemetry_Put_Sensor_Group(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Put Sensor Group  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/sensor-groups"
	payload := "{\"openconfig-telemetry:sensor-groups\":{\"sensor-group\":[{\"sensor-group-id\":\"sen2\",\"config\":{\"sensor-group-id\":\"sen2\"},\"sensor-paths\":{\"sensor-path\":[{\"path\":\"/openconfig-interfaces:interfaces\",\"config\":{\"path\":\"/openconfig-interfaces:interfaces\"}}]}}]}}"
	t.Run("Put Sensor Group", processSetRequest(url, payload, "PUT", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": map[string]interface{}{"NULL": "NULL"}}}
	expected2 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Put Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_GROUP|sen2", expected1, false))
	t.Run("Verify Put Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_PATH|sen2|/openconfig-interfaces:interfaces", expected2, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}

func Test_OC_Telemetry_Get_Sensor_Group(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": map[string]interface{}{"NULL": "NULL"}}}
	prereq2 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces": map[string]interface{}{"NULL": "NULL"}}}

	url := "/openconfig-telemetry:telemetry-system/sensor-groups"

	fmt.Println("++++++++++++++  Test OC Telemetry Get Sensor Group  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)

	expected := "{ \"openconfig-telemetry:sensor-groups\": { \"sensor-group\": [ { \"config\": { \"sensor-group-id\": \"sen2\" }, \"sensor-group-id\": \"sen2\", \"sensor-paths\": { \"sensor-path\": [ { \"config\": { \"path\": \"/openconfig-interfaces:interfaces\" }, \"path\": \"/openconfig-interfaces:interfaces\", \"state\": { \"path\": \"/openconfig-interfaces:interfaces\" } } ] }, \"state\": { \"sensor-group-id\": \"sen2\" } } ] } }"

	t.Run("Verify Get Sensor Group", processGetRequest(url, nil, expected, false))

	unloadConfigDB(rclient, prereq1)
	unloadConfigDB(rclient, prereq2)
}

func Test_OC_Telemetry_Delete_Sensor_Group(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen2": map[string]interface{}{"NULL": "NULL"}}}
	prereq2 := map[string]interface{}{"DIALOUT_SENSOR_PATH": map[string]interface{}{"sen2|/openconfig-interfaces:interfaces": map[string]interface{}{"NULL": "NULL"}}}

	url := "/openconfig-telemetry:telemetry-system/sensor-groups"

	fmt.Println("++++++++++++++  Test OC Telemetry Delete Sensor Group  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)

	t.Run("Delete Sensor Group", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)

	expected := make(map[string]interface{})

	t.Run("Verify Delete Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_GROUP|sen2", expected, false))
	t.Run("Verify Delete Sensor Group", verifyDbResult(rclient, "DIALOUT_SENSOR_GROUP|sen2|/openconfig-interfaces:interfaces", expected, false))
	// Teardown
	unloadConfigDB(rclient, prereq1)
	unloadConfigDB(rclient, prereq2)
}

func Test_OC_Telemetry_Post_Destination_Group(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|1.1.1.1|8080": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Post Destination Group  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/destination-groups"
	payload := "{\"openconfig-telemetry:destination-group\":[{\"group-id\":\"dest1\",\"config\":{\"group-id\":\"dest1\"},\"destinations\":{\"destination\":[{\"destination-address\":\"1.1.1.1\",\"destination-port\":8080,\"config\":{\"destination-address\":\"1.1.1.1\",\"destination-port\":8080}}]}}]}"
	t.Run("Post Destination Group", processSetRequest(url, payload, "POST", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": map[string]interface{}{"NULL": "NULL"}}}
	expected2 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|1.1.1.1|8080": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Post Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION_GROUP|dest1", expected1, false))
	t.Run("Verify Post Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION|dest1|1.1.1.1|8080", expected2, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}

func Test_OC_Telemetry_Patch_Destination_Group(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|2.2.2.2|1234": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Patch Destination Group  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/destination-groups"
	payload := "{\"openconfig-telemetry:destination-groups\":{\"destination-group\":[{\"group-id\":\"dest1\",\"config\":{\"group-id\":\"dest1\"},\"destinations\":{\"destination\":[{\"destination-address\":\"2.2.2.2\",\"destination-port\":1234,\"config\":{\"destination-address\":\"2.2.2.2\",\"destination-port\":1234}}]}}]}}"
	t.Run("Patch Sensor Group", processSetRequest(url, payload, "PATCH", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": map[string]interface{}{"NULL": "NULL"}}}
	expected2 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|2.2.2.2|1234": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Patch Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION_GROUP|dest1", expected1, false))
	t.Run("Verify Patch Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION|dest1|2.2.2.2|1234", expected2, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}

func Test_OC_Telemetry_Put_Destination_Group(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|2.2.2.2|1234": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Put Destination Group  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/destination-groups"
	payload := "{\"openconfig-telemetry:destination-groups\":{\"destination-group\":[{\"group-id\":\"dest1\",\"config\":{\"group-id\":\"dest1\"},\"destinations\":{\"destination\":[{\"destination-address\":\"2.2.2.2\",\"destination-port\":1234,\"config\":{\"destination-address\":\"2.2.2.2\",\"destination-port\":1234}}]}}]}}"
	t.Run("Put Sensor Group", processSetRequest(url, payload, "PUT", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": map[string]interface{}{"NULL": "NULL"}}}
	expected2 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|2.2.2.2|1234": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Put Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION_GROUP|dest1", expected1, false))
	t.Run("Verify Put Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION|dest1|2.2.2.2|1234", expected2, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}

func Test_OC_Telemetry_Get_Destination_Group(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": map[string]interface{}{"NULL": "NULL"}}}
	prereq2 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|2.2.2.2|1234": map[string]interface{}{"NULL": "NULL"}}}

	url := "/openconfig-telemetry:telemetry-system/destination-groups"

	fmt.Println("++++++++++++++  Test OC Telemetry Get Destination Group  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)

	expected := "{ \"openconfig-telemetry:destination-groups\": { \"destination-group\": [ { \"config\": { \"group-id\": \"dest1\" }, \"destinations\": { \"destination\": [ { \"config\": { \"destination-address\": \"2.2.2.2\", \"destination-port\": 1234 }, \"destination-address\": \"2.2.2.2\", \"destination-port\": 1234, \"state\": { \"destination-address\": \"2.2.2.2\", \"destination-port\": 1234 } } ] }, \"group-id\": \"dest1\", \"state\": { \"group-id\": \"dest1\" } } ] } }"

	t.Run("Verify Get Destination Group", processGetRequest(url, nil, expected, false))

	unloadConfigDB(rclient, prereq1)
	unloadConfigDB(rclient, prereq2)
}

func Test_OC_Telemetry_Delete_Destination_Group(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": map[string]interface{}{"NULL": "NULL"}}}
	prereq2 := map[string]interface{}{"DIALOUT_DESTINATION": map[string]interface{}{"dest1|2.2.2.2|1234": map[string]interface{}{"NULL": "NULL"}}}

	url := "/openconfig-telemetry:telemetry-system/destination-groups"

	fmt.Println("++++++++++++++  Test OC Telemetry Delete Destination Group  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)

	t.Run("Delete Destination Group", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)

	expected := make(map[string]interface{})

	t.Run("Verify Delete Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION_GROUP|dest1", expected, false))
	t.Run("Verify Delete Destination Group", verifyDbResult(rclient, "DIALOUT_DESTINATION|dest1|2.2.2.2|1234", expected, false))
	// Teardown
	unloadConfigDB(rclient, prereq1)
	unloadConfigDB(rclient, prereq2)
}

func Test_OC_Telemetry_Post_Persistent_Subscription(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen1": map[string]interface{}{"NULL": "NULL"}}}
	prereq2 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": map[string]interface{}{"NULL": "NULL"}}}

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": ""}}
	cleanuptbl3 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Post Persistent Subscription  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)

	url := "/openconfig-telemetry:telemetry-system/subscriptions/persistent-subscriptions"
	payload := "{\"openconfig-telemetry:persistent-subscription\":[{\"name\":\"sub1\",\"config\":{\"name\":\"sub1\",\"local-source-address\":\"1.1.1.1\",\"originated-qos-marking\":10,\"protocol\":\"STREAM_SSH\",\"encoding\":\"ENC_XML\"},\"sensor-profiles\":{\"sensor-profile\":[{\"sensor-group\":\"sen1\",\"config\":{\"sensor-group\":\"sen1\",\"sample-interval\":\"60\",\"heartbeat-interval\":\"50\",\"suppress-redundant\":true}}]},\"destination-groups\":{\"destination-group\":[{\"group-id\":\"dest1\",\"config\":{\"group-id\":\"dest1\"}}]}}]}"
	t.Run("Post Persistent Subscription", processSetRequest(url, payload, "POST", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": map[string]interface{}{"encoding": "ENC_XML", "local_source_address": "1.1.1.1", "originated_qos_marking": "10", "protocol": "STREAM_SSH"}}}
	expected2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": map[string]interface{}{"heartbeat_interval": "50", "sample_interval": "60", "suppress_redundant": "true"}}}
	expected3 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Post Persistent Subscription", verifyDbResult(rclient, "DIALOUT_PERSISTENT_SUBSCRIPTION|sub1", expected1, false))
	t.Run("Verify Post Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_SENSOR_PROFILE|sub1|sen1", expected2, false))
	t.Run("Verify Post Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_DESTINATION_GROUP|sub1|dest1", expected3, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
	unloadConfigDB(rclient, cleanuptbl3)
}

func Test_OC_Telemetry_Patch_Persistent_Subscription(t *testing.T) {

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub2|dest1": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub2|sen1": ""}}
	cleanuptbl3 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub2": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Patch Persistent Subscription  +++++++++++++")

	url := "/openconfig-telemetry:telemetry-system/subscriptions/persistent-subscriptions"
	payload := "{\"openconfig-telemetry:persistent-subscriptions\":{\"persistent-subscription\":[{\"name\":\"sub2\",\"config\":{\"name\":\"sub2\",\"local-source-address\":\"1.1.1.1\",\"originated-qos-marking\":10,\"protocol\":\"STREAM_SSH\",\"encoding\":\"ENC_XML\"},\"sensor-profiles\":{\"sensor-profile\":[{\"sensor-group\":\"sen1\",\"config\":{\"sensor-group\":\"sen1\",\"sample-interval\":\"60\",\"heartbeat-interval\":\"50\",\"suppress-redundant\":true}}]},\"destination-groups\":{\"destination-group\":[{\"group-id\":\"dest1\",\"config\":{\"group-id\":\"dest1\"}}]}}]}}"
	t.Run("Patch Sensor Group", processSetRequest(url, payload, "PATCH", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub2": map[string]interface{}{"encoding": "ENC_XML", "local_source_address": "1.1.1.1", "originated_qos_marking": "10", "protocol": "STREAM_SSH"}}}
	expected2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub2|sen1": map[string]interface{}{"heartbeat_interval": "50", "sample_interval": "60", "suppress_redundant": "true"}}}
	expected3 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub2|dest1": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Patch Persistent Subscription", verifyDbResult(rclient, "DIALOUT_PERSISTENT_SUBSCRIPTION|sub2", expected1, false))
	t.Run("Verify Patch Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_SENSOR_PROFILE|sub2|sen1", expected2, false))
	t.Run("Verify Patch Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_DESTINATION_GROUP|sub2|dest1", expected3, false))
	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
	unloadConfigDB(rclient, cleanuptbl3)
}

func Test_OC_Telemetry_Put_Persistent_Subscription(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": map[string]interface{}{"encoding": "ENC_XML", "local_source_address": "1.1.1.1", "originated_qos_marking": "10", "protocol": "STREAM_SSH"}}}
	prereq2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": map[string]interface{}{"heartbeat_interval": "50", "sample_interval": "60", "suppress_redundant": "true"}}}
	prereq3 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": map[string]interface{}{"NULL": "NULL"}}}

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": ""}}
	cleanuptbl3 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": ""}}

	fmt.Println("++++++++++++++  Test OC Telemetry Put Persistent Subscription  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)
	loadConfigDB(rclient, prereq3)

	url := "/openconfig-telemetry:telemetry-system/subscriptions/persistent-subscriptions"
	payload := "{\"openconfig-telemetry:persistent-subscriptions\":{\"persistent-subscription\":[{\"name\":\"sub1\",\"config\":{\"name\":\"sub1\",\"local-source-address\":\"2.2.2.2\",\"originated-qos-marking\":20,\"protocol\":\"STREAM_GRPC\",\"encoding\":\"ENC_JSON_IETF\"},\"sensor-profiles\":{\"sensor-profile\":[{\"sensor-group\":\"sen1\",\"config\":{\"sensor-group\":\"sen1\",\"sample-interval\":\"50\",\"heartbeat-interval\":\"40\",\"suppress-redundant\":false}}]},\"destination-groups\":{\"destination-group\":[{\"group-id\":\"dest1\",\"config\":{\"group-id\":\"dest1\"}}]}}]}}"
	t.Run("Put Sensor Group", processSetRequest(url, payload, "PUT", false))
	time.Sleep(1 * time.Second)
	expected1 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": map[string]interface{}{"encoding": "ENC_JSON_IETF", "local_source_address": "2.2.2.2", "originated_qos_marking": "20", "protocol": "STREAM_GRPC"}}}
	expected2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": map[string]interface{}{"heartbeat_interval": "40", "sample_interval": "50", "suppress_redundant": "false"}}}
	expected3 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": map[string]interface{}{"NULL": "NULL"}}}

	t.Run("Verify Put Persistent Subscription", verifyDbResult(rclient, "DIALOUT_PERSISTENT_SUBSCRIPTION|sub1", expected1, false))
	t.Run("Verify Put Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_SENSOR_PROFILE|sub1|sen1", expected2, false))
	t.Run("Verify Put Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_DESTINATION_GROUP|sub1|dest1", expected3, false))

	// Teardown
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
	unloadConfigDB(rclient, cleanuptbl3)
}

func Test_OC_Telemetry_Get_Persistent_Subscription(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": map[string]interface{}{"encoding": "ENC_XML", "local_source_address": "1.1.1.1", "originated_qos_marking": "10", "protocol": "STREAM_SSH"}}}
	prereq2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": map[string]interface{}{"heartbeat_interval": "50", "sample_interval": "60", "suppress_redundant": "true"}}}
	prereq3 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": map[string]interface{}{"NULL": "NULL"}}}

	url := "/openconfig-telemetry:telemetry-system/subscriptions/persistent-subscriptions"

	fmt.Println("++++++++++++++  Test OC Telemetry Get Persistent Subscription  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)
	loadConfigDB(rclient, prereq3)

	expected := "{ \"openconfig-telemetry:persistent-subscriptions\": { \"persistent-subscription\": [ { \"config\": { \"encoding\": \"openconfig-telemetry-types:ENC_XML\", \"local-source-address\": \"1.1.1.1\", \"name\": \"sub1\", \"originated-qos-marking\": 10, \"protocol\": \"openconfig-telemetry-types:STREAM_SSH\" }, \"destination-groups\": { \"destination-group\": [ { \"config\": { \"group-id\": \"dest1\" }, \"group-id\": \"dest1\", \"state\": { \"group-id\": \"dest1\" } } ] }, \"name\": \"sub1\", \"sensor-profiles\": { \"sensor-profile\": [ { \"config\": { \"heartbeat-interval\": \"50\", \"sample-interval\": \"60\", \"sensor-group\": \"sen1\", \"suppress-redundant\": true }, \"sensor-group\": \"sen1\", \"state\": { \"heartbeat-interval\": \"50\", \"sample-interval\": \"60\", \"sensor-group\": \"sen1\", \"suppress-redundant\": true } } ] }, \"state\": { \"encoding\": \"openconfig-telemetry-types:ENC_XML\", \"local-source-address\": \"1.1.1.1\", \"name\": \"sub1\", \"originated-qos-marking\": 10, \"protocol\": \"openconfig-telemetry-types:STREAM_SSH\" } } ] } }"

	t.Run("Verify Get Persistent Subscription", processGetRequest(url, nil, expected, false))

	unloadConfigDB(rclient, prereq1)
	unloadConfigDB(rclient, prereq2)
	unloadConfigDB(rclient, prereq3)
}

func Test_OC_Telemetry_Delete_Persistent_Subscription(t *testing.T) {

	prereq1 := map[string]interface{}{"DIALOUT_PERSISTENT_SUBSCRIPTION": map[string]interface{}{"sub1": map[string]interface{}{"encoding": "ENC_XML", "local_source_address": "1.1.1.1", "originated_qos_marking": "10", "protocol": "STREAM_SSH"}}}
	prereq2 := map[string]interface{}{"DIALOUT_SUBSCR_SENSOR_PROFILE": map[string]interface{}{"sub1|sen1": map[string]interface{}{"heartbeat_interval": "50", "sample_interval": "60", "suppress_redundant": "true"}}}
	prereq3 := map[string]interface{}{"DIALOUT_SUBSCR_DESTINATION_GROUP": map[string]interface{}{"sub1|dest1": map[string]interface{}{"NULL": "NULL"}}}

	cleanuptbl1 := map[string]interface{}{"DIALOUT_SENSOR_GROUP": map[string]interface{}{"sen1": ""}}
	cleanuptbl2 := map[string]interface{}{"DIALOUT_DESTINATION_GROUP": map[string]interface{}{"dest1": ""}}

	url := "/openconfig-telemetry:telemetry-system/subscriptions/persistent-subscriptions"

	fmt.Println("++++++++++++++  Test OC Telemetry Delete Persistent Subscription  +++++++++++++")

	// Setup - Prerequisite
	loadConfigDB(rclient, prereq1)
	loadConfigDB(rclient, prereq2)
	loadConfigDB(rclient, prereq3)

	t.Run("Delete Persistent Subscription", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)

	expected := make(map[string]interface{})

	t.Run("Verify Delete Persistent Subscription", verifyDbResult(rclient, "DIALOUT_PERSISTENT_SUBSCRIPTION|sub1", expected, false))
	t.Run("Verify Delete Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_SENSOR_PROFILE|sub1|sen1", expected, false))
	t.Run("Verify Delete Persistent Subscription", verifyDbResult(rclient, "DIALOUT_SUBSCR_DESTINATION_GROUP|sub1|dest1", expected, false))

	// Teardown
	unloadConfigDB(rclient, prereq1)
	unloadConfigDB(rclient, prereq2)
	unloadConfigDB(rclient, prereq3)
	unloadConfigDB(rclient, cleanuptbl1)
	unloadConfigDB(rclient, cleanuptbl2)
}
