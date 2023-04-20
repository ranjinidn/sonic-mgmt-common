////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Dell, Inc.                                                 //
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

package transformer

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"

	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
)

var XlateFuncs = make(map[string]reflect.Value)

var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
)

func XlateFuncBind(name string, fn interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(name + " is not valid Xfmr function.")
		}
	}()

	if _, ok := XlateFuncs[name]; !ok {
		v := reflect.ValueOf(fn)
		v.Type().NumIn()
		XlateFuncs[name] = v
	} else {
		xfmrLogInfo("Duplicate entry found in the XlateFunc map ", name)
	}
	return
}
func IsXlateFuncBinded(name string) bool {
	if _, ok := XlateFuncs[name]; !ok {
		return false
	} else {
		return true
	}
}
func XlateFuncCall(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := XlateFuncs[name]; !ok {
		log.Warning("Xfmr function does not exist: ", name)
		return nil, nil
	}
	if len(params) != XlateFuncs[name].Type().NumIn() {
		log.Warning("Error parameters not adapted")
		return nil, nil
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = XlateFuncs[name].Call(in)
	return result, nil
}

func TraverseDb(dbs [db.MaxDB]*db.DB, spec KeySpec, result *map[db.DBNum]map[string]map[string]db.Value, parentKey *db.Key, dbTblKeyGetCache map[db.DBNum]map[string]map[string]bool) error {
	var dataMap = make(RedisDbMap)

	for i := db.ApplDB; i < db.MaxDB; i++ {
		dataMap[i] = make(map[string]map[string]db.Value)
	}

	err := traverseDbHelper(dbs, spec, &dataMap, parentKey, dbTblKeyGetCache)
	if err != nil {
		xfmrLogDebug("Didn't get all data from traverseDbHelper")
		return err
	}
	/* db data processing */
	curMap := make(map[Operation]map[db.DBNum]map[string]map[string]db.Value)
	curMap[GET] = dataMap
	err = dbDataXfmrHandler(curMap)
	if err != nil {
		log.Warning("No conversion in dbdata-xfmr")
		return err
	}

	for oper, dbData := range curMap {
		if oper == GET {
			for dbNum, tblData := range dbData {
				mapCopy((*result)[dbNum], tblData)
			}
		}
	}
	return nil
}

func traverseDbHelper(dbs [db.MaxDB]*db.DB, spec KeySpec, result *map[db.DBNum]map[string]map[string]db.Value, parentKey *db.Key, dbTblKeyGetCache map[db.DBNum]map[string]map[string]bool) error {
	var err error
	var dbOpts db.Options = getDBOptions(spec.DbNum)

	separator := dbOpts.KeySeparator

	if spec.Key.Len() > 0 {
		// get an entry with a specific key
		if spec.Ts.Name != XFMR_NONE_STRING { // Do not traverse for NONE table
			data, err := dbs[spec.DbNum].GetEntry(&spec.Ts, spec.Key)
			queriedDbInfo := make(map[string]map[string]bool)
			queriedDbTblInfo := make(map[string]bool)
			queriedDbTblInfo[strings.Join(spec.Key.Comp, separator)] = true
			queriedDbInfo[spec.Ts.Name] = queriedDbTblInfo
			if dbTblKeyGetCache == nil {
				dbTblKeyGetCache = make(map[db.DBNum]map[string]map[string]bool)
			}
			dbTblKeyGetCache[spec.DbNum] = queriedDbInfo
			if err != nil {
				log.Warningf("Didn't get data for tbl(%v), key(%v) in traverseDbHelper", spec.Ts.Name, spec.Key)
				return err
			}

			if (*result)[spec.DbNum][spec.Ts.Name] == nil {
				(*result)[spec.DbNum][spec.Ts.Name] = map[string]db.Value{strings.Join(spec.Key.Comp, separator): data}
			} else {
				(*result)[spec.DbNum][spec.Ts.Name][strings.Join(spec.Key.Comp, separator)] = data
			}
		}
		if len(spec.Child) > 0 {
			for _, ch := range spec.Child {
				err = traverseDbHelper(dbs, ch, result, &spec.Key, dbTblKeyGetCache)
			}
		}
	} else {
		// TODO - GetEntry support with regex patten, 'abc*' for optimization
		if spec.Ts.Name != XFMR_NONE_STRING { //Do not traverse for NONE table
			keys, err := dbs[spec.DbNum].GetKeys(&spec.Ts)
			if err != nil {
				log.Warningf("Didn't get keys for tbl(%v) in traverseDbHelper", spec.Ts.Name)
				return err
			}
			xfmrLogDebug("keys for table %v in DB %v are %v", spec.Ts.Name, spec.DbNum, keys)
			for i := range keys {
				if parentKey != nil && !spec.IgnoreParentKey {
					// TODO - multi-depth with a custom delimiter
					if !strings.Contains(strings.Join(keys[i].Comp, separator), strings.Join((*parentKey).Comp, separator)) {
						continue
					}
				}
				spec.Key = keys[i]
				err = traverseDbHelper(dbs, spec, result, parentKey, dbTblKeyGetCache)
				if err != nil {
					xfmrLogDebug("Traversal didn't fetch for : %v", err)
				}
			}
		} else if len(spec.Child) > 0 {
			for _, ch := range spec.Child {
				err = traverseDbHelper(dbs, ch, result, &spec.Key, dbTblKeyGetCache)
			}
		}
	}
	return err
}

func XlateUriToKeySpec(uri string, requestUri string, ygRoot *ygot.GoStruct, t *interface{}, txCache interface{}) (*[]KeySpec, error) {

	var err error
	var retdbFormat = make([]KeySpec, 0)

	// In case of SONIC yang, the tablename and key info is available in the xpath
	if isSonicYang(uri) {
		/* Extract the xpath and key from input xpath */
		xpath, keyStr, tableName := sonicXpathKeyExtract(uri)
		if tblSpecInfo, ok := xDbSpecMap[tableName]; ok && keyStr != "" && hasKeyValueXfmr(tableName) {
			/* key from URI should be converted into redis-db key, to read data */
			keyStr, err = dbKeyValueXfmrHandler(CREATE, tblSpecInfo.dbIndex, tableName, keyStr)
			if err != nil {
				log.Warningf("Value-xfmr for table(%v) & key(%v) didn't do conversion.", tableName, keyStr)
				return &retdbFormat, err
			}
		}

		retdbFormat = fillSonicKeySpec(xpath, tableName, keyStr)
	} else {
		/* Extract the xpath and key from input xpath */
		retData, _ := xpathKeyExtract(nil, ygRoot, GET, uri, requestUri, nil, nil, txCache, nil)
		retdbFormat = fillKeySpecs(retData.xpath, retData.dbKey, &retdbFormat)
	}

	return &retdbFormat, err
}

func fillKeySpecs(yangXpath string, keyStr string, retdbFormat *[]KeySpec) []KeySpec {
	var err error
	if xYangSpecMap == nil {
		return *retdbFormat
	}
	_, ok := xYangSpecMap[yangXpath]
	if ok {
		xpathInfo := xYangSpecMap[yangXpath]
		if xpathInfo.tableName != nil {
			dbFormat := KeySpec{}
			dbFormat.Ts.Name = *xpathInfo.tableName
			dbFormat.DbNum = xpathInfo.dbIndex
			if len(xYangSpecMap[yangXpath].xfmrKey) > 0 || xYangSpecMap[yangXpath].keyName != nil {
				dbFormat.IgnoreParentKey = true
			} else {
				dbFormat.IgnoreParentKey = false
			}
			if keyStr != "" {
				if tblSpecInfo, ok := xDbSpecMap[dbFormat.Ts.Name]; ok && tblSpecInfo.hasXfmrFn {
					/* key from URI should be converted into redis-db key, to read data */
					keyStr, err = dbKeyValueXfmrHandler(CREATE, dbFormat.DbNum, dbFormat.Ts.Name, keyStr)
					if err != nil {
						log.Warningf("Value-xfmr for table(%v) & key(%v) didn't do conversion.", dbFormat.Ts.Name, keyStr)
					}
				}
				dbFormat.Key.Comp = append(dbFormat.Key.Comp, keyStr)
			}
			for _, child := range xpathInfo.childTable {
				if child == dbFormat.Ts.Name {
					continue
				}
				if xDbSpecMap != nil {
					if _, ok := xDbSpecMap[child]; ok {
						chlen := len(xDbSpecMap[child].yangXpath)
						if chlen > 0 {
							children := make([]KeySpec, 0)
							for _, childXpath := range xDbSpecMap[child].yangXpath {
								children = fillKeySpecs(childXpath, "", &children)
								dbFormat.Child = append(dbFormat.Child, children...)
							}
						}
					}
				}
			}
			*retdbFormat = append(*retdbFormat, dbFormat)
		} else {
			for _, child := range xpathInfo.childTable {
				if xDbSpecMap != nil {
					if _, ok := xDbSpecMap[child]; ok {
						chlen := len(xDbSpecMap[child].yangXpath)
						if chlen > 0 {
							for _, childXpath := range xDbSpecMap[child].yangXpath {
								*retdbFormat = fillKeySpecs(childXpath, "", retdbFormat)
							}
						}
					}
				}
			}
		}
	}
	return *retdbFormat
}

func fillSonicKeySpec(xpath string, tableName string, keyStr string) []KeySpec {

	var retdbFormat = make([]KeySpec, 0)

	if tableName != "" {
		dbFormat := KeySpec{}
		dbFormat.Ts.Name = tableName
		cdb := db.ConfigDB
		if _, ok := xDbSpecMap[tableName]; ok {
			cdb = xDbSpecMap[tableName].dbIndex
		}
		dbFormat.DbNum = cdb
		if keyStr != "" {
			dbFormat.Key.Comp = append(dbFormat.Key.Comp, keyStr)
		}
		retdbFormat = append(retdbFormat, dbFormat)
	} else {
		// If table name not available in xpath get top container name
		container := xpath
		if xDbSpecMap != nil {
			if _, ok := xDbSpecMap[container]; ok {
				dbInfo := xDbSpecMap[container]
				if dbInfo.yangType == YANG_CONTAINER {
					for dir := range dbInfo.dbEntry.Dir {
						_, ok := xDbSpecMap[dir]
						if ok && xDbSpecMap[dir].yangType == YANG_CONTAINER {
							cdb := xDbSpecMap[dir].dbIndex
							dbFormat := KeySpec{}
							dbFormat.Ts.Name = dir
							dbFormat.DbNum = cdb
							retdbFormat = append(retdbFormat, dbFormat)
						}
					}
				}
			}
		}
	}
	return retdbFormat
}

func XlateToDb(path string, oper int, d *db.DB, yg *ygot.GoStruct, yt *interface{}, jsonPayload []byte, txCache interface{}, skipOrdTbl *bool) (map[Operation]RedisDbMap, map[string]map[string]db.Value, map[string]map[string]db.Value, error) {

	var err error
	requestUri := path
	jsonData := make(map[string]interface{})
	opcode := Operation(oper)

	device := (*yg).(*ocbinds.Device)
	jsonStr, _ := ygot.EmitJSON(device, &ygot.EmitJSONConfig{
		Format:         ygot.RFC7951,
		Indent:         "  ",
		SkipValidation: true,
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})

	err = json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		errStr := "Error: failed to unmarshal json."
		err = tlerr.InternalError{Format: errStr}
		return nil, nil, nil, err
	}

	// Map contains table.key.fields
	var result = make(map[Operation]RedisDbMap)
	var yangDefValMap = make(map[string]map[string]db.Value)
	var yangAuxValMap = make(map[string]map[string]db.Value)
	switch opcode {
	case CREATE:
		xfmrLogInfo("CREATE case")
		err = dbMapCreate(d, yg, opcode, path, requestUri, jsonData, result, yangDefValMap, yangAuxValMap, txCache)
		if err != nil {
			log.Warning("Data translation from YANG to db failed for create request.")
		}

	case UPDATE:
		xfmrLogInfo("UPDATE case")
		err = dbMapUpdate(d, yg, opcode, path, requestUri, jsonData, result, yangDefValMap, yangAuxValMap, txCache)
		if err != nil {
			log.Warning("Data translation from YANG to db failed for update request.")
		}

	case REPLACE:
		xfmrLogInfo("REPLACE case")
		err = dbMapUpdate(d, yg, opcode, path, requestUri, jsonData, result, yangDefValMap, yangAuxValMap, txCache)
		if err != nil {
			log.Warning("Data translation from YANG to db failed for replace request.")
		}

	case DELETE:
		xfmrLogInfo("DELETE case")
		err = dbMapDelete(d, yg, opcode, path, requestUri, jsonData, result, txCache, skipOrdTbl)
		if err != nil {
			log.Warning("Data translation from YANG to db failed for delete request.")
		}
	}
	return result, yangDefValMap, yangAuxValMap, err
}

