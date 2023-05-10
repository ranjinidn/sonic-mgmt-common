////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Broadcom. The term Broadcom refers to Broadcom Inc. and/or //
//  its subsidiaries.                                                         //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//     http://www.apache.org/licenses/LICENSE-2.0                             //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

//go:build test
// +build test

package transformer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
)

func Test_node_exercising_subtree_xfmr_and_virtual_table(t *testing.T) {
	var pre_req_map, expected_map, cleanuptbl map[string]interface{}
	var url, url_body_json string

	fmt.Println("\n\n+++++++++++++ Performing Set on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")
	url = "/openconfig-test-xfmr:test-xfmr/interfaces"
	url_body_json = "{ \"openconfig-test-xfmr:interface\": [ { \"id\": \"Eth_0\", \"config\": { \"id\": \"Eth_0\" }, \"ingress-test-sets\": { \"ingress-test-set\": [ { \"set-name\": \"TestSet_01\", \"type\": \"TEST_SET_IPV4\", \"config\": { \"set-name\": \"TestSet_01\", \"type\": \"TEST_SET_IPV4\" } } ] } } ]}"
	expected_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": map[string]interface{}{"ports@": "Eth_0", "type": "IPV4"}}}
	cleanuptbl = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": ""}}
	t.Run("Test set on node exercising subtree-xfmr and virtual table.", processSetRequest(url, url_body_json, "POST", false, nil))
	time.Sleep(1 * time.Second)
	t.Run("Verify set on node exercising subtree-xfmr and virtual table.", verifyDbResult(rclient, "TEST_SET_TABLE|TestSet_01_TEST_SET_IPV4", expected_map, false))
	unloadDB(db.ConfigDB, cleanuptbl)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n+++++++++++++ Done Performing Set on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")

	fmt.Println("\n\n+++++++++++++ Performing Delete on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")
	pre_req_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": map[string]interface{}{
		"ports@": "Eth_0,Eth_1,Eth_3"},
		"TestSet_02_TEST_SET_IPV4": map[string]interface{}{
			"ports@": "Eth_1,Eth_4"}}}
	cleanuptbl = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": "", "TestSet_02_TEST_SET_IPV4": ""}}
	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)
	url = "/openconfig-test-xfmr:test-xfmr/interfaces/interface[id=Eth_1]"
	t.Run("Test delete on node exercising subtree-xfmr and virtual table.", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)
	expected_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": map[string]interface{}{
		"ports@": "Eth_0,Eth_3"}}}
	t.Run("Verify delete on node exercising subtree-xfmr and virtual table (TestSet_01).", verifyDbResult(rclient, "TEST_SET_TABLE|TestSet_01_TEST_SET_IPV4", expected_map, false))
	expected_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_02_TEST_SET_IPV4": map[string]interface{}{
		"ports@": "Eth_4"}}}
	t.Run("Verify delete on node exercising subtree-xfmr and virtual table (TestSet_02).", verifyDbResult(rclient, "TEST_SET_TABLE|TestSet_02_TEST_SET_IPV4", expected_map, false))
	unloadDB(db.ConfigDB, cleanuptbl)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n+++++++++++++ Done Performing Delete on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")

	fmt.Println("\n\n+++++++++++++ Performing Get on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")
	pre_req_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_03_TEST_SET_IPV6": map[string]interface{}{
		"ports@": "Eth_1"}}}

	loadDB(db.ConfigDB, pre_req_map)
	expected_get_json := "{\"openconfig-test-xfmr:ingress-test-set\":[{\"config\":{\"set-name\":\"TestSet_03\",\"type\":\"openconfig-test-xfmr:TEST_SET_IPV6\"},\"set-name\":\"TestSet_03\",\"state\":{\"set-name\":\"TestSet_03\",\"type\":\"openconfig-test-xfmr:TEST_SET_IPV6\"},\"type\":\"openconfig-test-xfmr:TEST_SET_IPV6\"}]}"
	url = "/openconfig-test-xfmr:test-xfmr/interfaces/interface[id=Eth_1]/ingress-test-sets/ingress-test-set[set-name=TestSet_03][type=TEST_SET_IPV6]"
	t.Run("Test get on node exercising subtree-xfmr and virtual table.", processGetRequest(url, expected_get_json, false))
	cleanuptbl = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_03_TEST_SET_IPV6": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	fmt.Println("\n\n+++++++++++++ Done Performing Get on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")

}

