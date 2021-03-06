package seating

/**
* this package is responsible for assigning seats to customer groups.
* can be called concurrently or synchronously but it's advised to create a
* goroutine as this function has multiple loops and db calls.
**/

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/mpourismaiel/guts-theater/store/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type seating struct {
	logger *zap.Logger
}

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

func (s *seating) createAvailableBlocks(rows map[string]*row, allocatedSeats allocatedSeatsType) []block {
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

func (s *seating) seatGroup(group *models.Group, section section, allocatedSeats allocatedSeatsType) (*models.Ticket, error) {
	availableBlocks := s.createAvailableBlocks(section.Rows, allocatedSeats)
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

func (s *seating) seatGroups(section section, m models.Models) []*models.Ticket {
	allocatedSeats := make(allocatedSeatsType)
	var tickets []*models.Ticket

	for _, group := range section.Groups {
		t, err := m.TicketGetByGroupId(group.ID)
		if err != nil {
			sentry.CaptureException(err)
			s.logger.Fatal("could not load tickets")
		}

		if t.ID != "" {
			m.TicketDelete(t)
		}

		ticket, err := s.seatGroup(group, section, allocatedSeats)
		if err != nil {
			fields := []zapcore.Field{
				zap.String("group", group.ID),
				zap.String("section", section.Section.Name),
			}
			s.logger.Warn("Unable to seat group", fields...)
			continue
		}

		for _, seat := range ticket.Seats {
			allocatedSeats[seat] = struct{}{}
		}
		tickets = append(tickets, ticket)
	}

	return tickets
}

func Process(m models.Models, logger *zap.Logger) {
	s := &seating{
		logger: logger,
	}
	sections, err := GetSections(m, logger)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	for _, section := range sections {
		tickets := s.seatGroups(*section, m)
		for _, v := range tickets {
			m.TicketSave(v)
		}
	}
}

func GetSections(m models.Models, logger *zap.Logger) ([]*section, error) {
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
