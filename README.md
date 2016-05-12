# ParseMailAndGo

This utility can recive new mails from pop3 server and paring its.
You can define rules for parsing mail, if any rules is matched then will be started cmd-line.

Config file for example:

{
	"Address":"pop3_server_address:110",
	"User":"autostart",
	"Pass":"qwerty",
	"Rules":[
		{
		"Subject":".*",
		"Action":"ls"
		},
		
		{
		"Subject":"shutdown",
		"Action":"shutdown -r now"
		},
		{
		"Body":"text",
		"Subject":"exe",
		"Action":"ls -l .exe"
		}
		]

}

Tegs: Address, User and Pass define parameters for pop3 server.
Subject and Body - regexp expression (https://en.wikipedia.org/wiki/Regular_expression)
Action - this teg define cmd-line, which will be started if expression in Subject and Body was true