func Test_node_exercising_tableName_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, expected_map, cleanuptbl map[string]interface{}
	var url, url_body_json string

	fmt.Println("\n\n+++++++++++++ Performing Set on Yang Node Exercising Table-Name, Key-Xfmr and Field-Xfmr ++++++++++++")
	url = "/openconfig-test-xfmr:test-xfmr/test-sets"
	url_body_json = "{ \"openconfig-test-xfmr:test-set\": [ { \"name\": \"TestSet_01\", \"type\": \"TEST_SET_IPV4\", \"config\": { \"name\": \"TestSet_01\", \"type\": \"TEST_SET_IPV4\", \"description\": \"TestSet_01Description\" } } ]}"
	expected_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": map[string]interface{}{"type": "IPV4", "description": "Description : TestSet_01Description"}}}
	cleanuptbl = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": ""}}
	t.Run("Test set on node exercising Table-Name, Key-Xfmr and Field-Xfmr", processSetRequest(url, url_body_json, "POST", false, nil))
	time.Sleep(1 * time.Second)
	t.Run("Verify set on node exercising Table-Name, Key-Xfmr and Field-Xfmr", verifyDbResult(rclient, "TEST_SET_TABLE|TestSet_01_TEST_SET_IPV4", expected_map, false))
	unloadDB(db.ConfigDB, cleanuptbl)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n+++++++++++++ Done Performing Set on Yang Node Exercising Table-Name, Key-Xfmr and Field-Xfmr ++++++++++++")

	fmt.Println("\n\n+++++++++++++ Performing Delete on Yang Node Exercising Table-Name ,Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": map[string]interface{}{
		"type":        "IPV4",
		"description": "Description : TestSet_01_description",
		"ports@":      "Eth_0"}}}
	cleanuptbl = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": ""}}
	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)
	url = "/openconfig-test-xfmr:test-xfmr/test-sets/test-set[name=TestSet_01][type=TEST_SET_IPV4]/config/description"
	t.Run("Test delete on node exercising Table-Name, Key-Xfmr and Field-Xfmr", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)
	expected_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_01_TEST_SET_IPV4": map[string]interface{}{
		"type":   "IPV4",
		"ports@": "Eth_0"}}}
	t.Run("Verify delete on node exercising Table-Name, Key-Xfmr and Field-Xfmr", verifyDbResult(rclient, "TEST_SET_TABLE|TestSet_01_TEST_SET_IPV4", expected_map, false))
	unloadDB(db.ConfigDB, cleanuptbl)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n+++++++++++++ Done Performing Delete on Yang Node Exercising Subtree-Xfmr and Virtual Table ++++++++++++")

	fmt.Println("\n\n+++++++++++++ Performing Get on Yang Node Exercising Table-Name, Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_03_TEST_SET_IPV6": map[string]interface{}{
		"type":        "IPV6",
		"description": "Description : TestSet_03Description",
		"ports@":      "Eth_3"}}}

	loadDB(db.ConfigDB, pre_req_map)
	expected_get_json := "{\"openconfig-test-xfmr:test-sets\":{\"test-set\":[{\"config\":{\"description\":\"TestSet_03Description\",\"name\":\"TestSet_03\",\"type\":\"openconfig-test-xfmr:TEST_SET_IPV6\"},\"name\":\"TestSet_03\",\"state\":{\"description\":\"TestSet_03Description\",\"name\":\"TestSet_03\",\"type\":\"openconfig-test-xfmr:TEST_SET_IPV6\"},\"type\":\"openconfig-test-xfmr:TEST_SET_IPV6\"}]}}"
	url = "/openconfig-test-xfmr:test-xfmr/test-sets"
	t.Run("Test get on node exercising Table-Name, Key-Xfmr and Field-Xfmr.", processGetRequest(url, expected_get_json, false))
	time.Sleep(1 * time.Second)
	cleanuptbl = map[string]interface{}{"TEST_SET_TABLE": map[string]interface{}{"TestSet_03_TEST_SET_IPV6": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	fmt.Println("\n\n+++++++++++++ Done Performing Get on Yang Node Exercising  Table-Name, Key-Xfmr and Field-Xfmr++++++++++++")
}

func Test_node_exercising_pre_xfmr_node(t *testing.T) {
	fmt.Println("\n\n+++++++++++++ Performing set on node exercising pre-xfmr ++++++++++++")
	err_str := "REPLACE not supported at this node."
	expected_err := tlerr.NotSupportedError{Format: err_str}
	//expected_err := tlerr.NotSupported("REPLACE not supported at this node.")
	url := "/openconfig-test-xfmr:test-xfmr/test-sets"
	url_body_json := "{ \"openconfig-test-xfmr:test-sets\": { \"test-set\": [ { \"name\": \"TestSet_03\", \"type\": \"TEST_SET_IPV4\", \"config\": { \"name\": \"TestSet_03\", \"type\": \"TEST_SET_IPV4\", \"description\": \"testSet_03 description\" } } ] }}"
	t.Run("Test set on node exercising pre-xfmr.", processSetRequest(url, url_body_json, "PUT", true, expected_err))
	fmt.Println("\n\n+++++++++++++ Done Performing set on node exercising pre-xfmr ++++++++++++")

}

func Test_node_with_child_tableXfmr_keyXfmr_fieldNameXfmrs_nonConfigDB_data(t *testing.T) {

	cleanuptbl := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": ""}, "TEST_SENSOR_A_TABLE": map[string]interface{}{"test_group_1|sensor_type_a_testA": ""}}
	url := "/openconfig-test-xfmr:test-xfmr/test-sensor-groups"

	fmt.Println("++++++++++++++  Test_set_on_node_with_child_table_key_field_xfmrs  +++++++++++++")

	// Setup - Prerequisite
	unloadDB(db.ConfigDB, cleanuptbl)

	// Payload
	post_payload := "{\"openconfig-test-xfmr:test-sensor-group\":[ { \"id\" : \"test_group_1\", \"config\": { \"id\": \"test_group_1\", \"group-colors\": [ \"red,blue,green\" ] }, \"test-sensor-types\": { \"test-sensor-type\": [ { \"type\": \"sensora_testA\", \"config\": { \"type\": \"sensora_testA\", \"exclude-filter\": \"filterB\" } } ] } } ]}"
	post_sensor_group_expected := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "10"}}}
	post_sensor_table_expected := map[string]interface{}{"TEST_SENSOR_A_TABLE": map[string]interface{}{"test_group_1|sensor_type_a_testA": map[string]interface{}{"exclude_filter": "filter_filterB"}}}

	t.Run("Set on Node having child table and field transformer mapping", processSetRequest(url, post_payload, "POST", false))
	time.Sleep(1 * time.Second)
	t.Run("Verify set on node with child table and field transformer", verifyDbResult(rclient, "TEST_SENSOR_GROUP|test_group_1", post_sensor_group_expected, false))
	t.Run("Verify set on node with child table and field transformer", verifyDbResult(rclient, "TEST_SENSOR_A_TABLE|test_group_1|sensor_type_a_testA", post_sensor_table_expected, false))

	// Teardown
	unloadDB(db.ConfigDB, cleanuptbl)

	fmt.Println("++++++++++++++  Test_get_on_node_with_table_key_field_xfmrs_nonConfigDB_data  +++++++++++++")

	prereq := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "10"}}, "TEST_SENSOR_A_TABLE": map[string]interface{}{"test_group_1|sensor_type_a_testA": map[string]interface{}{"exclude_filter": "filter_filterB"}}}
	nonconfig_prereq := map[string]interface{}{"TEST_SENSOR_GROUP_COUNTERS": map[string]interface{}{"test_group_1": map[string]interface{}{"frame-in": "12345", "frame-out": "678910"}}}

	url = "/openconfig-test-xfmr:test-xfmr/test-sensor-groups/test-sensor-group[id=test_group_1]"

	// Setup - Prerequisite
	loadDB(db.ConfigDB, prereq)
	loadDB(db.CountersDB, nonconfig_prereq)

	get_expected := "{\"openconfig-test-xfmr:test-sensor-group\":[{\"config\":{\"color-hold-time\":10,\"group-colors\":[\"red\",\"blue\",\"green\"],\"id\":\"test_group_1\"},\"id\":\"test_group_1\",\"state\":{\"color-hold-time\":10,\"counters\":{\"frame-in\":12345,\"frame-out\":678910},\"group-colors\":[\"red\",\"blue\",\"green\"],\"id\":\"test_group_1\"},\"test-sensor-types\":{\"test-sensor-type\":[{\"config\":{\"exclude-filter\":\"filterB\",\"type\":\"sensora_testA\"},\"state\":{\"exclude-filter\":\"filterB\",\"type\":\"sensora_testA\"},\"type\":\"sensora_testA\"}]}}]}"

	t.Run("Verify_get_on_node_with_child_table_key_field_xfmrs", processGetRequest(url, get_expected, false))

	// Teardown
	unloadDB(db.ConfigDB, cleanuptbl)
	unloadDB(db.CountersDB, nonconfig_prereq)

	fmt.Println("++++++++++++++  Test_delete_on_node_with_child_table_key_field_xfmrs  +++++++++++++")

	url = "/openconfig-test-xfmr:test-xfmr/test-sensor-groups"

	// Setup - Prerequisite - None
	unloadDB(db.ConfigDB, cleanuptbl)
	loadDB(db.ConfigDB, prereq)

	delete_expected := make(map[string]interface{})

	t.Run("Delete on node with child table, key and field xfmrs", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)
	t.Run("Verify delete on node with child table, key and field xfmrs", verifyDbResult(rclient, "TEST_SENSOR_GROUP|test_group_1", delete_expected, false))
	t.Run("Verify delete on node with child table, key and field xfmrs", verifyDbResult(rclient, "TEST_SENSOR_A_TABLE|test_group_1|sensor_type_a_testA", delete_expected, false))

	// Teardown
	unloadDB(db.ConfigDB, cleanuptbl)
}

