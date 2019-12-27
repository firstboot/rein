package main

type appInfoObj struct {
	exampleInfo                string
	versionInfo                string
	helpInfo                   string
	aboutInfo                  string
	exampleDetailInfo          string
	exampleDetailUpstreamInfo  string
	exampleDetailFileshareInfo string
	exampleDetailInpsInfo      string
	exampleDetailInpcInfo      string
	inpqDetail                 string
}

func appInfo() appInfoObj {

	const exampleInfo = `{
	"upstream": [
		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
	],
	"fileshare": [
		{"port": "9991", "path": "."}
	]
}`

	const exampleDetailInfo = `Enter a mode, show specific example, as follow:
-e-detail-upstream
-e-detail-inps
-e-detail-inpc
-e-detail-fileshare`

	const exampleDetailUpstreamInfo = `{
	"upstream": [
		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9990"}
	]
}`

	const exampleDetailFileshareInfo = `{
	"fileshare": [
		{"port": "9990", "path": "."}
	]
}`

	const exampleDetailInpsInfo = `{
	"inps": [
		{"ctrl": "0.0.0.0:17500"}
	]
}`

	const exampleDetailInpcInfo = `{
	"inpc": [
		{
			"ctrl": "127.0.0.1:17500",
			"source": "0.0.0.0:9800",
			"target": "192.168.1.122:22"
		}
	]
}`

	const versionInfo = `version: rein 1.0.6`

	const helpInfo = `-h: help(detail see: https://github.com/firstboot/rein)
-v: version
-c: conf(eg: rein -c rein.json)
-e: simple conf example (upstream and fileshare)
-e-detail:           detail conf example
-e-detail-upstream:  upstream mode conf example
-e-detail-inps:      inps mode conf example
-e-detail-inpc:      inpc mode conf example
-e-detail-fileshare: fileshare conf example
-inpq:               query inps or inpss, get source/target pairs
-inpq-detail:        inpq example`

	const aboutInfo = `author:  lz
e-mail:  linzhanggeorge@gmail.com
index:   https://github.com/firstboot/rein
help:    -h`

	const inpqDetail = `query inps or inpss, get source/target pairs.
eg:  
    ./rein -inpq x.x.x.x:17500
`

	return appInfoObj{
		exampleInfo,
		versionInfo,
		helpInfo,
		aboutInfo,
		exampleDetailInfo,
		exampleDetailUpstreamInfo,
		exampleDetailFileshareInfo,
		exampleDetailInpsInfo,
		exampleDetailInpcInfo,
		inpqDetail}

}
