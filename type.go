package main

type Saison struct {
	NbEps    int
	IsFinish bool
}

type Anime struct {
	Name    string
	Desc    string
	Genre   []string
	IsTrack bool
	Saisons map[string]Saison
}

type PlanningType struct {
	Lundi    []string
	Mardi    []string
	Mercredi []string
	Jeudi    []string
	Vendredi []string
	Samedi   []string
	Dimanche []string
	Inconue  []string
}