func Test_delete_on_node_with_default_value(t *testing.T) {

	cleanuptbl := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": ""}}

	url := "/openconfig-test-xfmr:test-xfmr/test-sensor-groups/test-sensor-group[id=test_group_1]/config/color-hold-time"
	prereq := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "30"}}}

	fmt.Println("++++++++++++++  Test_delete_on_node_with_default_value  +++++++++++++")

	// Setup - Prerequisite - None
	unloadDB(db.ConfigDB, cleanuptbl)
	loadDB(db.ConfigDB, prereq)

	// Payload
	del_sensor_group_expected := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "10"}}}

	t.Run("Delete on node having default value", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)
	t.Run("Verify delete on node with default value", verifyDbResult(rclient, "TEST_SENSOR_GROUP|test_group_1", del_sensor_group_expected, false))

	// Teardown
	unloadDB(db.ConfigDB, cleanuptbl)
}

func Test_post_xfmr(t *testing.T) {

	cleanuptbl := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": ""}}
	url := "/openconfig-test-xfmr:test-xfmr/test-sensor-groups/test-sensor-group[id=test_group_1]/config/color-hold-time"
	prereq := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "30"}}}

	fmt.Println("++++++++++++++  Test_post_xfmr  +++++++++++++")

	// Setup - Prerequisite
	unloadDB(db.ConfigDB, cleanuptbl)
	loadDB(db.ConfigDB, prereq)

	patch_payload := "{ \"openconfig-test-xfmr:color-hold-time\": 50}"
	patch_expected := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "50"}}}
	post_expected := map[string]interface{}{"TEST_SENSOR_A_TABLE": map[string]interface{}{"test_group_1|sensor_type_a_post50": map[string]interface{}{"description_a": "Added instance in post xfmr"}}}

	t.Run("Test_post_xfmr", processSetRequest(url, patch_payload, "PATCH", false))
	time.Sleep(1 * time.Second)
	t.Run("Verify Test_post_xfmr", verifyDbResult(rclient, "TEST_SENSOR_GROUP|test_group_1", patch_expected, false))
	t.Run("Verify Test_post_xfmr", verifyDbResult(rclient, "TEST_SENSOR_A_TABLE|test_group_1|sensor_type_a_post50", post_expected, false))

	unloadDB(db.ConfigDB, cleanuptbl)
}