func GetAndXlateFromDB(uri string, ygRoot *ygot.GoStruct, dbs [db.MaxDB]*db.DB, txCache interface{}) ([]byte, bool, error) {
	var err error
	var payload []byte
	var inParamsForGet xlateFromDbParams
	xfmrLogInfo("received xpath = ", uri)
	requestUri := uri

	keySpec, _ := XlateUriToKeySpec(uri, requestUri, ygRoot, nil, txCache)
	var dbresult = make(RedisDbMap)
	for i := db.ApplDB; i < db.MaxDB; i++ {
		dbresult[i] = make(map[string]map[string]db.Value)
	}

	inParamsForGet.dbTblKeyGetCache = make(map[db.DBNum]map[string]map[string]bool)

	for _, spec := range *keySpec {
		err := TraverseDb(dbs, spec, &dbresult, nil, inParamsForGet.dbTblKeyGetCache)
		if err != nil {
			xfmrLogDebug("TraverseDb() didn't fetch data.")
		}
	}

	isEmptyPayload := false
	payload, isEmptyPayload, err = XlateFromDb(uri, ygRoot, dbs, dbresult, txCache, inParamsForGet)
	if err != nil {
		return payload, true, err
	}

	return payload, isEmptyPayload, err
}

func XlateFromDb(uri string, ygRoot *ygot.GoStruct, dbs [db.MaxDB]*db.DB, data RedisDbMap, txCache interface{}, inParamsForGet xlateFromDbParams) ([]byte, bool, error) {

	var err error
	var result []byte
	var dbData = make(RedisDbMap)
	var cdb db.DBNum = db.ConfigDB
	var xpath string

	dbData = data
	requestUri := uri
	/* Check if the parent table exists for RFC compliance */
	var exists bool
	subOpMapDiscard := make(map[Operation]*RedisDbMap)
	exists, err = verifyParentTable(nil, dbs, ygRoot, GET, uri, dbData, txCache, subOpMapDiscard)
	xfmrLogDebug("verifyParentTable() returned - exists - %v, err - %v", exists, err)
	if err != nil {
		log.Warningf("Cannot perform GET Operation on URI %v due to - %v", uri, err)
		return []byte(""), true, err
	}
	if !exists {
		err = tlerr.NotFoundError{Format: "Resource Not found"}
		return []byte(""), true, err
	}

	if isSonicYang(uri) {
		lxpath, keyStr, tableName := sonicXpathKeyExtract(uri)
		xpath = lxpath
		if tableName != "" {
			dbInfo, ok := xDbSpecMap[tableName]
			if !ok {
				log.Warningf("No entry in xDbSpecMap for xpath %v", tableName)
			} else {
				cdb = dbInfo.dbIndex
			}
			tokens := strings.Split(xpath, "/")
			// Format /module:container/tableName/listname[key]/fieldName
			if tokens[SONIC_TABLE_INDEX] == tableName {
				fieldName := ""
				if len(tokens) > SONIC_FIELD_INDEX {
					fieldName = tokens[SONIC_FIELD_INDEX]
					dbSpecField := tableName + "/" + fieldName
					dbSpecFieldInfo, ok := xDbSpecMap[dbSpecField]
					if ok && fieldName != "" {
						yangNodeType := xDbSpecMap[dbSpecField].yangType
						if yangNodeType == YANG_LEAF_LIST {
							fieldName = fieldName + "@"
						}
						if (yangNodeType == YANG_LEAF_LIST) || (yangNodeType == YANG_LEAF) {
							dbData[cdb], err = extractFieldFromDb(tableName, keyStr, fieldName, data[cdb])
							// return resource not found when the leaf/leaf-list instance(not entire leaf-list GET) not found
							if (err != nil) && ((yangNodeType == YANG_LEAF) || ((yangNodeType == YANG_LEAF_LIST) && (strings.HasSuffix(uri, "]") || strings.HasSuffix(uri, "]/")))) {
								return []byte(""), true, err
							}
							if (yangNodeType == YANG_LEAF_LIST) && ((strings.HasSuffix(uri, "]")) || (strings.HasSuffix(uri, "]/"))) {
								leafListInstVal, valErr := extractLeafListInstFromUri(uri)
								if valErr != nil {
									return []byte(""), true, valErr
								}
								if dbSpecFieldInfo.xfmrValue != nil {
									inParams := formXfmrDbInputRequest(CREATE, cdb, tableName, keyStr, fieldName, leafListInstVal)
									retVal, err := valueXfmrHandler(inParams, *dbSpecFieldInfo.xfmrValue)
									if err != nil {
										log.Warningf("value-xfmr:fldpath(\"%v\") val(\"%v\"):err(\"%v\").", dbSpecField, leafListInstVal, err)
										return []byte(""), true, err
									}
									leafListInstVal = retVal
								}
								if leafListInstExists(dbData[cdb][tableName][keyStr].Field[fieldName], leafListInstVal) {
									/* Since translib already fills in ygRoot with queried leaf-list instance, do not
									   fill in resFldValMap or else Unmarshall of payload(resFldValMap) into ygotTgt in
									   app layer will create duplicate instances in result.
									*/
									log.Info("Queried leaf-list instance exists.")
									return []byte("{}"), false, nil
								} else {
									xfmrLogDebug("Queried leaf-list instance does not exist - %v", uri)
									return []byte(""), true, tlerr.NotFoundError{Format: "Resource not found"}
								}
							}
						}
					}
				}
			}
		}
	} else {
		lxpath, _, _ := XfmrRemoveXPATHPredicates(uri)
		xpath = lxpath
		if _, ok := xYangSpecMap[xpath]; ok {
			cdb = xYangSpecMap[xpath].dbIndex
		}
	}
	dbTblKeyGetCache := inParamsForGet.dbTblKeyGetCache
	inParamsForGet = formXlateFromDbParams(dbs[cdb], dbs, cdb, ygRoot, uri, requestUri, xpath, GET, "", "", &dbData, txCache, nil, false)
	inParamsForGet.xfmrDbTblKeyCache = make(map[string]tblKeyCache)
	inParamsForGet.dbTblKeyGetCache = dbTblKeyGetCache
	payload, isEmptyPayload, err := dbDataToYangJsonCreate(inParamsForGet)
	xfmrLogDebug("Payload generated : ", payload)

	if err != nil {
		log.Warning("Couldn't create json response from DB data.")
		return nil, isEmptyPayload, err
	}
	xfmrLogInfo("Created json response from DB data.")

	result = []byte(payload)
	return result, isEmptyPayload, err

}

