package quotes

import (
	"errors"
	"log"
	"time"
)

var ErrQuotesNotFound = errors.New("Quotes not found")

type Quote struct {
	ID         int64  `json:"ID"`
	Author     string `json:"Author"`
	Quote      string `json:"Quote"`
	Category   string `json:"Category"`
	DateCreate int64  `json:"DateCreate"`
}

type Quotes struct {
	quotes []Quote
	lastId int64
}

func NewQuotes() *Quotes {
	return &Quotes{}
}

func (q *Quotes) FindQuotesByID(id int64) (*Quote, error) {
	for i, quote := range q.quotes {
		if quote.ID == id {
			return &q.quotes[i], nil
		}
	}
	return nil, ErrQuotesNotFound
}

func (q *Quotes) FindQuotesByCategory(category string) []*Quote {
	var quoteCat []*Quote
	for _, quote := range q.quotes {
		if quote.Category == category {
			quoteCat = append(quoteCat, &quote)
		}
	}
	return quoteCat
}

func (q *Quotes) GetAllQuotes() []Quote {
	return q.quotes
}

func (q *Quotes) AddQuotes(category, quoteText, author string) int64 {
	q.lastId++
	q.quotes = append(q.quotes, Quote{
		ID:         q.lastId,
		Author:     author,
		Quote:      quoteText,
		Category:   category,
		DateCreate: time.Now().Unix(),
	})
	return q.lastId
}

func (q *Quotes) UpdateQuotes(quoteID int64, category, quoteText, author string) error {
	qt, err := q.FindQuotesByID(quoteID)
	if err != nil {
		return err
	}
	qt.Category = category
	qt.Quote = quoteText
	qt.Author = author
	qt.DateCreate = time.Now().Unix()

	return nil
}

func (q *Quotes) DeleteQuote(quoteID int64) error {
	for i, qt := range q.quotes {
		if qt.ID == quoteID {
			q.quotes = append(q.quotes[:i], q.quotes[i+1:]...)
			return nil
		}
	}
	return ErrQuotesNotFound
}

func (q *Quotes) DeleteQuoteByTime(timestampOffset int64) {

	currentTimestamp := time.Now().Unix()
	for i := 0; i < len(q.quotes); i++ {
		if q.quotes[i].DateCreate+timestampOffset < currentTimestamp {
			log.Print("Delete ID=", q.quotes[i].ID)
			q.quotes = append(q.quotes[:i], q.quotes[i+1:]...)

			i = 0
			continue
		}
	}
}
