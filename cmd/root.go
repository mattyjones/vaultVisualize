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
	"github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yieldbot/sensuplugin/sensuutil"
	"github.com/yieldbot/vaultVisualize/version"
	"io/ioutil"
	"log/syslog"
	"os"
)

var cfgFile string    // Configuration via Viper
var host string       // Hostname for logging
var debug bool        // debugging info
var insecureMode bool // strict cert verification
var token string      // Vault token
var dc string         // Datacenter
var port string       // Port to connect to
var path string       // Path to the secret
var tag string        // Consul tag
var outFile string    // Output file

// Create logging instances.
var syslogLog = logrus.New()
var txtlogLog = logrus.New()

// Prefix for env variables. All variables must be in the format of VAULT_FOO
const envPrefix string = "vault"

var RootCmd = &cobra.Command{
	Use:   "vaultVisualize",
	Short: fmt.Sprintf("Visualize the vault keyspaces - (%s)", version.AppVersion()),
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//Setup logging for the package.
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err != nil {
		fmt.Println(err)
		sensuutil.Exit("GENERALGOLANGERROR")
	}
	syslogLog.Hooks.Add(hook)
	syslogLog.Formatter = new(logrus.JSONFormatter)
	syslogLog.Out = ioutil.Discard // don't print anything just dump it to syslog

	// Set the hostname for use in logging within the package.
	host, err = os.Hostname()
	if err != nil {
		syslogLog.WithFields(logrus.Fields{
			"app":     "vaultVisualize",
			"version": version.AppVersion(),
			"error":   err,
		}).Error(`Could not determine the hostname of this machine as reported by the kernel.`)
		txtlogLog.WithFields(logrus.Fields{
			"app":     "vaultVisualize",
			"version": version.AppVersion(),
			"error":   err,
		}).Error(`Could not determine the hostname of this machine as reported by the kernel.`)
		sensuutil.Exit("GENERALGOLANGERROR")
	}

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./.vaultVisualize.yaml", "config file")
	RootCmd.PersistentFlags().BoolVar(&insecureMode, "insecureMode", false, "Skip cert verification (default is false)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().StringVar(&token, "token", "", "vault token")
	RootCmd.PersistentFlags().StringVar(&dc, "datacenter", "", "datacenter to connect to")
	RootCmd.PersistentFlags().StringVar(&tag, "tag", "", "consul tag to use")
	RootCmd.PersistentFlags().StringVar(&path, "path", "secret", "path to the secret w/o the leading slash")
	RootCmd.PersistentFlags().StringVar(&port, "port", "8200", "port to use")
	RootCmd.PersistentFlags().StringVar(&outFile, "outputFile", "", "file to write keysapces to")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debugging info (if any)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".vaultVisualize") // name of config file (without extension)
		viper.AddConfigPath("./")              // adding home directory as first search path
	}
	viper.SetEnvPrefix(envPrefix)
	viper.BindEnv("token")
	viper.BindEnv("skip_verify")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
