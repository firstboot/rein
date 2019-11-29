package main

import (
	"fmt"
	"log"
	"os"
)

var commandDealConfPath = ""

type commandDealObj struct {
	sysServPort string
}

func commandDeal() commandDealObj {
	return commandDealObj{"19090"}
}

func (obj commandDealObj) parse(args []string) {
	if len(args) <= 1 {
		infoStr := appInfo().aboutInfo
		fmt.Println(infoStr)
		os.Exit(0)
	}

	if args[1] == "-v" {
		fmt.Println(appInfo().versionInfo)
		os.Exit(0)
	}

	if args[1] == "-h" {
		helpStr := appInfo().helpInfo
		fmt.Println(helpStr)
		os.Exit(0)
	}

	if args[1] == "-e" {
		exampleStr := appInfo().exampleInfo
		fmt.Println(exampleStr)
		os.Exit(0)
	}

	if args[1] == "-e-detail" {
		text := appInfo().exampleDetailInfo
		fmt.Println(text)
		os.Exit(0)
	}

	if args[1] == "-e-detail-upstream" {
		text := appInfo().exampleDetailUpstreamInfo
		fmt.Println(text)
		os.Exit(0)
	}

	if args[1] == "-e-detail-fileshare" {
		text := appInfo().exampleDetailFileshareInfo
		fmt.Println(text)
		os.Exit(0)
	}

	if args[1] == "-e-detail-inps" {
		text := appInfo().exampleDetailInpsInfo
		fmt.Println(text)
		os.Exit(0)
	}

	if args[1] == "-e-detail-inpc" {
		text := appInfo().exampleDetailInpcInfo
		fmt.Println(text)
		os.Exit(0)
	}

	if args[1] == "-c" {
		if len(args) != 3 {
			fmt.Println("error: incorrect number of arguments!")
			os.Exit(0)
		}
		confPath := args[2]
		commandDealConfPath = confPath
		confMap := utilsConf().getConf(confPath)
		obj.confInnerDeal(confMap)
		return
	}

	if args[1] == "-inpq" {
		if len(args) != 3 {
			fmt.Println("error: incorrect number of arguments!")
			os.Exit(0)
		}
		inpsAddr := args[2]
		coroutineInpq().run(inpsAddr)
		return
	}

	if args[1] == "-inpq-detail" {
		text := appInfo().inpqDetail
		fmt.Println(text)
		os.Exit(0)
	}

	fmt.Println("error: incorrect arguments!")
	infoStr := appInfo().aboutInfo
	fmt.Println(infoStr)
	os.Exit(0)
}

