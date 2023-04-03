////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2020 Dell, Inc.                                                 //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//  http://www.apache.org/licenses/LICENSE-2.0                                //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

//go:build !campus_pkg

package transformer

import (
	"strconv"
	"strings"

	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	log "github.com/golang/glog"
)

func init() {
	XlateFuncBind("YangToDb_persistent_subscription_key_xfmr", YangToDb_persistent_subscription_key_xfmr)
	XlateFuncBind("DbToYang_persistent_subscription_key_xfmr", DbToYang_persistent_subscription_key_xfmr)
	XlateFuncBind("YangToDb_sensor_profile_key_xfmr", YangToDb_sensor_profile_key_xfmr)
	XlateFuncBind("DbToYang_sensor_profile_key_xfmr", DbToYang_sensor_profile_key_xfmr)
	XlateFuncBind("YangToDb_subscr_destination_key_xfmr", YangToDb_subscr_destination_key_xfmr)
	XlateFuncBind("DbToYang_subscr_destination_key_xfmr", DbToYang_subscr_destination_key_xfmr)
	XlateFuncBind("YangToDb_sensor_group_key_xfmr", YangToDb_sensor_group_key_xfmr)
	XlateFuncBind("DbToYang_sensor_group_key_xfmr", DbToYang_sensor_group_key_xfmr)
	XlateFuncBind("YangToDb_sensor_path_key_xfmr", YangToDb_sensor_path_key_xfmr)
	XlateFuncBind("DbToYang_sensor_path_key_xfmr", DbToYang_sensor_path_key_xfmr)
	XlateFuncBind("YangToDb_destination_group_key_xfmr", YangToDb_destination_group_key_xfmr)
	XlateFuncBind("DbToYang_destination_group_key_xfmr", DbToYang_destination_group_key_xfmr)
	XlateFuncBind("YangToDb_destination_key_xfmr", YangToDb_destination_key_xfmr)
	XlateFuncBind("DbToYang_destination_key_xfmr", DbToYang_destination_key_xfmr)
	XlateFuncBind("YangToDb_name_field_xfmr", YangToDb_name_field_xfmr)
	XlateFuncBind("DbToYang_name_field_xfmr", DbToYang_name_field_xfmr)
	XlateFuncBind("YangToDb_sensor_group_field_xfmr", YangToDb_sensor_group_field_xfmr)
	XlateFuncBind("DbToYang_sensor_group_field_xfmr", DbToYang_sensor_group_field_xfmr)
	XlateFuncBind("YangToDb_subscr_group_id_field_xfmr", YangToDb_subscr_group_id_field_xfmr)
	XlateFuncBind("DbToYang_subscr_group_id_field_xfmr", DbToYang_subscr_group_id_field_xfmr)
	XlateFuncBind("YangToDb_path_field_xfmr", YangToDb_path_field_xfmr)
	XlateFuncBind("DbToYang_path_field_xfmr", DbToYang_path_field_xfmr)
	XlateFuncBind("YangToDb_sensor_group_id_field_xfmr", YangToDb_sensor_group_id_field_xfmr)
	XlateFuncBind("DbToYang_sensor_group_id_field_xfmr", DbToYang_sensor_group_id_field_xfmr)
	XlateFuncBind("YangToDb_destination_port_field_xfmr", YangToDb_destination_port_field_xfmr)
	XlateFuncBind("DbToYang_destination_port_field_xfmr", DbToYang_destination_port_field_xfmr)
	XlateFuncBind("YangToDb_destination_address_field_xfmr", YangToDb_destination_address_field_xfmr)
	XlateFuncBind("DbToYang_destination_address_field_xfmr", DbToYang_destination_address_field_xfmr)
	XlateFuncBind("YangToDb_group_id_field_xfmr", YangToDb_group_id_field_xfmr)
	XlateFuncBind("DbToYang_group_id_field_xfmr", DbToYang_group_id_field_xfmr)
	XlateFuncBind("telemetry_system_post_xfmr", telemetry_system_post_xfmr)
}

