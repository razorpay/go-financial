package calculator

import (
	"errors"
	"time"

	"github.com/razorpay/go-financial/enums/paymentperiod"

	"github.com/razorpay/go-financial/enums/interesttype"

	"github.com/razorpay/go-financial/enums/frequency"
)

type Config struct {
	StartDate      time.Time
	EndDate        time.Time
	Frequency      frequency.Type
	AmountBorrowed int64
	InterestType   interesttype.Type
	Interest       int64
	PaymentPeriod  paymentperiod.Type
	Round          bool
	periods        int64       // derived
	startDates     []time.Time // derived
	endDates       []time.Time // derived
}

func (c *Config) SetPeriodsAndDates() error {
	sy, sm, sd := c.StartDate.Date()
	startDate := time.Date(sy, sm, sd, 0, 0, 0, 0, c.StartDate.Location())

	ey, em, ed := c.EndDate.Date()
	endDate := time.Date(ey, em, ed, 0, 0, 0, 0, c.EndDate.Location())

	days := int64(endDate.Sub(startDate).Hours()/24) + 1
	switch c.Frequency {
	case frequency.DAILY:
		c.periods = days
		for i := 0; i < int(c.periods); i++ {
			date := startDate.AddDate(0, 0, i)
			if i == 0 {
				c.startDates = append(c.startDates, c.StartDate)
				c.endDates = append(c.endDates, getEndDates(date, frequency.DAILY))
			} else {
				c.startDates = append(c.startDates, date)
				c.endDates = append(c.endDates, getEndDates(date, frequency.DAILY))
			}
		}
	case frequency.WEEKLY:
		if days%7 != 0 {
			return errors.New("uneven end date")
		}
		c.periods = days / 7
		for i := 0; i < int(c.periods); i++ {
			date := startDate.AddDate(0, 0, 7*i)
			if i == 0 {
				c.startDates = append(c.startDates, c.StartDate)
				c.endDates = append(c.endDates, getEndDates(date, frequency.WEEKLY))
			} else {
				c.startDates = append(c.startDates, date)
				c.endDates = append(c.endDates, getEndDates(date, frequency.WEEKLY))
			}

		}

	case frequency.MONTHLY:
		months, err := getMonthsBetweenDates(c.StartDate, c.EndDate)
		if err != nil {
			return err
		}
		c.periods = int64(*months)
		for i := 0; i < int(c.periods); i++ {
			date := startDate.AddDate(0, i, 0)
			if i == 0 {
				c.startDates = append(c.startDates, c.StartDate)
				c.endDates = append(c.endDates, getEndDates(date, frequency.MONTHLY))
			} else {
				c.startDates = append(c.startDates, date)
				c.endDates = append(c.endDates, getEndDates(date, frequency.MONTHLY))
			}

		}

	case frequency.ANNUALLY:
		years, err := getYearsBetweenDates(startDate, endDate)
		if err != nil {
			return err
		}
		c.periods = int64(*years)
		for i := 0; i < int(c.periods); i++ {
			date := startDate.AddDate(i, 0, 0)
			if i == 0 {
				c.startDates = append(c.startDates, c.StartDate)
				c.endDates = append(c.endDates, getEndDates(date, frequency.ANNUALLY))
			} else {
				c.startDates = append(c.startDates, date)
				c.endDates = append(c.endDates, getEndDates(date, frequency.ANNUALLY))
			}

		}
	default:
		return ErrInvalidFrequency

	}
	return nil
}

func getMonthsBetweenDates(start time.Time, end time.Time) (*int, error) {
	count := 0
	for start.Before(end) {
		start = start.AddDate(0, 1, 0)
		count++
	}
	finalDate := start.AddDate(0, 0, -1)
	if !finalDate.Equal(end) {
		return nil, ErrUnevenEndDate
	}
	return &count, nil
}

func getYearsBetweenDates(start time.Time, end time.Time) (*int, error) {
	count := 0
	for start.Before(end) {
		start = start.AddDate(1, 0, 0)
		count++
	}
	finalDate := start.AddDate(0, 0, -1)
	if !finalDate.Equal(end) {
		return nil, ErrUnevenEndDate
	}
	return &count, nil
}

func getEndDates(date time.Time, freq frequency.Type) time.Time {
	var nextDate time.Time
	switch freq {
	case frequency.DAILY:
		nextDate = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	case frequency.WEEKLY:
		date = date.AddDate(0, 0, 6)
		nextDate = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	case frequency.MONTHLY:
		date = date.AddDate(0, 1, 0).AddDate(0, 0, -1)
		nextDate = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	case frequency.ANNUALLY:
		date = date.AddDate(1, 0, 0).AddDate(0, 0, -1)
		nextDate = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	}
	return nextDate
}

func (c *Config) getInterestRatePerPeriodInDecimal() float64 {
	return float64(c.Interest) / 100 / 100 / float64(c.Frequency.Value())
}
