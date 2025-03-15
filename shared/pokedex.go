package shared

var Pokedex map[string]Pokemon

func init() {
	Pokedex = make(map[string]Pokemon)
}

func AddToPokedex(p *Pokemon) {
	Pokedex[p.Name] = *p
}