func extractFieldFromDb(tableName string, keyStr string, fieldName string, data map[string]map[string]db.Value) (map[string]map[string]db.Value, error) {

	var dbVal db.Value
	var dbData = make(map[string]map[string]db.Value)
	var err error

	if tableName != "" && keyStr != "" && fieldName != "" {
		if data[tableName][keyStr].Field != nil {
			fldVal, fldValExists := data[tableName][keyStr].Field[fieldName]
			if fldValExists {
				dbData[tableName] = make(map[string]db.Value)
				dbVal.Field = make(map[string]string)
				dbVal.Field[fieldName] = fldVal
				dbData[tableName][keyStr] = dbVal
			} else {
				log.Warningf("Field %v doesn't exist in table - %v, instance - %v", fieldName, tableName, keyStr)
				err = tlerr.NotFoundError{Format: "Resource not found"}
			}
		}
	}
	return dbData, err
}

func GetModuleNmFromPath(uri string) (string, error) {
	xfmrLogDebug("received URI %s to extract module name from ", uri)
	moduleNm, err := uriModuleNameGet(uri)
	return moduleNm, err
}

func GetOrdTblList(xfmrTbl string, uriModuleNm string) []string {
	var ordTblList []string
	processedTbl := false
	var sncMdlList []string = getYangMdlToSonicMdlList(uriModuleNm)

	for _, sonicMdlNm := range sncMdlList {
		sonicMdlTblInfo := xDbSpecTblSeqnMap[sonicMdlNm]
		for _, ordTblNm := range sonicMdlTblInfo.OrdTbl {
			if xfmrTbl == ordTblNm {
				xfmrLogInfo("Found sonic module(%v) whose ordered table list contains table %v", sonicMdlNm, xfmrTbl)
				ordTblList = sonicMdlTblInfo.DepTbl[xfmrTbl].DepTblWithinMdl
				processedTbl = true
				break
			}
		}
		if processedTbl {
			break
		}
	}
	return ordTblList
}

