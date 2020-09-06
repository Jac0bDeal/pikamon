package models

import "fmt"

type Pokemon struct {
	ID        int
	PokemonID int
	TrainerID int
	Name      string
}

func (p *Pokemon) ListingInfo() string {
	return fmt.Sprintf("**%s** | ID: %d", p.Name, p.ID)
}
