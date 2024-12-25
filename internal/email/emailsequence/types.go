package emailsequence

type SMTPConnectType string

const (
	SSL SMTPConnectType = "SSL"
)

type EmailAccountHandler struct {
	Name        string
	SMTPHost    SMTPServer
	SMTPPort    int
	IMAPHost    IMAPServer
	IMAPPort    int
	Username    string
	Password    string
	ConnectType SMTPConnectType
}

type SMTPServer struct {
	// Define SMTPServer fields here
}

type IMAPServer struct {
	// Define IMAPServer fields here
}

type EmailSendingSchema struct {
	To  []string `json:"to"`
	CC  string   `json:"cc,omitempty"`
	BCC string   `json:"bcc,omitempty"`
}

type SendEmailResponse struct {
	Message string `json:"message"`
}

type EmailAttachment struct {
	Filename    string `json:"filename"`
	Content     []byte `json:"content"`
	ContentType string `json:"content_type"`
}

type Email struct {
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
}

type Delay struct {
	Seconds int `json:"seconds"`
	Minutes int `json:"minutes"`
	Hours   int `json:"hours"`
	Days    int `json:"days"`
	Weeks   int `json:"weeks"`
}

type EventType string

const (
	OpenedEmail   EventType = "opened_email"
	ClickedOnLink EventType = "clicked_on_link"
)

type TriggerCondition string

const (
	Yes TriggerCondition = "yes"
	No  TriggerCondition = "no"
)

type Trigger struct {
	EventType EventType        `json:"event_type"`
	Condition TriggerCondition `json:"condition"`
}

type Goal struct {
	Description string `json:"description"`
}

type SequenceStep struct {
	Email   *Email   `json:"email,omitempty"`
	Trigger *Trigger `json:"trigger,omitempty"`
	Delay   *Delay   `json:"delay,omitempty"`
	Goal    *Goal    `json:"goal,omitempty"`
}

type Sequence struct {
	Name  string         `json:"name"`
	Steps []SequenceStep `json:"steps"`
}

type SequenceExecutor struct{}
