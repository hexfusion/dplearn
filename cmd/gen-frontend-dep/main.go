package main

import (
	"bytes"
	"flag"
	"strings"
	"text/template"
	"time"

	"github.com/gyuho/dplearn/pkg/fileutil"
	"github.com/gyuho/dplearn/pkg/gcp"

	"github.com/golang/glog"
)

func main() {
	outputPathPackageJSON := flag.String("output-package-json", "package.json", "Specify package.json output file path.")
	outputPathAngularCLIJSON := flag.String("output-angular-cli-json", "package.json", "Specify angular-cli.json output file path.")
	flag.Parse()

	cfg := configuration{
		NgCommandServeStart:     "ng serve --aot",
		NgCommandServeStartProd: "ng serve --aot --prod",

		// 0.0.0.0 means "all IPv4 addresses on the local machine".
		// If a host has two IP addresses, 192.168.1.1 and 10.1.2.1,
		// and a server running on the host listens on 0.0.0.0,
		// it will be reachable at both of those IPs
		// (Source https://en.wikipedia.org/wiki/0.0.0.0).
		Host:         "0.0.0.0",
		HostPort:     4200,
		HostProd:     "0.0.0.0",
		HostProdPort: 4200,
	}

	bts, err := gcp.GetComputeMetadata("instance/network-interfaces/0/access-configs/0/external-ip", 3, 300*time.Millisecond)
	if err != nil {
		glog.Warning(err)
	} else {
		ip := strings.TrimSpace(string(bts))
		glog.Infof("found public host IP %q", ip)

		// TODO: angular-cli does not work with public IP, so need to use 0.0.0.0
		// https://github.com/angular/angular-cli/issues/2587#issuecomment-252586913
		// https://github.com/webpack/webpack-dev-server/issues/882
	}

	buf := new(bytes.Buffer)
	tp := template.Must(template.New("tmplPackageJSON").Parse(tmplPackageJSON))
	if err := tp.Execute(buf, &cfg); err != nil {
		glog.Fatal(err)
	}
	txt := buf.Bytes()
	if err := fileutil.WriteToFile(*outputPathPackageJSON, txt); err != nil {
		glog.Fatal(err)
	}
	glog.Infof("wrote %q", *outputPathPackageJSON)

	buf.Reset()

	tp = template.Must(template.New("tmplAngularCLIJSON").Parse(tmplAngularCLIJSON))
	if err := tp.Execute(buf, &cfg); err != nil {
		glog.Fatal(err)
	}
	txt = buf.Bytes()
	if err := fileutil.WriteToFile(*outputPathAngularCLIJSON, txt); err != nil {
		glog.Fatal(err)
	}
	glog.Infof("wrote %q", *outputPathAngularCLIJSON)
}

type configuration struct {
	NgCommandServeStart     string
	NgCommandServeStartProd string
	Host                    string
	HostPort                int
	HostProd                string
	HostProdPort            int
}

const tmplPackageJSON = `{
    "name": "app-dplearn",
    "version": "0.9.9",
    "license": "Apache-2.0",
    "angular-cli": {},
    "bin": {
        "tslint": "./bin/tslint"
    },
    "scripts": {
        "start": "{{.NgCommandServeStart}} --port {{.HostPort}} --host {{.Host}}",
        "start-prod": "{{.NgCommandServeStartProd}} --port {{.HostProdPort}} --host {{.HostProd}} --disable-host-check",
        "lint": "tslint \"frontend/**/*.ts\"",
        "test": "ng test",
        "pree2e": "webdriver-manager update",
        "e2e": "protractor"
    },
    "private": true,
    "dependencies": {
        "@angular/common": "5.2.1",
        "@angular/compiler": "5.2.1",
        "@angular/compiler-cli": "5.2.1",
        "@angular/core": "5.2.1",
        "@angular/forms": "5.2.1",
        "@angular/http": "5.2.1",
        "@angular/platform-browser": "5.2.1",
        "@angular/platform-browser-dynamic": "5.2.1",
        "@angular/animations": "5.2.1",
        "@angular/router": "5.2.1",
        "@angular/tsc-wrapped": "4.4.6",
        "@angular/upgrade": "5.2.1",
        "@angular/cli": "1.7.0-beta.0",
        "@angular/cdk": "5.0.4",
        "@angular/material": "5.0.4",
        "@types/angular": "1.6.40",
        "@types/angular-animate": "1.5.9",
        "@types/angular-cookies": "1.4.5",
        "@types/angular-mocks": "1.5.11",
        "@types/angular-resource": "1.5.14",
        "@types/angular-route": "1.3.4",
        "@types/angular-sanitize": "1.3.7",
        "@types/node": "9.3.0",
        "@types/hammerjs": "2.0.35",
        "@types/jasmine": "2.8.2",
        "core-js": "2.5.3",
        "rxjs": "5.5.6",
        "typescript": "2.6.2",
        "ts-node": "4.0.2",
        "ts-helpers": "1.1.2",
        "zone.js": "0.8.19",
        "@types/hammerjs": "2.0.35",
        "@types/jasmine": "2.8.3",
        "core-js": "2.5.3",
        "rxjs": "5.5.6",
        "typescript": "2.6.2",
        "ts-node": "4.1.0",
        "ts-helpers": "1.1.2",
        "zone.js": "0.8.19"
    },
    "devDependencies": {
        "codelyzer": "4.0.2",
        "jasmine-core": "2.8.0",
        "jasmine-spec-reporter": "4.2.1",
        "karma": "2.0.0",
        "karma-chrome-launcher": "2.2.0",
        "karma-cli": "1.0.1",
        "karma-jasmine": "1.1.1",
        "karma-remap-istanbul": "0.6.0",
        "protractor": "5.2.2",
        "tslint": "5.9.1"
    },
    "description": "website",
    "main": "index.js",
    "repository": {
        "url": "https://github.com/gyuho/dplearn",
        "type": "git"
    },
    "author": "Gyu-Ho Lee <gyuhox@gmail.com>"
}
`

const tmplAngularCLIJSON = `{
    "project": {
        "version": "1.7.0-beta.0",
        "name": "app-dplearn"
    },
    "apps": [{
        "root": "frontend",
        "outDir": "dist",
        "assets": [
            "assets",
            "favicon.ico"
        ],
        "index": "index.html",
        "main": "main.ts",
        "test": "test.ts",
        "tsconfig": "tsconfig.json",
        "prefix": "app",
        "mobile": false,
        "styles": [
            "styles.css",
            "app-dplearn-theme.scss"
        ],
        "scripts": [],
        "environmentSource": "environments/environment.ts",
        "environments": {
            "prod": "environments/environment.prod.ts",
            "dev": "environments/environment.dev.ts"
        }
    }],
    "addons": [],
    "packages": [],
    "e2e": {
        "protractor": {
            "config": "./protractor.conf.js"
        }
    },
    "test": {
        "karma": {
            "config": "./karma.conf.js"
        }
    },
    "defaults": {
        "styleExt": "css",
        "prefixInterfaces": false,
        "lazyRoutePrefix": "+",
        "serve": {
            "proxyConfig": "proxy.config.json"
        }
    }
}
`
