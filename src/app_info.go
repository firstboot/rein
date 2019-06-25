package main

type appInfoObj struct {
	exampleInfo string
	versionInfo string
	helpInfo    string
	aboutInfo   string
}

func appInfo() appInfoObj {

	const exampleInfo = `{
		"stream": [
			{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
		],
		"fileshare": [
			{"port": "9991", "path": "."}
		]
}`

	const versionInfo = `version: rein 1.0.3`

	const helpInfo = `-h: help(detail see: https://github.com/firstboot/rein)
-v: version
-c: conf(eg: rein -c rein.json)
-e: conf example`

	const aboutInfo = `author:  lz
e-mail:  linzhanggeorge@gmail.com
index:   https://github.com/firstboot/rein
help:    -h`

	return appInfoObj{exampleInfo, versionInfo, helpInfo, aboutInfo}

}
