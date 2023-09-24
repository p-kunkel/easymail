<h1> Go Simple MIME </h1>
 A clean and easy-to-use package for creating MIME emails.
<hr><br><br><br>

<h2> Examples </h2><br>

Send a simple email message by native smtp client
```go
package main

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/p-kunkel/easymail"
	"github.com/p-kunkel/easymail/message"
)

func main() {
	var (
		sender     = "yourEmail@example.com"
		password = "yourPassword"

		smtpHost    = "smtp.host.com"
		smtpPort    = "587"
		smtpAddress = smtpHost + ":" + smtpPort
	)

	mail := easymail.New()

	mail.From(sender)
	mail.To("recipient@example.com")
	mail.Subject("Example of use easymail")
	mail.AppendPart(message.New("It was easy!"))

	auth := smtp.PlainAuth("", sender, password, smtpHost)
	if err := mail.SmtpSend(smtpAddress, auth); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Email Sent Successfully!")
}
```

Add more recipients
```go
	mail.To("recipient_1@example.com", "Receiver 2 <r_2@example.com>")
	mail.Cc("Receiver 3 <r_3@example.com>", "recipient_4@example.com")
	mail.Bcc("recipient_5@example.com", "recipient_6@example.com")
```

Get all recipients
```go
	mail.To("receiver@example.com")
	mail.Cc("Receiver 3 <receiver_3@example.com>")
	mail.Bcc("receiver_5@example.com")

	r := mail.Headers.GetRecipients()
```

Send HTML content
```go
    html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
</head>
<body>
	<h1>HTML</h1>
	<h2>It was easy!</h2>
</body>
</html>
`

	msg = message.New(html)
	mail.AppendPart(msg)
```

Add an attachment
```go
    // local file as an attachment
	localFile := attachment.New()
	localFile.ReadFile("path/to/your/file.txt")
	mail.AppendPart(localFile)


    // write byte and send as an attachment
	otherFile := attachment.New()
	otherFile.Write([]byte("everything what you want"))
	otherFile.ContentType = "text/plain"
	otherFile.Filename = "text_file.txt"
	mail.AppendPart(other)
```

Prepare raw MIME format ready to send
```go
	raw, err := mail.Raw()
```