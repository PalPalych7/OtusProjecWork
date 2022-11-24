package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

type Server struct {
	myCtx     context.Context
	myStorage sqlstorage.MyStorage
	myLogger  myLogger
	HTTPConf  string
	myHTTP    http.Server
}

type MyServer interface {
	Serve() error
	Stop() error
}

type myLogger interface {
	//	Trace(args ...interface{})
	//	Debug(args ...interface{})
	Info(args ...interface{})
	//	Print(args ...interface{})
	//	Warning(args ...interface{})
	Error(args ...interface{})
	//	Fatal(args ...interface{})
}

func NewServer(ctx context.Context, app sqlstorage.MyStorage, httpConf string, myLogger myLogger) MyServer {
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

func (s *Server) Serve() error {
	s.myHTTP.Addr = s.HTTPConf
	mux := http.NewServeMux()
	mux.HandleFunc("/AddBannerSlot", s.AddBannerSlotFunc)
	mux.HandleFunc("/GetBannerForSlot", s.GetBannerForSlotFunc)
	mux.HandleFunc("/BannerClick", s.BannerClickFunc)
	mux.HandleFunc("/DelBannerSlot", s.DelBannerSlotFunc)
	http.ListenAndServe(s.myHTTP.Addr, s.loggingMiddleware(mux))
	return nil
}

func (s *Server) Stop() error {
	err := s.myHTTP.Shutdown(s.myCtx)
	return err
}

func (s *Server) AddBannerSlotFunc(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("AddBannerSlot")
	myRaw := getBodyRow(req.Body)
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
	myRaw := getBodyRow(req.Body)
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
	myRaw1 := getBodyRow(req.Body)
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
	myRaw := getBodyRow(req.Body)
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