var YangToDb_persistent_subscription_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error

	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")

	key = name
	log.Info("YangToDb_persistent_subscription_key_xfmr key : ", key)

	return key, err
}

var DbToYang_persistent_subscription_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key

	rmap["name"] = key
	log.Info("DbToYang_persistent_subscription_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_sensor_profile_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error

	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	group := pathInfo.Var("sensor-group")

	key = name + "|" + group
	log.Info("YangToDb_sensor_profile_key_xfmr key : ", key)

	return key, err
}

var DbToYang_sensor_profile_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["sensor-group"] = TableKeys[1]
	}
	log.Info("DbToYang_sensor_profile_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_subscr_destination_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error

	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	group := pathInfo.Var("group-id")

	key = name + "|" + group
	log.Info("YangToDb_subscr_destination_key_xfmr key : ", key)

	return key, err
}

var DbToYang_subscr_destination_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["group-id"] = TableKeys[1]
	}
	log.Info("DbToYang_subscr_destination_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_sensor_group_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error
	pathInfo := NewPathInfo(inParams.uri)
	group := pathInfo.Var("sensor-group-id")

	key = group
	log.Info("YangToDb_sensor_group_key_xfmr key : ", key)

	return key, err
}

var DbToYang_sensor_group_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key

	rmap["sensor-group-id"] = key
	log.Info("DbToYang_sensor_group_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_sensor_path_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error
	pathInfo := NewPathInfo(inParams.uri)
	group := pathInfo.Var("sensor-group-id")
	path := pathInfo.Var("path")
	key = group + "|" + path
	raw_path, _, err := XfmrRemoveXPATHPredicates(path)
	log.Infof("%+v", pathInfo)
	log.Info("inParams.uri: ", inParams.uri)
	log.Info("raw_path: ", raw_path)

	log.Info("YangToDb_sensor_path_key_xfmr key : ", key)

	return key, err
}

var DbToYang_sensor_path_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["path"] = TableKeys[1]
	}
	log.Info("DbToYang_sensor_path_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_destination_group_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error
	pathInfo := NewPathInfo(inParams.uri)
	group := pathInfo.Var("group-id")

	key = group

	log.Info("YangToDb_destination_group_key_xfmr key : ", key)

	return key, err
}

var DbToYang_destination_group_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key

	rmap["group-id"] = key
	log.Info("DbToYang_destination_group_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_destination_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var key string
	var err error
	pathInfo := NewPathInfo(inParams.uri)
	group := pathInfo.Var("group-id")
	address := pathInfo.Var("destination-address")
	port := pathInfo.Var("destination-port")

	key = group + "|" + address + "|" + port
	log.Info("YangToDb_destination_key_xfmr key : ", key)

	return key, err
}

var DbToYang_destination_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 3 {
		rmap["destination-address"] = TableKeys[1]
		rmap["destination-port"], err = strconv.ParseInt(TableKeys[2], 10, 64)
	}
	log.Info("DbToYang_destination_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)
	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	dev := (*inParams.ygRoot).(*ocbinds.Device)
	cfg := dev.TelemetrySystem.Subscriptions.PersistentSubscriptions.PersistentSubscription[name].Config

	if cfg.LocalSourceAddress == nil && cfg.Encoding == ocbinds.OpenconfigTelemetryTypes_DATA_ENCODING_METHOD_UNSET && cfg.OriginatedQosMarking == nil && cfg.Protocol == ocbinds.OpenconfigTelemetryTypes_STREAM_PROTOCOL_UNSET {

		rmap["NULL"] = "NULL"
	}
	return rmap, err
}

