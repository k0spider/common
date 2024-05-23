package email

type EmailParameter struct {
	From     string
	To       []string
	Bcc      []string
	Cc       []string
	Subject  string
	TextBody []byte
	HtmlBody []byte
}
