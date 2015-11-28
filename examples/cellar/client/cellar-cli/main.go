package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/raphael/goa/examples/cellar/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

// PrettyPrint is true if the tool output should be formatted for human consumption.
var PrettyPrint bool

func main() {
	// Create command line parser
	app := kingpin.New("cellar-cli", "CLI client for the cellar service (http://goa.design/getting-started.html)")
	c := client.New()
	c.UserAgent = "cellar-cli/1.0"
	app.Flag("scheme", "Set the requests scheme").Short('s').Default("http").StringVar(&c.Scheme)
	app.Flag("host", "API hostname").Short('h').Default("cellar.goa.design").StringVar(&c.Host)
	app.Flag("timeout", "Set the request timeout, defaults to 20s").Short('t').Default("20s").DurationVar(&c.Timeout)
	app.Flag("dump", "Dump HTTP request and response.").BoolVar(&c.Dump)
	app.Flag("pp", "Pretty print response body").BoolVar(&PrettyPrint)
	commands := RegisterCommands(app)
	cmdName, err := app.Parse(os.Args[1:])
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
	cmd, ok := commands[cmdName]
	if !ok {
		kingpin.Fatalf("unknown command %s", cmdName)
	}
	resp, err := cmd.Run(c)
	if err != nil {
		kingpin.Fatalf("request failed: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kingpin.Fatalf("failed to read body: %s", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Let user know if something went wrong
		var sbody string
		if len(body) > 0 {
			sbody = ": " + string(body)
		}
		fmt.Printf("error: %d%s", resp.StatusCode, sbody)
	} else if len(body) > 0 {
		var out string
		if PrettyPrint {
			var jbody interface{}
			err = json.Unmarshal(body, &jbody)
			if err != nil {
				out = string(body)
			} else {
				var b []byte
				b, err = json.MarshalIndent(jbody, "", "    ")
				if err == nil {
					out = string(b)
				} else {
					out = string(body)
				}
			}
		} else {
			out = string(body)
		}
		fmt.Print(out)
	}

	// Figure out exit code
	exitStatus := 0
	switch {
	case resp.StatusCode == 401:
		exitStatus = 1
	case resp.StatusCode == 403:
		exitStatus = 3
	case resp.StatusCode == 404:
		exitStatus = 4
	case resp.StatusCode > 399 && resp.StatusCode < 500:
		exitStatus = 2
	case resp.StatusCode > 499:
		exitStatus = 5
	}
	os.Exit(exitStatus)
}

// RegisterCommands all the resource action subcommands to the application command line.
func RegisterCommands(app *kingpin.Application) map[string]client.ActionCommand {
	res := make(map[string]client.ActionCommand)
	var command, sub *kingpin.CmdClause
	command = app.Command("create", "create action")
	tmp1 := new(CreateAccountCommand)
	sub = command.Command("account", "Create new account")
	tmp1.RegisterFlags(sub)
	tmp1.Help = "Create new account"
	tmp1.Help += "\n\t--payload=PAYLOAD JSON encoded payload"

	res["create account"] = tmp1
	tmp2 := new(CreateBottleCommand)
	sub = command.Command("bottle", "Record new bottle")
	tmp2.RegisterFlags(sub)
	tmp2.Help = "Record new bottle"
	tmp2.Help += "\n\t--payload=PAYLOAD JSON encoded payload"

	res["create bottle"] = tmp2
	command = app.Command("delete", "delete action")
	tmp3 := new(DeleteAccountCommand)
	sub = command.Command("account", "")
	tmp3.RegisterFlags(sub)
	tmp3.Help = "Runs delete on a account resource."

	res["delete account"] = tmp3
	tmp4 := new(DeleteBottleCommand)
	sub = command.Command("bottle", "")
	tmp4.RegisterFlags(sub)
	tmp4.Help = "Runs delete on a bottle resource."

	res["delete bottle"] = tmp4
	command = app.Command("list", "List all bottles in account optionally filtering by year")
	tmp5 := new(ListBottleCommand)
	sub = command.Command("bottle", "List all bottles in account optionally filtering by year")
	tmp5.RegisterFlags(sub)
	tmp5.Help = "List all bottles in account optionally filtering by year"
	tmp5.Help += "\n\t--years=PARAM Filter by years"

	res["list bottle"] = tmp5
	command = app.Command("rate", "")
	tmp6 := new(RateBottleCommand)
	sub = command.Command("bottle", "")
	tmp6.RegisterFlags(sub)
	tmp6.Help = "Runs rate on a bottle resource."
	tmp6.Help += "\n\t--payload=PAYLOAD JSON encoded payload"

	res["rate bottle"] = tmp6
	command = app.Command("show", "show action")
	tmp7 := new(ShowAccountCommand)
	sub = command.Command("account", "Retrieve account with given id")
	tmp7.RegisterFlags(sub)
	tmp7.Help = "Retrieve account with given id"

	res["show account"] = tmp7
	tmp8 := new(ShowBottleCommand)
	sub = command.Command("bottle", "Retrieve bottle with given id")
	tmp8.RegisterFlags(sub)
	tmp8.Help = "Retrieve bottle with given id"

	res["show bottle"] = tmp8
	command = app.Command("update", "update action")
	tmp9 := new(UpdateAccountCommand)
	sub = command.Command("account", "Change account name")
	tmp9.RegisterFlags(sub)
	tmp9.Help = "Change account name"
	tmp9.Help += "\n\t--payload=PAYLOAD JSON encoded payload"

	res["update account"] = tmp9
	tmp10 := new(UpdateBottleCommand)
	sub = command.Command("bottle", "")
	tmp10.RegisterFlags(sub)
	tmp10.Help = "Runs update on a bottle resource."
	tmp10.Help += "\n\t--payload=PAYLOAD JSON encoded payload"

	res["update bottle"] = tmp10

	return res
}
