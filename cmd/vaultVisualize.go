// Copyright © 2017 Yieldbot <devops@yieldbot.com>
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
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"github.com/yieldbot/sensuplugin/sensuutil"
	"github.com/yieldbot/vaultVisualize/version"
	"os"
	"strconv"
)

// vaultVisualizeCmd represents the vaultVisualize command
var vaultVisualizeCmd = &cobra.Command{
	Use:   "vaultVisualize",
	Short: "Visualize the vault keyspaces",
	Long:  `Walk a vault tree and present a json blob or build a graphical representation of all keyspaces.`,
	Run: func(vaultVisualize *cobra.Command, args []string) {

		// Set the baseline config for the vault client
		cfg := api.DefaultConfig()

		// configure tls for the client using viper
		tls := &api.TLSConfig{}
		c, err := strconv.ParseBool(setCertMode())
		if err != nil {
			syslogLog.WithFields(logrus.Fields{
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"error":   err,
			}).Error(`Could not set cert verify mode`)
			txtlogLog.WithFields(logrus.Fields{
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"error":   err,
			}).Error(`Could not set cert verify mode`)
			sensuutil.Exit("CONFIGERROR")
		}
		tls.Insecure = c
		cfg.ConfigureTLS(tls)

		// Set the vault server address using consul and viper
		cfg.Address = buildUrl()

		// Create a client token
		cli, err := api.NewClient(cfg)
		if err != nil {
			syslogLog.WithFields(logrus.Fields{
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"error":   err,
			}).Error(`Could not create new client`)
			txtlogLog.WithFields(logrus.Fields{
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"error":   err,
			}).Error(`Could not create new client`)
			sensuutil.Exit("GENERALGOLANGERROR")
		}

		// Set the client auth token using viper
		a := setAuth()
		if a != "" {
			cli.SetToken(a)
		} else {
			cli.SetToken(os.Getenv("VAULT_TOKEN"))
		}

		s, e := cli.Logical().List("secret")
		if e != nil {
			syslogLog.WithFields(logrus.Fields{
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"error":   err,
			}).Error(`Could not list vault keys, check the path`)
			txtlogLog.WithFields(logrus.Fields{
				"app":     "vaultVisualize",
				"version": version.AppVersion(),
				"error":   err,
			}).Error(`Could not list vault keys, check the path`)
			sensuutil.Exit("GENERALGOLANGERROR")
		}
		crawler.client = cli
		root := secret{
			path:     "secret",
			parent:   nil,
			children: nil,
		}
		secrets["secret"] = &root
		keys := s.Data["keys"]
		for _, key := range keys.([]interface{}) {
			crawl(&root, key.(string))
		}
		rPath := root.path
		graph := gph.NewGraph()
		graph.SetDir(true)
		graph.SetName("Vault")
		ct := 0
		for _, child := range root.children {
			//outputSTD(child)
			if outFile != "" {
				outputFile(child, outFile)
			}
			fmt.Println(child)
			graphOut(graph, child, rPath, ct)
			ct = ct + 1

		}
		output := graph.String()
		fmt.Println(output)

	},
}

func init() {
	RootCmd.AddCommand(vaultVisualizeCmd)
}