func GetXfmrOrdTblList(xfmrTbl string) []string {
	/* get the table hierarchy read from json file */
	var ordTblList []string
	if _, ok := sonicOrdTblListMap[xfmrTbl]; ok {
		ordTblList = sonicOrdTblListMap[xfmrTbl]
	}
	return ordTblList
}

func GetTablesToWatch(xfmrTblList []string, uriModuleNm string) []string {
	var depTblList []string
	depTblMap := make(map[string]bool) //create to avoid duplicates in depTblList, serves as a Set
	processedTbl := false
	var sncMdlList []string
	var lXfmrTblList []string

	sncMdlList = getYangMdlToSonicMdlList(uriModuleNm)

	// remove duplicates from incoming list of tables
	xfmrTblMap := make(map[string]bool) //create to avoid duplicates in xfmrTblList
	for _, xfmrTblNm := range xfmrTblList {
		xfmrTblMap[xfmrTblNm] = true
	}
	for xfmrTblNm := range xfmrTblMap {
		lXfmrTblList = append(lXfmrTblList, xfmrTblNm)
	}

	for _, xfmrTbl := range lXfmrTblList {
		processedTbl = false
		//can be optimized if there is a way to know all sonic modules, a given OC-Yang spans over
		for _, sonicMdlNm := range sncMdlList {
			sonicMdlTblInfo := xDbSpecTblSeqnMap[sonicMdlNm]
			for _, ordTblNm := range sonicMdlTblInfo.OrdTbl {
				if xfmrTbl == ordTblNm {
					xfmrLogInfo("Found sonic module(%v) whose ordered table list contains table %v", sonicMdlNm, xfmrTbl)
					ldepTblList := sonicMdlTblInfo.DepTbl[xfmrTbl].DepTblAcrossMdl
					for _, depTblNm := range ldepTblList {
						depTblMap[depTblNm] = true
					}
					//assumption that a table belongs to only one sonic module
					processedTbl = true
					break
				}
			}
			if processedTbl {
				break
			}
		}
		if !processedTbl {
			depTblMap[xfmrTbl] = false
		}
	}
	for depTbl := range depTblMap {
		depTblList = append(depTblList, depTbl)
	}
	return depTblList
}

