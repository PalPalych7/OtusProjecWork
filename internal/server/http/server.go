package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

type Server struct {
	myCtx     context.Context
	myStorage Storage
	myLogger  Logger
	HTTPConf  string
	myHTTP    http.Server
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Storage interface {
	Connect() error
	AddBannerSlot(slotID int, bannerID int) error
	DelBannerSlot(slotID int, bannerID int) error
	BannerClick(slotID int, bannerID int, socGroupID int) error
	GetBannerForSlot(slotID int, socGroupID int) (int, error)
	GetBannerStat() ([]sqlstorage.BannerStatStruct, error)
	ChangeSendStatID(ID int) error
	Close() error
}

func NewServer(ctx context.Context, app Storage, httpConf string, myLogger Logger) *Server {
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
	mux.HandleFunc("/AddBannerSlot", s.AddBannerSlot)
	mux.HandleFunc("/GetBannerForSlot", s.GetBannerForSlot)
	mux.HandleFunc("/BannerClick", s.BannerClick)
	mux.HandleFunc("/DelBannerSlot", s.DelBannerSlot)

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

func (s *Server) AddBannerSlot(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("AddBannerSlot")
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	mySB := SlotBanner{}
	if err := json.Unmarshal(myRaw, &mySB); err != nil {
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("I am heare", mySB)
	myErr := s.myStorage.AddBannerSlot(mySB.SlotID, mySB.BannerID)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) BannerClick(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("BannerClick")
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myFBC := ForBannerClick{}
	if err := json.Unmarshal(myRaw, &myFBC); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.BannerClick(myFBC.SlotID, myFBC.BannerID, myFBC.SocGroupID)
	s.myLogger.Info("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DelBannerSlot(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("DelBannerSlot")
	myRaw1 := getBodyRaw(req.Body)
	if myRaw1 == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	mySB := SlotBanner{}
	if err1 := json.Unmarshal(myRaw1, &mySB); err1 != nil {
		s.myLogger.Info(myRaw1)
		s.myLogger.Error("Error json.Unmarshal - " + err1.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.DelBannerSlot(mySB.SlotID, mySB.BannerID)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetBannerForSlot(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("GetBannerForSlot")
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myFGBS := ForGetBanner{}
	if err := json.Unmarshal(myRaw, &myFGBS); err != nil {
		s.myLogger.Info(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.myLogger.Info("myFGBS=", myFGBS)
	bannerID, myEr := s.myStorage.GetBannerForSlot(myFGBS.SlotID, myFGBS.SocGroupID)
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
