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

func NewServer(ctx context.Context, app sqlstorage.MyStorage, httpConf string, myLogger logger.MyLogger) MyServer {
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
		s.myLogger.Error("Ошибка при обработке тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct1 := SlotBanner{}
	if err1 := json.Unmarshal(myRaw1, &myStruct1); err1 != nil {
		s.myLogger.Info(myRaw1)
		s.myLogger.Error("Ошибка при переводе json в структуру - " + err1.Error())
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