func Test_sonic_yang_node_operations(t *testing.T) {

	cleanuptbl := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_id_123": ""}, "TEST_SENSOR_A_TABLE": map[string]interface{}{"sensor_id_123|sensor_type_a_123": ""}}
	prereq := map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_id_123": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "25"}}}
	url := "/sonic-test-xfmr:sonic-test-xfmr/TEST_SENSOR_A_TABLE"

	fmt.Println("++++++++++++++  Test_set_on_sonic_table_yang_node +++++++++++++")

	// Setup - Prerequisite
	unloadDB(db.ConfigDB, cleanuptbl)
	loadDB(db.ConfigDB, prereq)

	// Payload
	post_payload := "{ \"sonic-test-xfmr:TEST_SENSOR_A_TABLE_LIST\": [ { \"id\": \"sensor_id_123\", \"type\": \"sensor_type_a_123\", \"exclude_filter\": \"filter_123\", \"description_a\": \"description test field for sensor A table\" } ]}"
	post_sensor_table_expected := map[string]interface{}{"TEST_SENSOR_A_TABLE": map[string]interface{}{"sensor_id_123|sensor_type_a_123": map[string]interface{}{"exclude_filter": "filter_123", "description_a": "description test field for sensor A table"}}}

	t.Run("Set on sonic table yang node", processSetRequest(url, post_payload, "POST", false))
	time.Sleep(1 * time.Second)
	t.Run("Verify set on sonic table yang node", verifyDbResult(rclient, "TEST_SENSOR_A_TABLE|sensor_id_123|sensor_type_a_123", post_sensor_table_expected, false))

	// Teardown
	unloadDB(db.ConfigDB, cleanuptbl)

	fmt.Println("++++++++++++++  Test_delete_on_sonic_module  +++++++++++++")

	cleanuptbl = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": ""}, "TEST_SENSOR_A_TABLE": map[string]interface{}{"test_group_1|sensor_type_a_testA": ""}, "TEST_SENSOR_B_TABLE": map[string]interface{}{"test_group_1|sensor_type_b_testB": ""}, "TEST_SET_TABLE": map[string]interface{}{"test_set_1": ""}}

	url = "/sonic-test-xfmr:sonic-test-xfmr"
	prereq = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"test_group_1": map[string]interface{}{"colors@": "red,blue,green", "color-hold-time": "30"}}, "TEST_SENSOR_A_TABLE": map[string]interface{}{"test_group_1|sensor_type_a_testA": map[string]interface{}{"test_group_1|sensor_type_a_testA": map[string]interface{}{"exclude_filter": "filter_filterB"}}}, "TEST_SENSOR_B_TABLE": map[string]interface{}{"test_group_1|sensor_type_b_testB": map[string]interface{}{"exclude_filter": "filter_filterB"}}, "TEST_SET_TABLE": map[string]interface{}{"quert_TEST_SET_IPV4": map[string]interface{}{"type": "IPV4"}}}

	// Setup - Prerequisite
	unloadDB(db.ConfigDB, cleanuptbl)
	loadDB(db.ConfigDB, prereq)

	delete_expected := make(map[string]interface{})

	t.Run("Delete on sonic module", processDeleteRequest(url, false))
	time.Sleep(1 * time.Second)
	t.Run("Verify delete on sonic module table1", verifyDbResult(rclient, "TEST_SENSOR_GROUP|test_group_1", delete_expected, false))
	t.Run("Verify delete on sonic module table2", verifyDbResult(rclient, "TEST_SENSOR_A_TABLE|test_group_1|sensor_type_a_testA", delete_expected, false))
	t.Run("Verify delete on sonic module table3", verifyDbResult(rclient, "TEST_SENSOR_B_TABLE|test_group_1|sensor_type_b_testB", delete_expected, false))
	t.Run("Verify delete on sonic module table4", verifyDbResult(rclient, "TEST_SET_TABLE|quert_TEST_SET_IPV4", delete_expected, false))

	// Teardown
	unloadDB(db.ConfigDB, cleanuptbl)

	fmt.Println("++++++++++++++  Test_get_on_sonic_table_with_key_xfmr  +++++++++++++")

	cleanuptbl = map[string]interface{}{"TEST_SENSOR_MODE_TABLE": map[string]interface{}{"mode:testsensor123:3543": ""}}
	prereq = map[string]interface{}{"TEST_SENSOR_MODE_TABLE": map[string]interface{}{"mode:testsensor123:3543": map[string]interface{}{"description": "Test sensor mode"}}}
	url = "/sonic-test-xfmr:sonic-test-xfmr/TEST_SENSOR_MODE_TABLE"

	// Setup - Prerequisite
	loadDB(db.CountersDB, prereq)

	get_expected := "{\"sonic-test-xfmr:TEST_SENSOR_MODE_TABLE\":{\"TEST_SENSOR_MODE_TABLE_LIST\":[{\"description\":\"Test sensor mode\",\"id\":3543,\"mode\":\"mode:testsensor123\"}]}}"
	t.Run("Get on Sonic table with key xfmr", processGetRequest(url, get_expected, false))

	// Teardown
	unloadDB(db.CountersDB, cleanuptbl)
}

