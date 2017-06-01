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
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"os"
	"testing"
)

// init sets environment variables that will be used for the test and configures viper to read in
// a test config from an external file.
func init() {
	// Clear the air and set some stock testing values
	os.Setenv("TEST_VAULT_SKIP_VERIFY", "")
	os.Setenv("TEST_VAULT_TOKEN", "")
	viper.SetEnvPrefix("test_vault")
	viper.BindEnv("token")
	viper.BindEnv("skip_verify")

	viper.SetConfigName("vaultVisualizeConfig") // name of config file (without extension)
	viper.AddConfigPath("../test")              // path to look for the config file in
	err := viper.ReadInConfig()                 // Find and read the config file
	if err != nil {                             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func TestSetCertMode(t *testing.T) {

	Convey("When the insecureMode flag is not set and VAULT_SKIP_VERIFY is not set and viper is set to false", t, func() {

		os.Setenv("TEST_VAULT_SKIP_VERIFY", "")
		insecureMode = false

		i := setCertMode()

		Convey("setCertMode should return a string and the certificate mode should be secure", func() {
			So(i, should.Equal, "false")
		})

		Convey("setCertMode should return a string and the certificate mode should not be insecure", func() {
			So(i, should.NotEqual, "true")
		})

		Convey("setCertMode should not return a bool and the certificate mode should not be secure", func() {
			So(i, should.NotEqual, false)
		})

		Convey("setCertMode should not return a bool and the certificate mode should not be insecure", func() {
			So(i, should.NotEqual, true)
		})
	})

	Convey("When the insecureMode flag is set and VAULT_SKIP_VERIFY is not set and viper is false", t, func() {

		os.Setenv("TEST_VAULT_SKIP_VERIFY", "")
		insecureMode = true

		i := setCertMode()

		Convey("setCertMode should return a string and the certificate mode should be insecure", func() {
			So(i, should.Equal, "true")
		})

		Convey("setCertMode should return a string and the certificate mode should not be secure", func() {
			So(i, should.NotEqual, "false")
		})

		Convey("setCertMode should not return a bool and he certificate mode should be insecure", func() {
			So(i, should.NotEqual, true)
		})

		Convey("setCertMode should not return a bool and the certificate mode should not be secure", func() {
			So(i, should.NotEqual, false)
		})
	})

	Convey("When the insecureMode flag is not set and VAULT_SKIP_VERIFY is set to TRUE and viper is set to false", t, func() {

		os.Setenv("TEST_VAULT_SKIP_VERIFY", "true")
		insecureMode = false

		i := setCertMode()

		Convey("setCertMode should return a string and the certificate mode should be insecure", func() {
			So(i, should.Equal, "true")
		})

		Convey("setCertMode should return a string and the certificate mode should not be secure", func() {
			So(i, should.NotEqual, "false")
		})

		Convey("setCertMode should not return a bool and he certificate mode should be insecure", func() {
			So(i, should.NotEqual, true)
		})

		Convey("setCertMode should not return a bool and the certificate mode should not be secure", func() {
			So(i, should.NotEqual, false)
		})
	})
}

func TestSetAuth(t *testing.T) {

	// Goal: --token (set via the commandline) should win the day
	Convey("When TEST_VAULT_TOKEN is \"1234\", viper is \"4321\", and --token is 138", t, func() {
		os.Setenv("TEST_VAULT_TOKEN", "1234")
		token = "138"
		a := setAuth()

		Convey("The token should equal to the string \"138\"", func() {
			So(a, should.Equal, "138")
		})
		Convey("The token should not be equal to the string \"1234\"", func() {
			So(a, should.NotEqual, "1234")
		})
		Convey("The token should not be equal to the string \"4321\"", func() {
			So(a, should.NotEqual, "4321")
		})
		Convey("The token should not be equal to nil", func() {
			So(a, should.NotEqual, nil)
		})
		Convey("The token should not be equal to \"\"", func() {
			So(a, should.NotEqual, "")
		})
		Convey("The token should not be equal to integer 37", func() {
			So(a, should.NotEqual, 37)
		})
	})

	// Goal: token (set via viper) should win the day
	Convey("When TEST_VAULT_TOKEN is \"\", viper is \"4321\", and --token is not set", t, func() {
		os.Setenv("TEST_VAULT_TOKEN", "")
		token = ""
		a := setAuth()

		Convey("The token should be equal to the string \"4321\"", func() {
			So(a, should.Equal, "4321")
		})

		Convey("The token should not be equal to the string \"1234\"", func() {
			So(a, should.NotEqual, "1234")
		})
		Convey("The token should not be equal to the string \"\"", func() {
			So(a, should.NotEqual, "")
		})
		Convey("The token should not be equal to nil", func() {
			So(a, should.NotEqual, nil)
		})
		Convey("The token should not be equal to integer 37", func() {
			So(a, should.NotEqual, 37)
		})
	})

	// Goal: token (set via env var) should win the day
	Convey("When TEST_VAULT_TOKEN is \"1234\", viper is \"4321\", and --token is not set", t, func() {
		os.Setenv("TEST_VAULT_TOKEN", "1234")
		token = ""
		a := setAuth()

		Convey("The token should be equal to the string \"1234\"", func() {
			So(a, should.Equal, "1234")
		})

		Convey("The token should not be equal to the string \"4321\"", func() {
			So(a, should.NotEqual, "4321")
		})
		Convey("The token should not be equal to the string \"\"", func() {
			So(a, should.NotEqual, "")
		})
		Convey("The token should not be equal to nil", func() {
			So(a, should.NotEqual, nil)
		})
		Convey("The token should not be equal to integer 37", func() {
			So(a, should.NotEqual, 37)
		})
	})
}

func TestBulidUrl(t *testing.T) {

	Convey("When the consul tag is \"foo\", dc is \"bar\", and port is \"42\" and the debug flag is not set", t, func() {
		dc = "bar"
		tag = "foo"
		port = "42"
		u := buildUrl()

		Convey("The url should be \"https://foo.vault.service.bar.consul:42\"", func() {
			So(u, should.Equal, "https://foo.vault.service.bar.consul:42")
		})
		Convey("The url should not be \"https://bar.vault.service.foo.consul:42\"", func() {
			So(u, should.NotEqual, "https://bar.vault.service.foo.consul:42")
		})
		Convey("The url should be \"http://bar.vault.service.foo.consul:42\"", func() {
			So(u, should.NotEqual, "http://bar.vault.service.foo.consul:42")
		})
		Convey("The url should not be \"\"", func() {
			So(u, should.NotEqual, "")
		})
	})
}

func TestStringParse(t *testing.T) {

	Convey("When the tInString is foo-bar", t, func() {
		tInString := "foo-bar"
		tOutString := ""
		tOutString = stringParse(tInString)

		Convey("The output should be foo_bar", func() {
			So(tOutString, should.Equal, "foo_bar")
		})
		Convey("The output should not be foo-bar", func() {
			So(tOutString, should.NotEqual, "foo-bar")
		})
		Convey("The output should not be an empty string", func() {
			So(tOutString, should.NotEqual, "")
		})
	})
}

func TestColorPick(t *testing.T) {

	Convey("When setting the color of a node or element and the input is 0", t, func() {
		var out string
		out = colorPick(0)

		Convey("The output should contain red", func() {
			So(out, should.ContainSubstring, "red")
		})
		Convey("The output should not contain blue", func() {
			So(out, should.NotContainSubstring, "blue")
		})
		Convey("The output should not be an empty string", func() {
			So(out, should.NotEqual, "")
		})
		Convey("The output should be properly escaped", func() {
			So(out, should.Equal, "\"red\"")
		})
	})
}
