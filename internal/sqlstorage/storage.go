package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
)

type DBConf struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUserName string
	DBPassward string
}

type Storage struct {
	DBConf    DBConf
	DBConnect *sql.DB
	Ctx       context.Context
	MyBandit  manyarmedbandit.MyBandit
}

type BannerStatStruct struct {
	ID         int
	SlotID     int
	BannerID   int
	SocGroupID int
	StatType   string
	RecDate    string
}

func New(ctx context.Context, myDBConf DBConf, myBandit manyarmedbandit.MyBandit) *Storage {
	return &Storage{
		DBConf: myDBConf, Ctx: ctx, MyBandit: myBandit,
	}
}

func rowsToStruct(rows *sql.Rows) ([]manyarmedbandit.BannerStruct, error) {
	var myBannerList []manyarmedbandit.BannerStruct
	var bannerID, ShowCount, ClickCount int
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&bannerID, &ShowCount, &ClickCount); err != nil {
			return nil, err
		}
		myBanner := manyarmedbandit.BannerStruct{
			BannerID:   bannerID,
			ShowCount:  ShowCount,
			ClickCount: ClickCount,
		}
		myBannerList = append(myBannerList, myBanner)
	}
	return myBannerList, nil
}

func rawsToStat(raws *sql.Rows) ([]BannerStatStruct, error) {
	var myBannerList []BannerStatStruct
	var id, slotID, bannerID, socGroupID int
	var statType, recDate string
	defer raws.Close()
	for raws.Next() {
		if err := raws.Scan(&id, &slotID, &bannerID, &socGroupID, &statType, &recDate); err != nil {
			return nil, err
		}
		myBannerList = append(myBannerList, BannerStatStruct{id, slotID, bannerID, socGroupID, statType, recDate})
	}
	return myBannerList, nil
}

func (s *Storage) Connect() error {
	var err error
	myStr := "postgres://" + s.DBConf.DBUserName + ":" + s.DBConf.DBPassward + "@"
	myStr += s.DBConf.DBHost + ":" + s.DBConf.DBPort + "/" + s.DBConf.DBName + "?sslmode=disable"
	fmt.Println("start connect to postgrace:", myStr)
	s.DBConnect, err = sql.Open("postgres", myStr)
	if err == nil {
		err = s.DBConnect.PingContext(s.Ctx)
	}
	return err
}

func (s *Storage) AddBannerSlot(slotID int, bannerID int) error {
	fmt.Println("slotID=", slotID)
	query := `
			insert into slot_banner(slot_id,  banner_id)
			values($1, $2)
		`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotID, bannerID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) DelBannerSlot(slotID int, bannerID int) error {
	query := `
			delete from slot_banner
			where slot_id = $1 
			  and  banner_id=$2
		`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotID, bannerID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) BannerClick(slotID int, bannerID int, socGroupID int) error {
	query := `
			insert into banner_stat(slot_id,  banner_id, soc_group_id, stat_type)
			values($1, $2, $3, 'C')
		`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotID, bannerID, socGroupID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) GetBannerForSlot(slotID int, socGroupID int) (int, error) {
	queryStat := `
		select  sb.banner_id, count(distinct bs_s.id) show_count, count(distinct bs_c.id) click_count
		from slot_banner sb
		left join banner_stat bs_s
			on sb.slot_id=bs_s.slot_id
			and sb.banner_id=bs_s.banner_id
			and bs_s.soc_group_id=$1
			and bs_s.stat_type = 'S'
		left join banner_stat bs_c
			on sb.slot_id=bs_c.slot_id
			and sb.banner_id=bs_c.banner_id
			and bs_c.soc_group_id=$1
			and bs_c.stat_type = 'C'
		where sb.slot_id=$2 
		group by sb.banner_id;
	`
	myStat, errStat := s.DBConnect.QueryContext(s.Ctx, queryStat, socGroupID, slotID)
	if errStat != nil {
		return 0, errStat
	}

	myBannerList, errStruct := rowsToStruct(myStat)
	if errStruct != nil {
		return 0, errStruct
	}
	fmt.Println(myBannerList)
	if len(myBannerList) == 0 {
		return 0, nil
	}

	arrNum := s.MyBandit.GetBannerNum(myBannerList)
	myBannerID := myBannerList[arrNum].BannerID
	query := `
					insert into banner_stat(slot_id,  banner_id, soc_group_id, stat_type)
					values($1, $2, $3, 'S')
	`
	_, err := s.DBConnect.ExecContext(s.Ctx, query, slotID, myBannerID, socGroupID)
	return myBannerID, err
}

func (s *Storage) GetBannerStat() ([]BannerStatStruct, error) {
	queryStat := `
		select id, slot_id, banner_id, soc_group_id, stat_type, rec_date
		from banner_stat
		where id>(select max(banner_stat_id) from send_stat_max_id)	
	`
	myStat, errStat := s.DBConnect.QueryContext(s.Ctx, queryStat)
	if errStat != nil {
		return nil, errStat
	}

	myBannerStatList, errStruct := rawsToStat(myStat)
	if errStruct != nil {
		return nil, errStruct
	}
	return myBannerStatList, nil
}

func (s *Storage) ChangeSendStatID(id int) error {
	query := `
			update send_stat_max_id
			set banner_stat_id = $1
		`
	_, err := s.DBConnect.ExecContext(s.Ctx, query, id)
	return err
}

func (s *Storage) Close() error {
	fmt.Println("Start close postgr")
	err := s.DBConnect.Close()
	fmt.Println("Finiah close postgr")
	return err
}
