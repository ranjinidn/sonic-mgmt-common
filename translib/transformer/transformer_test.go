//////////////////////////////////////////////////////////////////////////
//
// Copyright 2020 Dell, Inc.
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

package transformer_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/go-redis/redis/v7"
	"github.com/openconfig/ygot/ytypes"

	"testing"
)

var rclient *redis.Client
var port_map map[string]interface{}
var filehandle *os.File
var ygSchema *ytypes.Schema
var rclientDBNum map[db.DBNum]*redis.Client

var loadDeviceDataMap bool
var deviceDataMap = map[string]interface{}{
	"DEVICE_METADATA": map[string]interface{}{
		"localhost": map[string]interface{}{
			"hwsku":    "Force10-S6000",
			"hostname": "sonic",
			"type":     "LeafRouter",
			"platform": "x86_64-dell_s6000_s1220-r0",
			"mac":      "00:de:11:01:06:02",
		},
	},
}

func getDBOptions(dbNo db.DBNum, isWriteDisabled bool) db.Options {
	var opt db.Options

	switch dbNo {
	case db.ApplDB, db.CountersDB, db.AsicDB:
		opt = getDBOptionsWithSeparator(dbNo, "", ":", ":", isWriteDisabled)
		break
	case db.FlexCounterDB, db.ConfigDB, db.StateDB, db.ErrorDB:
		opt = getDBOptionsWithSeparator(dbNo, "", "|", "|", isWriteDisabled)
		break
	}

	return opt
}

func getDBOptionsWithSeparator(dbNo db.DBNum, initIndicator string, tableSeparator string, keySeparator string, isWriteDisabled bool) db.Options {
	return (db.Options{
		DBNo:               dbNo,
		InitIndicator:      initIndicator,
		TableNameSeparator: tableSeparator,
		KeySeparator:       keySeparator,
		IsWriteDisabled:    isWriteDisabled,
	})
}

