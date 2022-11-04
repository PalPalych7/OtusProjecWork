package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

type Server struct {
	myCtx     context.Context
	myStorage sqlstorage.MyStorage
	myLogger  logger.MyLogger
	HTTPConf  string
	myHTTP    http.Server
}

type MyServer interface {
	Start() error
	Stop() error
}

func NewServer(ctx context.Context, app sqlstorage.MyStorage, httpConf string, myLogger logger.MyLogger) MyServer /**Server*/ {
	return &Server{myCtx: ctx, myStorage: app, myLogger: myLogger, HTTPConf: httpConf}
}

func getBodyRow(reqBody io.ReadCloser) []byte {
	raw, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil
	}
	defer reqBody.Close()
	return raw
}

func (s *Server) Start() error {
	s.myHTTP.Addr = s.HTTPConf
	fmt.Println("serv=", s.HTTPConf)
	s.myLogger.Info("serv=", s.HTTPConf)

	mux := http.NewServeMux()

	mux.HandleFunc("/AddBannerSlot", s.AddBannerSlotFunc)
	mux.HandleFunc("/GetBannerForSlot", s.GetBannerForSlotFunc)
	mux.HandleFunc("/BannerClick", s.BannerClickFunc)
	mux.HandleFunc("/DelBannerSlot", s.DelBannerSlotFunc)
	http.ListenAndServe(s.myHTTP.Addr, s.loggingMiddleware(mux)) //nolint
	return nil
}

func (s *Server) Stop() error {
	fmt.Println("start finish server")
	err := s.myHTTP.Shutdown(s.myCtx)
	fmt.Println("end finish server")
	return err
}

func (s *Server) AddBannerSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("AddBannerSlot")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := SlotBanner{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Ошибка перевода json в структуру - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.myLogger.Info("myStruct=", myStruct, "slotId=", myStruct.SlotId)
	myErr := s.myStorage.AddBannerSlot(myStruct.SlotId, myStruct.BannerId)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) BannerClickFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("BannerClick")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForBannerClick{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Ошибка перевода json в структуру - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.BannerClick(myStruct.SlotId, myStruct.BannerId, myStruct.SocGroupId)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DelBannerSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("DelBannerSlot")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := SlotBanner{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Ошибка перевода json в структуру - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.myLogger.Info("myStruct=", myStruct, "slotId=", myStruct.SlotId)
	myErr := s.myStorage.DelBannerSlot(myStruct.SlotId, myStruct.BannerId)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetBannerForSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("GetBannerForSlot")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForGetBanner{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Ошибка перевода json в структуру - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.myLogger.Info("myStruct=", myStruct)
	myErr, bannerId := s.myStorage.GetBannerForSlot(myStruct.SlotId, myStruct.SocGroupId)
	s.myLogger.Info(myErr, bannerId)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rawResp, err3 := json.Marshal(bannerId)
		if err3 == nil {
			rw.Write(rawResp)
		}
	}
}

/*
func (s *Server) CreateEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("CreateEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForCreate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Info(myRaw)
		s.App.Error("Ошибка перевода json в структуру - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.App.CreateEvent(s.myCtx, myStruct.Title, myStruct.StartDate, myStruct.Details, int(myStruct.UserID))
	if myErr != nil {
		s.App.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) UpdateEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("UpdateEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForUpdate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("myStruct=", myStruct)
	myErr := s.App.UpdateEvent(s.myCtx, myStruct.EventID, myStruct.Title, myStruct.StartDate, myStruct.Details, int(myStruct.UserID)) //nolint
	if myErr != nil {
		s.App.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DeleteEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("DeleteEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForDelete{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.App.DeleteEvent(s.myCtx, myStruct.EventID)
	if myErr != nil {
		s.App.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetEventByDateFunc(rw http.ResponseWriter, req *http.Request) { //nolint:dupl
	s.App.Info("GetEventByDate")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := StartDate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventList, myErr := s.App.GetEventByDate(s.myCtx, myStruct.StartDateStr)
	if myErr == nil {
		rawResp, err3 := json.Marshal(eventList)
		if err3 == nil {
			rw.Write(rawResp)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			s.App.Error(err3)
			return
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		s.App.Error(myErr)
	}
}

func (s *Server) GetEventMonthFunc(rw http.ResponseWriter, req *http.Request) { //nolint:dupl
	s.App.Info("GetEventMonth")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса!")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := StartDate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру!")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventList, myErr := s.App.GetEventMonth(s.myCtx, myStruct.StartDateStr)
	if myErr == nil {
		rawResp, err3 := json.Marshal(eventList)
		if err3 == nil {
			rw.Write(rawResp)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			s.App.Error(err3)
			return
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		s.App.Error(myErr)
	}
}

func (s *Server) GetEventByWeekFunc(rw http.ResponseWriter, req *http.Request) { //nolint:dupl
	s.App.Info("GetEventByWeekFunc")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := StartDate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventList, myErr := s.App.GetEventWeek(s.myCtx, myStruct.StartDateStr)
	if myErr == nil {
		rawResp, err3 := json.Marshal(eventList)
		if err3 == nil {
			rw.Write(rawResp)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			s.App.Error(err3)
			return
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		s.App.Error(myErr)
	}
}
*/
