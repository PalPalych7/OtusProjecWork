package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type mySuite struct {
	suite.Suite
	ctx       context.Context
	client    http.Client
	hostName  string
	DBConnect *sql.DB
}

type SlotBanner struct {
	SlotID   int
	BannerID int
}

type ForBannerClick struct {
	SlotID     int
	BannerID   int
	SocGroupID int
}

type ForGetBanner struct {
	SlotID     int
	SocGroupID int
}

var (
	err      error
	bodyRaw  []byte
	req      *http.Request
	resp     *http.Response
	countRec int
	bannerID int
)

func (s *mySuite) CheckCountRec(myQueryText string, expCount int) {
	mySQLRows, err := s.DBConnect.QueryContext(s.ctx, myQueryText) //nolint
	s.Require().NoError(err)
	defer mySQLRows.Close()
	mySQLRows.Next()
	err = mySQLRows.Scan(&countRec)
	fmt.Println("countRec=", countRec)
	s.Require().NoError(err)
	s.Require().Equal(expCount, countRec)
}

func (s *mySuite) SetupSuite() {
	fmt.Println("start setup suit")
	s.client = http.Client{
		Timeout: time.Second * 5,
	}

	s.hostName = "http://mainSevice:5000/"
	s.ctx = context.Background()
	myStr := "postgres://testuser:123456@postgres_db:5432/otusfinalproj?sslmode=disable" // через докер
	//	myStr := "postgres://testuser:123456@localhost:5432/otusfinalproj?sslmode=disable" // локально

	fmt.Println("start connect to postgrace:", myStr)
	s.DBConnect, err = sql.Open("postgres", myStr)
	fmt.Println("result: s.DBConnect:", err)
	if err == nil {
		err = s.DBConnect.PingContext(s.ctx)
	}
	s.Require().NoError(err)

	s.CheckCountRec("select count(*) RC from banner", 20)
	_, err = s.DBConnect.ExecContext(s.ctx, "delete from slot_banner")
	s.Require().NoError(err)
	s.CheckCountRec("select count(*) RC from slot_banner", 0)

	_, err = s.DBConnect.ExecContext(s.ctx, "delete from banner_stat")
	s.Require().NoError(err)
	s.CheckCountRec("select count(*) RC from banner_stat", 0)
	fmt.Println("finish setup suit")
}

func (s *mySuite) TearDownSuite() {
	fmt.Println("start TearDownSuite")
	_, err = s.DBConnect.ExecContext(s.ctx, "delete from slot_banner;delete from banner_stat;")
	s.Require().NoError(err)
	s.DBConnect.Close()
	fmt.Println("finish TearDownSuite")
}

func (s *mySuite) SendRequest(myMethodName string, myStruct interface{}) []byte {
	bodyRaw, err = json.Marshal(myStruct)

	s.Require().NoError(err)
	req, err = http.NewRequestWithContext(s.ctx, http.MethodPost, s.hostName+myMethodName, bytes.NewBuffer(bodyRaw))
	s.Require().NoError(err)

	resp, err = s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	bodyRaw, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	return bodyRaw
}

func (s *mySuite) AddSlotBanner(mySlotBanner SlotBanner) {
	// добавление баннера к слоту
	bodyRaw = s.SendRequest("AddBannerSlot", mySlotBanner)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) DelSlotBanner(mySlotBanner SlotBanner) { // удалени баннера из слота
	bodyRaw = s.SendRequest("DelBannerSlot", mySlotBanner)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) GetBannerForSlot(mySlotSoc ForGetBanner) int { // получения баннера для показа в слоте
	bodyRaw = s.SendRequest("GetBannerForSlot", mySlotSoc)
	s.Require().NotEmpty(bodyRaw)
	err = json.Unmarshal(bodyRaw, &bannerID)
	s.Require().NoError(err)
	return bannerID
}

func (s *mySuite) BannerClick(myBannerClick ForBannerClick) { // кликg по баннеру
	bodyRaw = s.SendRequest("BannerClick", myBannerClick)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) Test1AddBanner() {
	fmt.Println("statrt Test1AddBanner")
	for i := 1; i <= 10; i++ {
		s.AddSlotBanner(SlotBanner{1, i})
	}
	// к слоту привязано 10 баннеров
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1", 10)
	// привязан баннер с id=1
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1 and banner_id=1", 1)
	s.AddSlotBanner(SlotBanner{1, 1})
	// после повторной попытке ничего не изменилось
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1 and banner_id=1", 1)
	fmt.Println("finish Test1AddBanner")
}

func (s *mySuite) Test2DelBanner() {
	fmt.Println("start TestDelSlotBanner")
	//  добавим баннер к слоту
	s.AddSlotBanner(SlotBanner{2, 2})
	// убедимся что он добавился
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2 and banner_id=2", 1)
	// отвяжем баннер от слота
	s.DelSlotBanner(SlotBanner{2, 2})
	// убедимся что отвязался
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2 and banner_id=2", 0)
	fmt.Println("finish TestDelSlotBanner")
}

func (s *mySuite) Test3GetBannerForSlot() {
	fmt.Println("start Test3GetBanner")
	// убедимся, что к слоту 2 не првязан ни один баннер
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2", 0)
	// поскольку к слоту 2 не првязан ни один баннер должен вернуть 0
	bannerID = s.GetBannerForSlot(ForGetBanner{2, 1})
	s.Require().Equal(0, bannerID)

	// добавим во второй слот баннер с ID=3
	s.AddSlotBanner(SlotBanner{2, 3})
	// теперь должен вернуть ID=3 (так как это единственный баннер
	bannerID = s.GetBannerForSlot(ForGetBanner{2, 1})
	s.Require().Equal(3, bannerID)
	// убедимся что этот показ отразился в статистике (1 раз)
	s.CheckCountRec("select count(*) RC from banner_stat where stat_type='S' and slot_id=2 and banner_id=3", 1)
	fmt.Println("finish Test3GetBanner")
}

func (s *mySuite) Test4BannerClick() {
	fmt.Println("start Test4BannerClick")
	// убедимся, что к в слоте 1 для баннера 2 для соц группы 3 ещё не было кликов
	s.CheckCountRec(`select 
	         		 	count(*) RC 
	                 from banner_stat where stat_type='C' 
					   	and slot_id=1 
					   	and banner_id=2 
					   	and soc_group_id=3`, 0)
	//  кликнем в слоте 1 на баннер 2 для соц группы 3
	s.BannerClick(ForBannerClick{1, 2, 3})
	// убедимся, что теперь сохранился 1 клик
	s.CheckCountRec(`select 
						count(*) RC 
					 from banner_stat 
					 where stat_type='C' 
					 	and slot_id=1 
						and banner_id=2 
						and soc_group_id=3`, 1)
	fmt.Println("finish Test4BannerClick")
}

func TestService(t *testing.T) {
	suite.Run(t, new(mySuite))
}
