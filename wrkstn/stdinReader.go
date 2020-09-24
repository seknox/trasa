package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var count int

// read Creates a new buffered I/O reader and reads messages from Stdin.
func read() {

	v := bufio.NewReader(os.Stdin)
	// adjust buffer size to accommodate your json payload size limits; default is 4096
	s := bufio.NewReaderSize(v, bufferSize)

	getMsgSize := make([]byte, 4)
	msgLengthSize := int(0)

	// we're going to indefinitely read the first 4 bytes in buffer, which gives us the message length.
	// if stdIn is closed we'll exit the loop and shut down host
	for b, err := s.Read(getMsgSize); b > 0 && err == nil; b, err = s.Read(getMsgSize) {
		// convert message length bytes to integer value
		msgLengthSize = readMessageLength(getMsgSize)

		// If message length exceeds size of buffer, the message will be truncated.
		// This will likely cause an error when we attempt to unmarshal message to JSON.
		// if msgLengthSize > bufferSize {
		// 	logrus.Errorf("Message size of %d exceeds buffer size of %d.", msgLengthSize, bufferSize)
		// }

		// read the content of the message from buffer
		content := make([]byte, msgLengthSize)
		_, err := io.ReadFull(s, content)
		//_, err := s.Read(content)
		if err != nil && err != io.EOF {
			logrus.Errorf("Failed to read content: %v", err)
		}

		// message has been read, now parse and process
		parseMessage(content)
	}

}

// readMessageLength reads and returns the message length value in native byte order.
func readMessageLength(msg []byte) int {
	var length uint32
	buf := bytes.NewBuffer(msg)
	err := binary.Read(buf, nativeEndian, &length)
	if err != nil {
		logrus.Errorf("Unable to read bytes representing message length: %v", err)
	}
	return int(length)
}

// IncomingMessage represents a message sent to the native host.
type IncomingMessage struct {
	Intent string      `json:"intent"`
	Data   browserData `json:"data"`
}

// parseMessage parses incoming message
func parseMessage(msg []byte) {
	iMsg := decodeMessage(msg)
	// start building outgoing json message
	oMsg := OutgoingMessage{
		Status: true,
		Data:   "",
	}

	switch iMsg.Intent {
	case "getHygiene":

		resp, err := getEncryptedHygiene(iMsg)
		if err != nil {
			oMsg.Intent = iMsg.Intent
			oMsg.Status = false
			oMsg.Data = resp

		} else {
			oMsg.Data = resp
		}

	case "enrolDevice":

		resp, err := enrolOrSyncDevice(iMsg, "enrolDevice")
		if err != nil {
			logrus.Error(err)
			oMsg.Status = false
			oMsg.Data = resp

		} else {
			oMsg.Data = resp
		}

	case "syncDevice":

		resp, err := enrolOrSyncDevice(iMsg, "syncDevice")
		if err != nil {
			oMsg.Status = false
			oMsg.Data = resp

		} else {
			oMsg.Data = resp
		}

	default:
		oMsg.Data = "{}"
	}

	send(oMsg)
}

// decodeMessage unmarshals incoming json request and returns query value.
func decodeMessage(msg []byte) IncomingMessage {
	var iMsg IncomingMessage
	err := json.Unmarshal(msg, &iMsg)
	if err != nil {
		logrus.Errorf("Unable to unmarshal json to struct: %v. data: %s", err, string(msg))
	}
	return iMsg
}

// send sends an OutgoingMessage to os.Stdout.
func send(msg OutgoingMessage) {
	byteMsg := dataToBytes(msg)
	writeMessageLength(byteMsg)

	var msgBuf bytes.Buffer
	_, err := msgBuf.Write(byteMsg)
	if err != nil {
		logrus.Errorf("Unable to write message length to message buffer: %v", err)
	}

	_, err = msgBuf.WriteTo(os.Stdout)
	if err != nil {
		logrus.Errorf("Unable to write message buffer to Stdout: %v", err)
	}
}

// OutgoingMessage respresents a response to an incoming message query.
type OutgoingMessage struct {
	Intent string `json:"intent"`
	Status bool   `json:"status"`
	Data   string `json:"data"`
}
