// Copyright © 2016 Yieldbot <devops@yieldbot.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	gph "github.com/awalterschulze/gographviz"
	"github.com/spf13/viper"
	"github.com/yieldbot/sensuplugin/sensuutil"
	"github.com/yieldbot/vaultVisualize/version"
	"os"
	pth "path"
	"strconv"
	"strings"
)

// debugOut prints user defined and derived variable values including secrets.
// When this flag is set no connection will be made to vault, values will just
// be calculated.
func debugOut(vaultUrl string) {
	txtlogLog.WithFields(logrus.Fields{
		"Consul Tag":        tag,
		"version":           version.AppVersion(),
		"Consul Datacenter": dc,
		"Port":              port,
		"Path":              path,
		"Full Url":          vaultUrl,
	}).Info()
	sensuutil.Exit("DEBUG")
}

// buildUrl returns the url string of the vault server based upon values given
// on the commandline or specified in a configuration file.
func buildUrl() string {
	if (dc == "" || port == "" || tag == "") && !debug {
		syslogLog.WithFields(logrus.Fields{
			"host":       host,
			"app":        "vaultVisualize",
			"version":    version.AppVersion(),
			"datacenter": dc,
			"port":       port,
		}).Error(`Missing config variable, check current values`)

		txtlogLog.WithFields(logrus.Fields{
			"host":       host,
			"app":        "vaultVisualize",
			"version":    version.AppVersion(),
			"datacenter": dc,
			"port":       port,
		}).Error(`Missing config variable, check current values`)
		sensuutil.Exit("CONFIGERROR")
	}
	u := "https://" + tag + ".vault.service." + dc + ".consul:" + port
	return u
}

// setAuth returns the vault token from either the commandline or viper. Viper
// will search either a configuration file or an environment variable.
func setAuth() string {
	t := ""
	if token != "" {
		t = token
	} else if viper.Get("token") != nil {
		t = viper.Get("token").(string)
	}
	return t
}

// setCertMode will determine if the ssl certificates should be checked, this will
// default to false.
func setCertMode() string {

	//fmt.Println(viper.Get("skip_verify")) ADD TO DEBUG MODE
	//fmt.Println(strconv.FormatBool(insecureMode)) ADD TO DEBUG MODE
	if insecureMode {
		return strconv.FormatBool(insecureMode)
	} else if viper.Get("skip_verify") != nil {
		return viper.Get("skip_verify").(string)
	}
	return strconv.FormatBool(false)
}

// crawl will iterate over the vault keyspace starting at the root.
func crawl(root *secret, path string) {
	//fmt.Printf("root is %s\n", root.path) ADD TO DEBUG
	var s secret
	if root.path == "secret" {
		s = secret{
			path:     pth.Join(root.path, path),
			parent:   root,
			children: nil,
		}
	} else {
		s = secret{
			path:     path,
			parent:   root,
			children: nil,
		}
	}
	//fmt.Printf("crawling %s\n", s.path) ADD TO DEBUG
	crawler.Lock()
	secrets[s.path] = &s
	root.addChild(&s)
	sec, _ := crawler.client.Logical().List(s.path)
	keys := sec.Data["keys"]
	//spew.Dump(sec) ADD TO DEBUG
	crawler.Unlock()
	if keys == nil {
		return
	}
	if len(keys.([]interface{})) == 0 {
		fmt.Println("No additional keys")
		return
	}
	for _, key := range keys.([]interface{}) {
		crawl(&s, fmt.Sprintf("%s/%s", s.path, strings.Trim(key.(string), "/")))
	}
	return
}

// addChild adds the new child node to the mapt
func (s *secret) addChild(child *secret) {
	s.children = append(s.children, child)
}

// outputSTD prints a human readable to stdout
func outputSTD(node *secret) {
	for _, child := range node.children {
		fmt.Println(child.path)
		outputSTD(child)
	}
}

// output to a file
func outputFile(node *secret, outFile string) {
	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		syslogLog.WithFields(logrus.Fields{
			"host":    host,
			"app":     "vaultVisualize",
			"version": version.AppVersion(),
			"file":    outFile,
			"error":   err,
		}).Error(`Could not create or open the specified output file`)

		txtlogLog.WithFields(logrus.Fields{
			"host":    host,
			"app":     "vaultVisualize",
			"version": version.AppVersion(),
			"file":    outFile,
			"error":   err,
		}).Error(`Could not create or open the specified output file`)
		sensuutil.Exit("GENERALGOLANGERROR")
	}

	defer f.Close()

	for _, child := range node.children {
		if _, err = f.WriteString(child.path + "\n"); err != nil {
			syslogLog.WithFields(logrus.Fields{
				"host":    host,
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"file":    outFile,
				"error":   err,
			}).Error(`Could not write to the specified output file`)

			txtlogLog.WithFields(logrus.Fields{
				"host":    host,
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"file":    outFile,
				"error":   err,
			}).Error(`Could not write to the specified output file`)
			sensuutil.Exit("GENERALGOLANGERROR")
		}

		outputFile(child, outFile)
	}

}

func stringParse(s string) string {
	return strings.Replace(s, "-", "_", -1)
}

func colorPick(i int) string {

	colorMap[0] = "\"red\""
	colorMap[1] = "\"blue\""
	colorMap[2] = "\"green\""
	colorMap[3] = "\"yellow\""
	colorMap[4] = "\"orange\""
	colorMap[5] = "\"purple\""
	colorMap[6] = "\"brown\""

	return colorMap[i]
}

func graphOut(g *gph.Graph, node *secret, rPath string, ct int) {
	pSplit := strings.Split(node.path, "/")
	lastEl := pSplit[len(pSplit)-1]
	lastEl = stringParse(lastEl)
	//params["color"] = colorPick(ct)
	params["style"] = "\"bold\""
	//
	//fmt.Println("Color:" + colorPick(ct))
	//fmt.Println("Node:" + lastEl)
	//fmt.Println(params)

	g.AddNode("Vault", rPath, params)
	g.AddNode("Vault", lastEl, params)
	g.AddEdge(rPath, lastEl, true, nil)
	rPath = lastEl
	for _, child := range node.children {
		ct = ct + 1
		graphOut(g, child, rPath, ct)
	}

}