func (obj commandDealObj) confInnerDeal(confMap map[string]interface{}) {
	if utilsConf().isExistKeyOfMap("fileshare", confMap) == true {
		obj.modelFileShareDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("upstream", confMap) == true {
		obj.modelStreamDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("http", confMap) == true {
		obj.modelHTTPDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("system", confMap) == true {
		obj.systemConfDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("inps", confMap) == true {
		obj.modelInpsDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("inpc", confMap) == true {
		obj.modelInpcDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("inpcc", confMap) == true {
		obj.modelInpccDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("inpss", confMap) == true {
		obj.modelInpssDeal(confMap)
	}

	if utilsConf().isExistKeyOfMap("inpcs", confMap) == true {
		obj.modelInpcsDeal(confMap)
	}
}

func (obj commandDealObj) modelInpcsDeal(confMap map[string]interface{}) {
	for k, v := range confMap["inpcs"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		ctrlAddr := v.(map[string]interface{})["ctrl"].(string)
		fmt.Println(ctrlAddr)
		go coroutineInpcs().run(ctrlAddr)
	}
}

func (obj commandDealObj) modelInpssDeal(confMap map[string]interface{}) {
	for k, v := range confMap["inpss"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		ctrlAddr := v.(map[string]interface{})["ctrl"].(string)
		fmt.Println(ctrlAddr)
		go coroutineInpss().run(ctrlAddr)
	}
}

func (obj commandDealObj) modelInpccDeal(confMap map[string]interface{}) {
	for k, v := range confMap["inpcc"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		ctrlAddr := v.(map[string]interface{})["ctrl"].(string)
		source := v.(map[string]interface{})["source"].(string)
		target := v.(map[string]interface{})["target"].(string)
		fmt.Println(ctrlAddr, source, target)
		go coroutineInpc().run(ctrlAddr, source, target)
	}
}

func (obj commandDealObj) modelInpsDeal(confMap map[string]interface{}) {
	for k, v := range confMap["inps"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		ctrlAddr := v.(map[string]interface{})["ctrl"].(string)
		fmt.Println(ctrlAddr)
		go coroutineInps().run(ctrlAddr)
	}
}

func (obj commandDealObj) modelInpcDeal(confMap map[string]interface{}) {
	for k, v := range confMap["inpc"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		ctrlAddr := v.(map[string]interface{})["ctrl"].(string)
		source := v.(map[string]interface{})["source"].(string)
		target := v.(map[string]interface{})["target"].(string)
		fmt.Println(ctrlAddr, source, target)
		go coroutineInpc().run(ctrlAddr, source, target)
	}
}

func (obj commandDealObj) systemConfDeal(confMap map[string]interface{}) {
	systemMap := confMap["system"]
	port := systemMap.(map[string]interface{})["port"].(string)
	obj.sysServPort = port
	username := systemMap.(map[string]interface{})["username"].(string)
	password := systemMap.(map[string]interface{})["password"].(string)
	log.Println(port, username, password)
	go httpServer(port, username, password).run()

}

func (obj commandDealObj) modelFileShareDeal(confMap map[string]interface{}) {
	for k, v := range confMap["fileshare"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		port := v.(map[string]interface{})["port"].(string)
		path := v.(map[string]interface{})["path"].(string)
		fmt.Println(port, path)
		go coroutineFileShare().run(port, path)
	}
}

func (obj commandDealObj) modelStreamDeal(confMap map[string]interface{}) {
	for k, v := range confMap["upstream"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		source := v.(map[string]interface{})["source"].(string)
		target := v.(map[string]interface{})["target"].(string)
		fmt.Println(source, target)
		go coroutineStream().run(source, target)
	}
}

func (obj commandDealObj) modelHTTPDeal(confMap map[string]interface{}) {

	for k, v := range confMap["http"].([]interface{}) {
		// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
		fmt.Println(k, v)
		source := v.(map[string]interface{})["source"].(string)
		locations := v.(map[string]interface{})["locations"]
		// fmt.Println(source, locations)

		rootTargetMap := make(map[string]string) // root:target
		// rootFilterOldMap := make(map[string]string) // root-(old:inx):old
		// rootFilterNewMap := make(map[string]string) // root-(old:inx):new

		for _, lv := range locations.([]interface{}) {
			root := lv.(map[string]interface{})["root"].(string)
			target := lv.(map[string]interface{})["target"].(string)
			rootTargetMap[root] = target
			// filters := lv.(map[string]interface{})["filters"]
			// for fk, fv := range filters.([]interface{}) {
			// 	old := fv.(map[string]interface{})["old"].(string)
			// 	new := fv.(map[string]interface{})["new"].(string)
			// 	// fmt.Println(source, lk, root, target, fk, old, new)
			// 	// encodeString := base64.StdEncoding.EncodeToString(root)
			// 	// decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
			// 	rootFilterOldMap[root+"-"+strconv.Itoa(fk)] = old
			// 	rootFilterNewMap[root+"-"+strconv.Itoa(fk)] = new
			// }
		}

		// log.Println(source, rootTargetMap, rootFilterOldMap, rootFilterOldMap)

		// go coroutineHTTP(rootTargetMap, rootFilterOldMap, rootFilterOldMap).run(source)
		go coroutineHTTP(rootTargetMap).run(source)
	}
}
