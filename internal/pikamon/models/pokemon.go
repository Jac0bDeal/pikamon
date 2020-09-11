package models

import "fmt"

type Pokemon struct {
	ID        string
	PokemonID int
	TrainerID string
	Name      string
}

func (p *Pokemon) ListingInfo() string {
	return fmt.Sprintf("**%s** | ID: %s", p.Name, p.ID)
}