var DbToYang_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_sensor_group_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)
	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	group := pathInfo.Var("sensor-group")
	dev := (*inParams.ygRoot).(*ocbinds.Device)
	cfg := dev.TelemetrySystem.Subscriptions.PersistentSubscriptions.PersistentSubscription[name].SensorProfiles.SensorProfile[group].Config

	if cfg.HeartbeatInterval == nil && cfg.SampleInterval == nil && cfg.SuppressRedundant == nil {
		rmap["NULL"] = "NULL"
	}
	return rmap, err
}

var DbToYang_sensor_group_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["sensor-group"] = TableKeys[1]
	}

	return rmap, err
}

var YangToDb_subscr_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_subscr_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["group-id"] = TableKeys[1]
	}

	return rmap, err
}

var YangToDb_path_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_path_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["path"] = TableKeys[1]
	}

	return rmap, err
}

var YangToDb_sensor_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_sensor_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})

	rmap["sensor-group-id"] = inParams.key

	return rmap, err
}

var YangToDb_destination_port_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_destination_port_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["destination-port"], err = strconv.ParseInt(TableKeys[2], 10, 64)
	}

	return rmap, err
}

var YangToDb_destination_address_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_destination_address_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["destination-address"] = TableKeys[1]
	}

	return rmap, err
}

var YangToDb_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["group-id"] = inParams.key

	return rmap, err
}

func telemetry_system_post_xfmr_del(inParams XfmrParams, tblName string, key string, numOfKeys int) map[string]map[string]db.Value {
	keyList, _ := inParams.d.GetKeysPattern(&(db.TableSpec{Name: tblName}), db.Key{Comp: []string{key + "*"}})
	var keys string
	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]

	for _, k := range keyList {
		subOpMap := make(map[db.DBNum]map[string]map[string]db.Value)
		subOpMapDel := make(map[string]map[string]db.Value)
		subOpMapDel[tblName] = make(map[string]db.Value)
		if len(k.Comp) == 0 {
			continue
		}
		if numOfKeys == 2 {
			keys = k.Comp[0] + "|" + k.Comp[1]
		} else {
			keys = k.Comp[0] + "|" + k.Comp[1] + "|" + k.Comp[2]
		}
		log.Info("telemetry_system_post_xfmr_del key : ", keys)
		subOpMapDel[tblName][keys] = db.Value{}
		subOpMap[db.ConfigDB] = subOpMapDel
		inParams.subOpDataMap[DELETE] = &subOpMap
		retDbDataMap = subOpMapDel
	}
	return retDbDataMap
}

var telemetry_system_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {
	requestUri,  _, _ := XfmrRemoveXPATHPredicates(inParams.requestUri)
	pathInfo := NewPathInfo(inParams.uri)
	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]

	if inParams.oper == DELETE {
		if strings.HasSuffix(requestUri, "sensor-paths") || strings.HasSuffix(requestUri, "sensor-path") {
			senGrp := pathInfo.Var("sensor-group-id")
			retDbDataMap = telemetry_system_post_xfmr_del(inParams, "DIALOUT_SENSOR_PATH", senGrp, 2)
		}
		if strings.HasSuffix(requestUri, "destinations") || strings.HasSuffix(requestUri, "destination") {
			destGrp := pathInfo.Var("group-id")
			retDbDataMap = telemetry_system_post_xfmr_del(inParams, "DIALOUT_DESTINATION", destGrp, 3)
		}
		if strings.HasSuffix(requestUri, "sensor-profiles") || strings.HasSuffix(requestUri, "sensor-profile") {
			name := pathInfo.Var("name")
			retDbDataMap = telemetry_system_post_xfmr_del(inParams, "DIALOUT_SUBSCR_SENSOR_PROFILE", name, 2)
		}
		if strings.HasSuffix(requestUri, "destination-groups") || strings.HasSuffix(requestUri, "destination-group") {
			name := pathInfo.Var("name")
			retDbDataMap = telemetry_system_post_xfmr_del(inParams, "DIALOUT_SUBSCR_DESTINATION_GROUP", name, 2)
		}
	}
	return retDbDataMap, nil
}
