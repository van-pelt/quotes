package server

import (
	"encoding/json"
	"fmt"
	"github.com/van-pelt/quotes/pkg/quotes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const timeout = 15 * time.Second

var QT quotes.Quotes

type ErrMess struct {
	ErrCode int    `json:"ErrorCode"`
	ErrMess string `json:"ErrorMessage"`
}

func Run() {
	QT = *quotes.NewQuotes()
	QT.AddQuotes("cat1", "text 1", "kizildur")
	QT.AddQuotes("cat1", "text 2", "dazdranagor")
	QT.AddQuotes("cat2", "text 3", "boris_britva")
	QT.AddQuotes("cat3", "text 4", "kizildur")
	QT.AddQuotes("cat3", "Если тебя где-то не ждут в рванных носках,то и в целых туда идтить не нужно", "Д.Стетхем")
	handler := http.NewServeMux()
	handler.HandleFunc("/quotes", HandleGetAllQuotes)                     //GET
	handler.HandleFunc("/quotes/category/", HandleGetAllQuotesByCategory) //GET
	handler.HandleFunc("/quotes/add", HandleAddQuotes)                    //POST
	handler.HandleFunc("/quotes/update/", HandleUpdateQuotes)             //PUT
	handler.HandleFunc("/quotes/delete/", HandleDeleteQuotes)             //DELETE
	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		IdleTimeout:    timeout,
		MaxHeaderBytes: 128,
	}

	log.Printf("Listening on http://%s\n", s.Addr)

	/*go func() {
		log.Printf("Start scedule")
		for {
			time.Sleep(10 * time.Second)
			QT.DeleteQuoteByTime(10)

		}
	}()*/

	err := s.ListenAndServe()
	if err != nil {
		log.Print(fmt.Errorf("ListenAndServe():%w", err))
	}

}

func HandleGetAllQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		dd := QT.GetAllQuotes()
		q, _ := json.Marshal(dd)
		w.Write(q)
	} else {
		HandleMethodIsNotAllowed(w, r)
	}
}

func HandleGetAllQuotesByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		var re = regexp.MustCompile(`(?s)/quotes/category/\w+`)
		if len(re.FindStringIndex(r.URL.Path)) > 0 {
			s := strings.Split(r.URL.Path, "/")
			if len(s) == 4 {
				dd := QT.FindQuotesByCategory(s[3])
				q, _ := json.Marshal(dd)
				w.WriteHeader(http.StatusOK)
				w.Write(q)
			}
		}
	} else {
		HandleMethodIsNotAllowed(w, r)
	}
}

func HandleAddQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		quote := quotes.Quote{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ErrHandler(w, err, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &quote)
		if err != nil {
			ErrHandler(w, err, http.StatusInternalServerError)
		}
		id := QT.AddQuotes(quote.Category, quote.Quote, quote.Author)
		w.WriteHeader(http.StatusOK)
		addJson, _ := json.Marshal(id)
		w.Write(addJson)
	} else {
		HandleMethodIsNotAllowed(w, r)
	}
}

func HandleUpdateQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPut {
		var re = regexp.MustCompile(`(?s)/quotes/update/\d+`)
		if len(re.FindStringIndex(r.URL.Path)) > 0 {
			s := strings.Split(r.URL.Path, "/")
			if len(s) == 4 {
				id, err := strconv.ParseInt(s[3], 10, 64)
				if err != nil {
					ErrHandler(w, err, http.StatusInternalServerError)
					return
				}
				quote := quotes.Quote{}
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					ErrHandler(w, err, http.StatusInternalServerError)
					return
				}
				err = json.Unmarshal(body, &quote)
				if err != nil {
					ErrHandler(w, err, http.StatusInternalServerError)
					return
				}
				err = QT.UpdateQuotes(id, quote.Category, quote.Quote, quote.Author)
				if err != nil {
					ErrHandler(w, err, http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				addJson, _ := json.Marshal("OK")
				w.Write(addJson)
			}
		}
	} else {
		HandleMethodIsNotAllowed(w, r)
	}
}

func HandleDeleteQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodDelete {
		var re = regexp.MustCompile(`(?s)/quotes/delete/\d+`)
		if len(re.FindStringIndex(r.URL.Path)) > 0 {
			s := strings.Split(r.URL.Path, "/")
			if len(s) == 4 {
				id, err := strconv.ParseInt(s[3], 10, 64)
				if err != nil {
					ErrHandler(w, err, http.StatusInternalServerError)
					return
				}
				err = QT.DeleteQuote(id)
				if err != nil {
					ErrHandler(w, err, http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				addJson, _ := json.Marshal("OK")
				w.Write(addJson)
			}
		}
	} else {
		HandleMethodIsNotAllowed(w, r)
	}
}

func ErrHandler(w http.ResponseWriter, err error, statusCode int) {
	er := ErrMess{
		ErrCode: statusCode,
		ErrMess: "Error:" + err.Error(),
	}
	q, _ := json.Marshal(er)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(q))
	return
}

func HandleMethodIsNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	msg, _ := json.Marshal(fmt.Sprintf("Method %s not allowed", r.Method))
	w.Write(msg)
}
