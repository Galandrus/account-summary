package libs

import "account-summary/src/models"

type EmailSenderInterface interface {
	SendAccountSummaryEmail(to string, subject string, summary models.AccountSummary) error
}