func CallRpcMethod(path string, body []byte, dbs [db.MaxDB]*db.DB) ([]byte, error) {
	const (
		RPC_XFMR_RET_ARGS     = 2
		RPC_XFMR_RET_VAL_INDX = 0
		RPC_XFMR_RET_ERR_INDX = 1
	)
	var err error
	var ret []byte
	var data []reflect.Value
	var rpcFunc = ""

	// TODO - check module name
	if isSonicYang(path) {
		if rpcFuncNm, ok := xDbRpcSpecMap[path]; ok {
			rpcFunc = rpcFuncNm
		}
	} else {
		if rpcFuncNm, ok := xYangRpcSpecMap[path]; ok {
			rpcFunc = rpcFuncNm
		}
	}

	if rpcFunc != "" {
		xfmrLogInfo("RPC callback invoked (%v) \r\n", rpcFunc)
		data, err = XlateFuncCall(rpcFunc, body, dbs)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			if len(data) == RPC_XFMR_RET_ARGS {
				// rpc xfmr callback returns err as second value in return data list from <xfmr_func>.Call()
				if data[RPC_XFMR_RET_ERR_INDX].Interface() != nil {
					err = data[RPC_XFMR_RET_ERR_INDX].Interface().(error)
					if err != nil {
						log.Warningf("Transformer function(\"%v\") returned error - %v.", rpcFunc, err)
					}
				}
			}

			if data[RPC_XFMR_RET_VAL_INDX].Interface() != nil {
				retVal, retOk := data[RPC_XFMR_RET_VAL_INDX].Interface().([]byte)
				if retOk {
					ret = retVal
				}
			}
		}
	} else {
		log.Warning("Not supported RPC", path)
		err = tlerr.NotSupported("Not supported RPC")
	}
	return ret, err
}

