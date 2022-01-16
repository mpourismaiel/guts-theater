package seating

import (
	"encoding/json"
	"log"

	"mpourismaiel.dev/guts/store/models"
)

type row struct {
	Row   *models.Row    `json:"row"`
	Seats []*models.Seat `json:"seats"`
}

type section struct {
	Section *models.Section `json:"section"`
	Rows    map[string]*row `json:"rows"`
	Groups  []*models.Group `json:"groups"`
}

func Process(m models.Models) {
	log.Println("Starting process")
	dbSections, err := m.SectionGetAll()
	if err != nil {
		log.Fatalf("failed to load sections: %v", err)
		return
	}

	log.Println("Fetched sections")
	sections := make(map[string]*section)
	for _, dbSection := range dbSections {
		sections[dbSection.Name] = &section{
			Section: dbSection,
		}

		rows, err := m.RowGetBySection(dbSection.Name)
		if err != nil {
			log.Fatalf("failed to load rows: %v", err)
			return
		}

		log.Println("Fetched rows for section:", dbSection.Name)
		sections[dbSection.Name].Rows = make(map[string]*row)
		for _, dbRow := range rows {
			sections[dbSection.Name].Rows[dbRow.Name] = &row{
				Row: dbRow,
			}
		}

		seats, err := m.SeatGetBySection(dbSection.Name)
		if err != nil {
			log.Fatalf("failed to load seats: %v", err)
			return
		}

		log.Println("Fetched seats for section:", dbSection.Name)
		for _, dbSeat := range seats {
			sections[dbSection.Name].Rows[dbSeat.Row].Seats = append(sections[dbSection.Name].Rows[dbSeat.Row].Seats, dbSeat)
		}

		groups, err := m.GroupGetBySection(dbSection.Name)
		if err != nil {
			log.Fatalf("failed to load groups: %v", err)
			return
		}

		log.Println("Fetched groups for section:", dbSection.Name)
		sections[dbSection.Name].Groups = append(sections[dbSection.Name].Groups, groups...)
	}

	res, err := json.Marshal(sections)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(string(res))
}
