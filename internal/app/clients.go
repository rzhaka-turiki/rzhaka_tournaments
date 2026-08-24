package app

import (
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/apexverifier"
	"github.com/rzhaka-turiki/rzhaka_tournaments/internal/client/matchapi"
)

type Clients struct {
	ApexVerifier apexverifier.Client
	MatchAPI     matchapi.Client
}