func getAllDbs(isGetCase bool) ([db.MaxDB]*db.DB, error) {
	var dbs [db.MaxDB]*db.DB
	var err error
	var isWriteDisabled bool

	if isGetCase {
		isWriteDisabled = true
	} else {
		isWriteDisabled = false
	}

	//Create Application DB connection
	dbs[db.ApplDB], err = db.NewDB(getDBOptions(db.ApplDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	//Create ASIC DB connection
	dbs[db.AsicDB], err = db.NewDB(getDBOptions(db.AsicDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	//Create Counter DB connection
	dbs[db.CountersDB], err = db.NewDB(getDBOptions(db.CountersDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	//Create Log Level DB connection
//	dbs[db.LogLevelDB], err = db.NewDB(getDBOptions(db.LogLevelDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	isWriteDisabled = true

	//Create Config DB connection
	dbs[db.ConfigDB], err = db.NewDB(getDBOptions(db.ConfigDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	if isGetCase {
		isWriteDisabled = true
	} else {
		isWriteDisabled = false
	}

	//Create Flex Counter DB connection
	dbs[db.FlexCounterDB], err = db.NewDB(getDBOptions(db.FlexCounterDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	//Create State DB connection
	dbs[db.StateDB], err = db.NewDB(getDBOptions(db.StateDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	//Create Error DB connection
	dbs[db.ErrorDB], err = db.NewDB(getDBOptions(db.ErrorDB, isWriteDisabled))

	if err != nil {
		closeAllDbs(dbs[:])
		return dbs, err
	}

	return dbs, err
}

// Closes the dbs, and nils out the arr.
func closeAllDbs(dbs []*db.DB) {
	for dbsi, d := range dbs {
		if d != nil {
			d.DeleteDB()
			dbs[dbsi] = nil
		}
	}
}

func TestMain(t *testing.M) {
	fmt.Println("----- Setting up transformer tests -----")
	if err := setup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up transformer testing state: %v.\n", err)
		os.Exit(1)
	}
	t.Run()
	teardown()
	os.Exit(0)
}

// setups state each of the tests uses
func setup() error {
	fmt.Println("----- Performing setup -----")
	var err error
	if ygSchema, err = ocbinds.Schema(); err != nil {
		panic("Error in getting the schema: " + err.Error())
		return err
	}

	if err := initDbConfig(); err != nil {
		return err
	}

	//Clear all tables which are used for testing
	clearDb()

	/* Prepare the Redis database. */
	prepareDb()

	return nil
}

func teardown() error {
	if rclient == nil {
		return nil
	}

	fmt.Println("----- Performing teardown -----")
	unloadConfigDB(rclient, port_map)
	if loadDeviceDataMap == true {
		unloadConfigDB(rclient, deviceDataMap)
	}
	clearDb()
	rclient.Close()
	rclient.FlushDB()
	for dbNum := range rclientDBNum {
		if rclientDBNum[dbNum] != nil {
			rclientDBNum[dbNum].Close()
			rclientDBNum[dbNum].FlushDB()
		}
	}

	return nil
}

/* Separator for keys. */
func getSeparator() string {
	return "|"
}

/* Converts JSON config to map which can be loaded to Redis */
func loadConfig(key string, in []byte) map[string]interface{} {
	var fvp map[string]interface{}

	err := json.Unmarshal(in, &fvp)
	if err != nil {
		fmt.Printf("Failed to Unmarshal %v err: %v", in, err)
	}
	if key != "" {
		kv := map[string]interface{}{}
		kv[key] = fvp
		return kv
	}
	return fvp
}

/* Unloads the Config DB based on JSON File. */
func unloadConfigDB(rclient *redis.Client, mpi map[string]interface{}) {
	for key, fv := range mpi {
		switch fv.(type) {
		case map[string]interface{}:
			for subKey, subValue := range fv.(map[string]interface{}) {
				newKey := key + getSeparator() + subKey
				_, err := rclient.Del(newKey).Result()

				if err != nil {
					fmt.Printf("Invalid data for db: %v : %v %v", newKey, subValue, err)
				}

			}
		default:
			fmt.Printf("Invalid data for db: %v : %v", key, fv)
		}
	}

}

func unloadDB(dbNum db.DBNum, mpi map[string]interface{}) {
	client := rclientDBNum[dbNum]
	opts := getDBOptions(dbNum, false)
	for key, fv := range mpi {
		switch fv.(type) {
		case map[string]interface{}:
			for subKey, subValue := range fv.(map[string]interface{}) {
				newKey := key + opts.KeySeparator + subKey
				_, err := client.Del(newKey).Result()

				if err != nil {
					fmt.Printf("Invalid data for db: %v : %v %v", newKey, subValue, err)
				}

			}
		default:
			fmt.Printf("Invalid data for db: %v : %v", key, fv)
		}
	}

}

/* Loads the Config DB based on JSON File. */
func loadConfigDB(rclient *redis.Client, mpi map[string]interface{}) {
	for key, fv := range mpi {
		switch fv.(type) {
		case map[string]interface{}:
			for subKey, subValue := range fv.(map[string]interface{}) {
				newKey := key + getSeparator() + subKey
				_, err := rclient.HMSet(newKey, subValue.(map[string]interface{})).Result()

				if err != nil {
					fmt.Printf("Invalid data for db: %v : %v %v", newKey, subValue, err)
				}

			}
		default:
			fmt.Printf("Invalid data for db: %v : %v", key, fv)
		}
	}
}

/* Loads the redis DB based on JSON File. */
func loadDB(dbNum db.DBNum, mpi map[string]interface{}) {
	client := rclientDBNum[dbNum]
	opts := getDBOptions(dbNum, false)
	for key, fv := range mpi {
		switch fv.(type) {
		case map[string]interface{}:
			for subKey, subValue := range fv.(map[string]interface{}) {
				newKey := key + opts.KeySeparator + subKey
				_, err := client.HMSet(newKey, subValue.(map[string]interface{})).Result()

				if err != nil {
					fmt.Printf("Invalid data for db: %v : %v %v", newKey, subValue, err)
				}

			}
		default:
			fmt.Printf("Invalid data for db: %v : %v", key, fv)
		}
	}
}

var dbConfig struct {
	Instances map[string]map[string]interface{} `json:"INSTANCES"`
	Databases map[string]map[string]interface{} `json:"DATABASES"`
}

func initDbConfig() error {
	dbConfigFile := "/run/redis/sonic-db/database_config.json"
	if path, ok := os.LookupEnv("DB_CONFIG_PATH"); ok {
		dbConfigFile = path
	}

	fmt.Println("dbConfigFile =", dbConfigFile)
	dbConfigJson, err := ioutil.ReadFile(dbConfigFile)
	if err == nil {
		err = json.Unmarshal(dbConfigJson, &dbConfig)
	}

	return err
}

func getConfigDbClient() *redis.Client {
	return getDbClient(int(db.ConfigDB))
}

func getDbClient(dbNum int) *redis.Client {
	addr := "localhost:6379"
	pass := ""
	for _, d := range dbConfig.Databases {
		if id, ok := d["id"]; !ok || int(id.(float64)) != dbNum {
			continue
		}

		dbi := dbConfig.Instances[d["instance"].(string)]
		addr = fmt.Sprintf("%v:%v", dbi["hostname"], dbi["port"])
		if p, ok := dbi["password_path"].(string); ok {
			pwd, _ := ioutil.ReadFile(p) //TODO handle IO error
			pass = string(pwd)
		}
		break
	}

	rclient := redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        addr,
		Password:    pass,
		DB:          dbNum,
		DialTimeout: 0,
	})
	_, err := rclient.Ping().Result()
	if err != nil {
		fmt.Printf("failed to connect to redis server %v", err)
	}
	return rclient
}

/* Prepares the database in Redis Server. */
func prepareDb() {
	rclient = getConfigDbClient()
	if rclient == nil {
		fmt.Printf("error in getConfigDbClient")
		return
	}

	rclientDBNum = make(map[db.DBNum]*redis.Client)
	rclientDBNum[db.CountersDB] = getDbClient(int(db.CountersDB))
	if rclientDBNum[db.CountersDB] == nil {
		fmt.Printf("error in getDbClient(int(db.CountersDB)")
		return
	}

	rclientDBNum[db.StateDB] = getDbClient(int(db.StateDB))
	if rclientDBNum[db.StateDB] == nil {
		fmt.Printf("error in getDbClient(int(db.StateDB)")
		return
	}

	rclientDBNum[db.ApplDB] = getDbClient(int(db.ApplDB))
	if rclientDBNum[db.ApplDB] == nil {
		fmt.Printf("error in getDbClient(int(db.ApplDB )")
		return
	}

	fileName := "testdata/port_table.json"
	PortsMapByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("read file %v err: %v", fileName, err)
	}

	//Load device data map on which application of deviation files depends
	dm, err := rclient.Keys("DEVICE_METADATA|localhost").Result()
	if (err != nil) || (len(dm) == 0) {
		loadConfigDB(rclient, deviceDataMap)
		loadDeviceDataMap = true
	}

	port_map = loadConfig("", PortsMapByte)

	portKeys, err := rclient.Keys("PORT|*").Result()
	//Load only the port config which are not there in Redis
	if err == nil {
		portMapKeys := port_map["PORT"].(map[string]interface{})
		for _, portKey := range portKeys {
			//Delete the port key which is already there in Redis
			delete(portMapKeys, portKey[len("PORTS|")-1:])
		}
		port_map["PORT"] = portMapKeys
	}

	loadConfigDB(rclient, port_map)
}

func clearDb() {

	tblList := []string{
		"ACL_RULE",
		"ACL_TABLE",
		"BGP_GLOBALS",
		"BUFFER_PG",
		"CABLE_LENGTH",
		"CFG_L2MC_TABLE",
		"INTERFACE",
		"MIRROR_SESSION",
		"PORTCHANNEL",
		"PORTCHANNEL_MEMBER",
		"PORT_QOS_MAP",
		"QUEUE",
		"SCHEDULER",
		"STP",
		"STP_PORT",
		"STP_VLAN",
		"TAM_COLLECTOR_TABLE",
		"TAM_INT_IFA_FLOW_TABLE",
		"VLAN",
		"VLAN_INTERFACE",
		"VLAN_MEMBER",
		"VRF",
		"VXLAN_TUNNEL",
		"VXLAN_TUNNEL_MAP",
		"WRED_PROFILE",
	}

	for _, tbl := range tblList {
		_, err := exec.Command("/bin/sh", "-c",
			"sonic-db-cli CONFIG_DB del `sonic-db-cli CONFIG_DB keys '"+
				tbl+"|*' | cut -d ' ' -f 2`").Output()

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
