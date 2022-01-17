package seating

import (
	"fmt"
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
	Groups  []*models.Group
}

type block struct {
	hasAisle bool
	rank     string
	seats    []*models.Seat
	row      *models.Row
}

type allocatedSeatsType map[string]struct{}

func createAvailableBlocks(rows map[string]*row, allocatedSeats allocatedSeatsType) []block {
	var availableBlocks []block
	for _, row := range rows {
		availableBlockInRow := block{
			row: row.Row,
		}

		for j, seat := range row.Seats {
			_, isSeatAllocated := allocatedSeats[seat.ID]

			if seat.Broken || isSeatAllocated {
				if len(availableBlockInRow.seats) > 0 {
					availableBlocks = append(availableBlocks, availableBlockInRow)
					availableBlockInRow.seats = []*models.Seat{}
				}

				continue
			}

			if len(availableBlockInRow.seats) > 0 && seat.Rank != availableBlockInRow.seats[len(availableBlockInRow.seats)-1].Rank {
				availableBlocks = append(availableBlocks, availableBlockInRow)
				availableBlockInRow.seats = []*models.Seat{}
				continue
			}

			if len(availableBlockInRow.seats) == 0 {
				availableBlockInRow = block{
					row:  row.Row,
					rank: seat.Rank,
				}
			}

			if seat.Aisle {
				availableBlockInRow.hasAisle = true
			}

			availableBlockInRow.seats = append(availableBlockInRow.seats, seat)

			if j == len(row.Seats)-1 {
				availableBlocks = append(availableBlocks, availableBlockInRow)
			}
		}
	}

	return availableBlocks
}

func seatGroup(group *models.Group, section section, allocatedSeats allocatedSeatsType) (*models.Ticket, error) {
	availableBlocks := createAvailableBlocks(section.Rows, allocatedSeats)
	for _, block := range availableBlocks {
		if int(group.Count) > len(block.seats) || group.Rank != block.rank || (group.Aisle && !block.hasAisle) {
			continue
		}

		availableSeats := block.seats[0:group.Count]
		if group.Aisle && !availableSeats[0].Aisle {
			availableSeats = block.seats[len(block.seats)-int(group.Count)+1:]
		}

		var seats []string
		for _, seat := range availableSeats {
			seats = append(seats, seat.ID)
		}

		ticket := models.Ticket{
			GroupId: group.ID,
			Seats:   seats,
		}

		return &ticket, nil
	}

	return &models.Ticket{}, fmt.Errorf("no available block found")
}

func seatGroups(section section, m models.Models) []*models.Ticket {
	allocatedSeats := make(allocatedSeatsType)
	var tickets []*models.Ticket

	for _, group := range section.Groups {
		t, err := m.TicketGetByGroupId(group.ID)
		if err != nil {
			panic(err)
		}

		if t.ID != "" {
			m.TicketDelete(t)
		}

		ticket, err := seatGroup(group, section, allocatedSeats)
		if err != nil {
			log.Println("group", group)
			log.Println("section", section)
			log.Println("Unable to seat group", group.ID, "in section", section.Section.Name)
		}

		for _, seat := range ticket.Seats {
			allocatedSeats[seat] = struct{}{}
		}
		tickets = append(tickets, ticket)
	}

	return tickets
}

func Process(m models.Models) {
	sections, err := GetSections(m)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, s := range sections {
		tickets := seatGroups(*s, m)
		for _, v := range tickets {
			m.TicketSave(v)
		}
	}
}

func GetSections(m models.Models) ([]*section, error) {
	dbSections, err := m.SectionGetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load sections: %v", err)
	}

	var sections []*section
	for _, dbSection := range dbSections {
		s := section{
			Section: dbSection,
		}

		rows, err := m.RowGetBySection(dbSection.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to load rows: %v", err)
		}

		s.Rows = make(map[string]*row)
		for _, dbRow := range rows {
			s.Rows[dbRow.Name] = &row{
				Row: dbRow,
			}
		}

		seats, err := m.SeatGetBySection(dbSection.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to load seats: %v", err)
		}

		for _, dbSeat := range seats {
			if s.Rows[dbSeat.Row] == nil {
				continue
			}
			s.Rows[dbSeat.Row].Seats = append(s.Rows[dbSeat.Row].Seats, dbSeat)
		}

		groups, err := m.GroupGetBySection(dbSection.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to load groups: %v", err)
		}

		s.Groups = append(s.Groups, groups...)
		sections = append(sections, &s)
	}

	return sections, nil
}