func Test_leaflist_node(t *testing.T) {
	var pre_req_map, expected_map, cleanuptbl map[string]interface{}
	var url, url_body_json string

	url = "/openconfig-test-xfmr:test-xfmr/test-sensor-groups/test-sensor-group[id=sensor_group_01]/config/group-colors"
	url_body_json = "{ \"openconfig-test-xfmr:group-colors\": [ \"red\",\"black\" ]}"
	pre_req_map = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_group_01": map[string]interface{}{
		"colors@": "red,green"}}}
	expected_map = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_group_01": map[string]interface{}{"colors@": "red,green,black"}}}
	cleanuptbl = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_group_01": ""}}
	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)
	t.Run("Test patch on leaf-list.", processSetRequest(url, url_body_json, "PATCH", false, nil))
	time.Sleep(1 * time.Second)
	t.Run("Verify patch on leaf-list.", verifyDbResult(rclient, "TEST_SENSOR_GROUP|sensor_group_01", expected_map, false))
	unloadDB(db.ConfigDB, cleanuptbl)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n+++++++++++++ Done Performing Patch/Update on Yang leaf-list Node demonstrating leaf-list contents merge ++++++++++++")

	fmt.Println("\n\n+++++++++++++ Performing Put/Replace on Yang leaf-list Node demonstrating leaf-list contents swap ++++++++++++")
	url = "/openconfig-test-xfmr:test-xfmr/test-sensor-groups/test-sensor-group[id=sensor_group_01]/config/group-colors"
	url_body_json = "{ \"openconfig-test-xfmr:group-colors\": [ \"blue\",\"yellow\" ]}"
	pre_req_map = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_group_01": map[string]interface{}{
		"colors@": "red,green"}}}
	expected_map = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_group_01": map[string]interface{}{"colors@": "blue,yellow"}}}
	cleanuptbl = map[string]interface{}{"TEST_SENSOR_GROUP": map[string]interface{}{"sensor_group_01": ""}}
	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)
	t.Run("Test replace on leaf-list.", processSetRequest(url, url_body_json, "PUT", false, nil))
	time.Sleep(1 * time.Second)
	t.Run("Verify replace on leaf-list.", verifyDbResult(rclient, "TEST_SENSOR_GROUP|sensor_group_01", expected_map, false))
	unloadDB(db.ConfigDB, cleanuptbl)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n+++++++++++++ Done Performing Put/Replace on Yang leaf-list Node demonstrating leaf-list contents swap ++++++++++++")

	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)

	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)
}
