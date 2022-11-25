package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

type Server struct {
	myCtx     context.Context
	myStorage myStorage
	myLogger  myLogger
	HTTPConf  string
	myHTTP    http.Server
}

type myLogger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type myStorage interface {
	Connect() error
	AddBannerSlot(slotID int, bannerID int) error
	DelBannerSlot(slotID int, bannerID int) error
	BannerClick(slotID int, bannerID int, socGroupID int) error
	GetBannerForSlot(slotID int, socGroupID int) (int, error)
	GetBannerStat() ([]sqlstorage.BannerStatStruct, error)
	ChangeSendStatID(ID int) error
	Close() error
}

func NewServer(ctx context.Context, app myStorage, httpConf string, myLogger myLogger) *Server {
	return &Server{myCtx: ctx, myStorage: app, myLogger: myLogger, HTTPConf: httpConf}
}

func getBodyRaw(reqBody io.ReadCloser) []byte {
	raw, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil
	}
	defer reqBody.Close()
	return raw
}

func (s *Server) Serve() error {
	s.myHTTP.Addr = s.HTTPConf
	mux := http.NewServeMux()
	mux.HandleFunc("/AddBannerSlot", s.AddBannerSlotFunc)
	mux.HandleFunc("/GetBannerForSlot", s.GetBannerForSlotFunc)
	mux.HandleFunc("/BannerClick", s.BannerClickFunc)
	mux.HandleFunc("/DelBannerSlot", s.DelBannerSlotFunc)

	server := &http.Server{
		Addr:              s.myHTTP.Addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           s.loggingMiddleware(mux),
	}

	err := server.ListenAndServe()
	if err != nil {
		s.myLogger.Error(err)
	}
	return err
}

func (s *Server) Stop() error {
	err := s.myHTTP.Shutdown(s.myCtx)
	return err
}

func (s *Server) AddBannerSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("AddBannerSlot")
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := SlotBanner{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.myLogger.Info("myStruct=", myStruct, "slotId=", myStruct.SlotID)
	myErr := s.myStorage.AddBannerSlot(myStruct.SlotID, myStruct.BannerID)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) BannerClickFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("BannerClick")
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForBannerClick{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.BannerClick(myStruct.SlotID, myStruct.BannerID, myStruct.SocGroupID)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DelBannerSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("DelBannerSlot")
	myRaw1 := getBodyRaw(req.Body)
	if myRaw1 == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct1 := SlotBanner{}
	if err1 := json.Unmarshal(myRaw1, &myStruct1); err1 != nil {
		s.myLogger.Info(myRaw1)
		s.myLogger.Error("Error json.Unmarshal - " + err1.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.DelBannerSlot(myStruct1.SlotID, myStruct1.BannerID)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetBannerForSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("GetBannerForSlot")
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForGetBanner{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.myLogger.Info("myStruct=", myStruct)
	bannerID, myEr := s.myStorage.GetBannerForSlot(myStruct.SlotID, myStruct.SocGroupID)
	s.myLogger.Info(bannerID, myEr)
	if myEr != nil {
		s.myLogger.Error(myEr)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rawResp, err3 := json.Marshal(bannerID)
		if err3 == nil {
			rw.Write(rawResp)
		}
	}
}
