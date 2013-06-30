// SMTP client and server program
package main

import (
	"fmt"
	"github.com/rogerhub/mailman/server"
	"github.com/rogerhub/mailman/simpleconf"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		var UsageString = "Usage: mailman [-h|--help] [-x|--exampleconf] [configfile]\n"
		fmt.Fprintf(os.Stderr, UsageString)
		return
	} else if os.Args[1] == "--help" || os.Args[1] == "-h" {
		var HelpString = ""
		fmt.Fprintf(os.Stderr, HelpString)
		return
	} else if os.Args[1] == "--exampleconf" || os.Args[1] == "-x" {
		var TemplateString = "" + 
			";;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;\n" +
			"; mm.conf Default Settings    ;\n" +
			";;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;\n" +
			"; Use this output as a template to configure mailman\n" + 
			"; Pipe it to a file by running: $ mailman -x > mm.conf\n" +
			"; Then, run $ mailman mm.conf\n" +
			"\n" +
			"[server]\n" +
			"; Which port and addresses to listen on\n" +
			"; This can be a port number like :25, a semantic port like :http, or\n" +
			"; something like localhost:7677.\n" +
			"listen=0.0.0.0:2225\n" +
			"\n" +
			"; The hostname that identifies this server\n" +
			"; Used during the initial handshake with the client. Use something sensible.\n" +
			"server_hostname=localhost\n" +
			"\n" +
			"; Greeting used at initial handshake\n" +
			"motd=Greetings from Mailman!\n" +
			"\n" +
			"[saver]\n" +
			"; How do you want to save messages?\n" +
			"; Valid options: discard, file, postgres\n" +
			"mail_saver=file\n" +
			"\n" +
			"[file_saver]\n" +
			"; The file saver just logs basic data about the email to a file. This should be\n" +
			"; used primarily for debugging purposes, whereas a database is more suitable for\n" +
			"; live environments.\n" +
			"file_saver_target=/var/tmp/mailman.out\n"
		fmt.Fprintf(os.Stdout, TemplateString)
		return
	} else {
		settings, err := simpleconf.ParseFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		server.MailmanStart(settings)
	}
}
