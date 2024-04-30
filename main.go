package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var AnimeList map[string]Anime
var Planning PlanningType

func main() {
	getPlanning()

	// fmt.Println(Planning)

	an := getAnimeData("blue-exorcist")

	fmt.Println(an)

	// getAnData(ans["Lundi"][0])
	// for _, animes := range ans {
	// 	for _, anime := range animes {
	// 		getAnData(anime)
	// 	}
	// }

	// Lecture JSON
	// dat, _ := os.ReadFile("./models/animeModel.json")
	// var anime Anime
	// err := json.Unmarshal(dat, &anime)
	// fmt.Println(err)
	// fmt.Println(anime)

	// Ecriture JSON
	// var anime Anime
	// anime.Name = "Name"
	// anime.Desc = "Desc"
	// anime.Genre = []string{"Action", "Aventure"}
	// anime.IsTrack = true
	// anime.Saisons = make(map[string]Saison)
	// anime.Saisons["films"] = Saison{10, false}
	// dat, _ := json.Marshal(anime)
	// err := os.WriteFile("./test.json", dat, 0644)
	// fmt.Println(err)
}

func getAnimeData(animeLink string) Anime {
	c := colly.NewCollector()
	var anime Anime
	anime.Saisons = make(map[string]Saison)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)

		if len(strings.Split(r.URL.Path, "/")) == 3 {

			anime.Name = "Blue Exorcist"
			anime.IsTrack = false

			c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
				anime.Desc = e.Attr("content")
			})

			c.OnHTML("h2 + a", func(e *colly.HTMLElement) {
				anime.Genre = strings.Split(strings.ReplaceAll(e.Text, " ", ""), ",")
			})

			c.OnHTML("body", func(e *colly.HTMLElement) {
				var re = regexp.MustCompile(`(?m)panneauAnime\((.*?)\);`)

				tuiles := re.FindAllStringSubmatch(e.Text, -1)

				for key, tuile := range tuiles {
					if key != 0 {
						for _, tst := range tuile {
							if !strings.HasPrefix(tst, "panneauAnime") {

								str := strings.ReplaceAll(tst, " ", "")
								str = strings.ReplaceAll(str, "\"", "")
								str = strings.Split(str, ",")[1]
								str = strings.Split(str, "/")[0]

								anime.Saisons[str] = Saison{0, false}
							}

						}
					}
				}
			})
		}

	})

	c.Visit("https://anime-sama.fr/catalogue/" + animeLink)

	return anime
}

func getPlanning() {
	c := colly.NewCollector()

	c.OnHTML("#planningClass", func(e *colly.HTMLElement) {
		textSemaine := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(e.Text, "\t", ""), "\n", ""), " ", "")

		semaine := strings.Split(textSemaine, "Dimanche")
		textDimanche := semaine[1]
		semaine = strings.Split(semaine[0], "Samedi")
		textSamedi := semaine[1]
		semaine = strings.Split(semaine[0], "Vendredi")
		textVendredi := semaine[1]
		semaine = strings.Split(semaine[0], "Jeudi")
		textJeudi := semaine[1]
		semaine = strings.Split(semaine[0], "Mercredi")
		textMercredi := semaine[1]
		semaine = strings.Split(semaine[0], "Mardi")
		textMardi := semaine[1]
		semaine = strings.Split(semaine[0], "Lundi")
		textLundi := semaine[1]

		animesList := make(map[string]string)
		animesList["Lundi"] = textLundi
		animesList["Mardi"] = textMardi
		animesList["Mercredi"] = textMercredi
		animesList["Jeudi"] = textJeudi
		animesList["Vendredi"] = textVendredi
		animesList["Samedi"] = textSamedi
		animesList["Dimanche"] = textDimanche

		var re = regexp.MustCompile(`(?m)cartePlanningAnime\((.*?)\);`)

		for key, value := range animesList {
			cartes := re.FindAllStringSubmatch(value, -1)
			var tab []string

			for _, carte := range cartes {
				if strings.HasSuffix(carte[1], "\"VF\"") {

					an := strings.Split(strings.ReplaceAll(carte[1], "\"", ""), ",")

					tab = append(tab, strings.Split(an[1], "/")[0])
				}
			}

			switch key {
			case "Lundi":
				Planning.Lundi = tab
			case "Mardi":
				Planning.Mardi = tab
			case "Mercredi":
				Planning.Mercredi = tab
			case "Jeudi":
				Planning.Jeudi = tab
			case "Vendredi":
				Planning.Vendredi = tab
			case "Samedi":
				Planning.Samedi = tab
			case "Dimanche":
				Planning.Dimanche = tab
			}

		}
	})

	c.OnHTML("#planningClass + h2 + div", func(e *colly.HTMLElement) {
		text := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(e.Text, "\t", ""), "\n", ""), " ", "")

		var re = regexp.MustCompile(`(?m)cartePlanningAnime\((.*?)\);`)

		cartes := re.FindAllStringSubmatch(text, -1)

		var animes []string

		for _, carte := range cartes {
			if strings.HasSuffix(carte[1], "\"VF\"") {

				an := strings.Split(strings.ReplaceAll(carte[1], "\"", ""), ",")

				animes = append(animes, strings.Split(an[1], "/")[0])
			}
		}

		Planning.Inconue = animes
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://anime-sama.fr/planning/")
}
