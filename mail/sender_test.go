package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xiao-yangg/simplebank/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() { // skip testing if short flag = True
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test Email"
	content := `
	<h1>Hello World</h1>
	<p>This is a test message</p>
	`
	to := []string{"receiver1@gmail.com", "receiver2@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}