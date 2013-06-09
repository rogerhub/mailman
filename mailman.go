/**
 *  Core mailman library
 */

package mailman

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type MailmanBuffer struct {
	ReversePath string
	ForwardPaths *ExpandingBuffer
	Data string
	Host string
	State CommandState
}

type CommandState int

const (
	// With a new connection
	CS_New CommandState = iota
	// Server handshake completed, ready for mail
	CS_ReadyForMail
	// Mail command received
	CS_ReadyForRcpt
	// One or more rcpt commands received
	CS_ReadyForData
	// Data command received
	CS_Data
)

var replyCodeExpansions = map[string] string {
	"500": "Syntax error, command unrecognized",
	"501": "Syntax error in parameters or arguments",
	"502": "Syntax error, command unrecognized",
	"503": "Bad sequence of commands",
	"250": "OK",
	"221": "Service closing transmission channel",
	"354": "Start mail input; end with <CRLF>.<CRLF>",
}

func writeString (w *bufio.Writer, s string) {
	if replyCodeExpansions[s] != "" {
		s += " " + replyCodeExpansions[s]
	}
	w.WriteString(s + "\n")
	w.Flush()
}

var defaultSettings = map[string] string {
	"checkAddress": "false",
	"validAddresses" : "postmaster root@localhost",
	"domain": "localhost",
	"motd": "Greetings",
	"maxMessageSize": "10000000",
}

func getSetting (key string, settings *map[string] string) (string) {
	if (*settings)[key] == "" {
		return defaultSettings[key]
	}
	return (*settings)[key]
}

func (b *MailmanBuffer) resetMailBuffers () {
	if b.State != CS_New {
		b.State = CS_ReadyForMail
	}
	b.ReversePath = ""
	b.ForwardPaths = NewExpandingBuffer()
	b.Data = ""
}

func (b *MailmanBuffer) processCommands (conn net.Conn, commands chan string,
	quits chan byte, settings *map[string] string) (bool) {
	w := bufio.NewWriter(conn)
	command := ""
	dataMode := false
	dataBuffer := make([]string, 50)
	dataBufferIndex := 0
	stopIteration := false
	writeString(w, fmt.Sprintf("220 %s %s",
		getSetting("domain", settings), getSetting("motd", settings)))
	for {
		if stopIteration {
			break
		}
		select {
			case <- quits:
				stopIteration = true
			case command = <-commands:
				if dataMode {
					if command == "." {
						if dataBufferIndex != 0 {
							b.Data += strings.Join(dataBuffer[:dataBufferIndex], "\r\n")
						}
						go saveMail(b.ReversePath, b.ForwardPaths, b.Data, b.Host)
						b.resetMailBuffers()
						dataMode = false
						writeString(w, "250")
					} else {
						if dataBufferIndex == len(dataBuffer) {
							b.Data += strings.Join(dataBuffer, "\r\n")
							dataBufferIndex = 0
						} else {
							dataBuffer[dataBufferIndex] = command
							dataBufferIndex ++
						}
					}
				} else {
					lowercaseCommand := strings.ToLower(command)
					if strings.HasPrefix(lowercaseCommand, "quit") {
						writeString(w, "221")
						stopIteration = true
					} else if strings.HasPrefix(lowercaseCommand, "ehlo") {
						if !strings.HasPrefix(lowercaseCommand, "ehlo ") ||
							len(lowercaseCommand) < 6 {
							writeString(w, "501 syntax: ehlo hostname")
						} else {
							writeString(w, fmt.Sprintf("250-%s", 
								getSetting("domain", settings)))
							writeString(w, "250-PIPELINING")
							writeString(w, fmt.Sprintf("250 SIZE %s",
								getSetting("maxMessageSize", settings)))
							b.State = CS_ReadyForMail
							b.Host = command[5:]
						}
					} else if strings.HasPrefix(lowercaseCommand, "helo") {
						if !strings.HasPrefix(lowercaseCommand, "helo ") ||
							len(lowercaseCommand) < 6 {
							writeString(w, "501 syntax: helo hostname")
						} else {
							writeString(w, fmt.Sprintf("250 %s",
								getSetting("domain", settings)))
							b.State = CS_ReadyForMail
							b.Host = command[5:]
						}
					} else if strings.HasPrefix(lowercaseCommand, "mail") {
						if !strings.HasPrefix(lowercaseCommand, "mail from:") {
							writeString(w, "501")
						} else if (b.State != CS_ReadyForMail) {
							writeString(w, "503")
						} else {
							address := command[10:]
							if !ValidateEmail(address) {
								writeString(w, "501")
							} else {
								b.ReversePath = address
								b.State = CS_ReadyForRcpt
								writeString(w, "250")
							}
						}
					} else if strings.HasPrefix(lowercaseCommand, "rcpt") {
						if !strings.HasPrefix(lowercaseCommand, "rcpt to:") {
							writeString(w, "501")
						} else if b.State != CS_ReadyForRcpt && 
							b.State != CS_ReadyForData {
							writeString(w, "503")
						} else {
							address := command[8:]
							if !ValidateEmail(address) {
								writeString(w, "501")
							} else if "true" == getSetting("checkAddress", 
								settings) && ListContains(address, 
								getSetting("validAddresses", settings)) {
								writeString(w, "550 No such user here")
							} else {
								b.ForwardPaths.InsertItem(address)
								b.State = CS_ReadyForData
								writeString(w, "250")
							}
						}
					} else if strings.HasPrefix(lowercaseCommand, "data") {
						if b.State != CS_ReadyForData {
							writeString(w, "503")
						} else {
							dataMode = true
							writeString(w, "354")
						}
					} else if strings.HasPrefix(lowercaseCommand, "rset") {
						b.resetMailBuffers()
						writeString(w, "250")
					} else {
						writeString(w, "502")
					}
				}
		}
	}
	conn.Close()
	fmt.Printf("Connection to %s terminated by client.\n", conn.RemoteAddr())
	return true
}

func handle (b *MailmanBuffer, conn net.Conn, settings *map[string] string) (bool) {
	commands := make(chan string, 100)
	quits := make(chan byte, 2)
	go b.processCommands(conn, commands, quits, settings)
	reader := bufio.NewReader(conn)
	fullLine := ""
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			// Connection was probably closed..
			quits <- byte(1)
			break
		}
		if !isPrefix {
			commands <- string(line)
			fullLine = ""
		} else {
			fullLine += string(line)
		}
	}
	// Handle final line
	if fullLine != "" {
		commands <- fullLine
	}
	return true
}

func MailmanStart (port int, settings *map[string] string) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Connected to client %s on %s\n",
			conn.RemoteAddr(), conn.LocalAddr())
		buffer := new(MailmanBuffer)
		buffer.ForwardPaths = NewExpandingBuffer()
		go handle(buffer, conn, settings)
	}
}