func AddModelCpbltInfo() map[string]*mdlInfo {
	return xMdlCpbltMap
}

func IsTerminalNode(uri string) (bool, error) {
	xpath, _, err := XfmrRemoveXPATHPredicates(uri)
	if xpathData, ok := xYangSpecMap[xpath]; ok {
		if !xpathData.hasNonTerminalNode {
			return true, nil
		}
	} else {
		log.Warningf("xYangSpecMap data not found for xpath : %v", xpath)
		errStr := "xYangSpecMap data not found for xpath."
		err = tlerr.InternalError{Format: errStr}
	}

	return false, err
}

func IsLeafNode(uri string) bool {
	result := false
	yngNdType, err := getYangNodeTypeFromUri(uri)
	if (err == nil) && (yngNdType == YANG_LEAF) {
		result = true
	}
	return result
}

func IsLeafListNode(uri string) bool {
	result := false
	yngNdType, err := getYangNodeTypeFromUri(uri)
	if (err == nil) && (yngNdType == YANG_LEAF_LIST) {
		result = true
	}
	return result
}

func tableKeysToBeSorted(tblNm string) bool {
	/* function to decide whether to sort table keys.
	Required when a sonic table has more than 1 lists
	with keys having leaf-refs to each other, i.e table has primary and secondary keys
	*/
	areTblKeysToBeSorted := false
	TBL_LST_CNT_NO_SEC_KEY := 1 //Tables having primary and secondary keys have more than one lists defined in sonic yang
	if dbSpecInfo, ok := xDbSpecMap[tblNm]; ok {
		if len(dbSpecInfo.listName) > TBL_LST_CNT_NO_SEC_KEY {
			areTblKeysToBeSorted = true
		}
	} else {
		log.Warning("xDbSpecMap data not found for ", tblNm)
	}
	xfmrLogInfo("Table %v keys should be sorted - %v", tblNm, areTblKeysToBeSorted)
	return areTblKeysToBeSorted
}

func SortSncTableDbKeys(tableName string, dbKeyMap map[string]db.Value) []string {
	var ordDbKey []string

	if tableKeysToBeSorted(tableName) {

		m := make(map[string]int)
		for tblKey := range dbKeyMap {
			keyList := strings.Split(tblKey, "|")
			m[tblKey] = len(keyList)
		}

		type kv struct {
			Key   string
			Value int
		}

		var ss []kv
		for k, v := range m {
			ss = append(ss, kv{k, v})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		for _, kv := range ss {
			ordDbKey = append(ordDbKey, kv.Key)
		}

	} else {

		// Restore the order as in the original map in case of single list in table case and error case
		if len(ordDbKey) == 0 {
			for tblKey := range dbKeyMap {
				ordDbKey = append(ordDbKey, tblKey)
			}
		}
	}

	return ordDbKey
}
