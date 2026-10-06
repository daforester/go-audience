package requests

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/models"
)

type NewSubmissionRequest struct {
	Kind models.SubmissionKind
	Text string
}
