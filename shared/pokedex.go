package shared

var pokedex map[string]Pokemon

func init() {
	pokedex = make(map[string]Pokemon)
}

func AddToPokedex(p *Pokemon) {
	pokedex[p.Name] = *p
}
