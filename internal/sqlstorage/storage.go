package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
)

type Storage struct {
	DBName     string
	DBUserName string
	DBPassword string
	DBConnect  *sql.DB
	Ctx        context.Context
	MyBandit   manyArmedBandit.MyBandit
}

type MyStorage interface {
	Connect() error
	AddBannerSlot(slotId int, bannerId int) error
	DelBannerSlot(slotId int, bannerId int) error
	BannerClick(slotId int, bannerId int, socGroupId int) error
	GetBannerForSlot(slotId int, socGroupId int) (error, int)
	Close() error
}

func rowsToStruct(rows *sql.Rows) ([]manyArmedBandit.BannerStruct, error) {
	var myBannerList []manyArmedBandit.BannerStruct
	var BannerId, ShowCount, ClickCount int
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&BannerId, &ShowCount, &ClickCount); err != nil {
			return nil, err
		}
		myBannerList = append(myBannerList, manyArmedBandit.BannerStruct{BannerId, ShowCount, ClickCount}) //nolint
	}
	return myBannerList, nil
}

func New(ctx context.Context, dbName, dbUserName, dbPassword string, myBandit manyArmedBandit.MyBandit) MyStorage {
	return &Storage{
		DBName: dbName, DBUserName: dbUserName, DBPassword: dbPassword, Ctx: ctx, MyBandit: myBandit,
	}
}

func (s *Storage) Connect() error {
	var err error
	myStr := "user=" + s.DBUserName + " dbname=" + s.DBName + " password=" + s.DBPassword + " sslmode=disable"
	s.DBConnect, err = sql.Open("postgres", myStr)
	if err == nil {
		err = s.DBConnect.PingContext(s.Ctx)
	}
	return err
}

func (s *Storage) AddBannerSlot(slotId int, bannerId int) error {
	fmt.Println("slotId=", slotId)
	query := `
			insert into slot_banner(slot_id,  banner_id)
			values($1, $2)
		`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotId, bannerId)
	fmt.Println(result, err)
	return err
}

func (s *Storage) DelBannerSlot(slotId int, bannerId int) error {
	query := `
			delete from slot_banner
			where slot_id = $1 
			  and  banner_id=$2
		`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotId, bannerId)
	fmt.Println(result, err)
	return err
}

func (s *Storage) BannerClick(slotId int, bannerId int, soc_group_id int) error {
	query := `
			insert into banner_stat(slot_id,  banner_id, soc_group_id, stat_type)
			values($1, $2, $3, 'C')
		`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotId, bannerId, soc_group_id)
	fmt.Println(result, err)
	return err
}

func (s *Storage) GetBannerForSlot(slotId int, socGroupId int) (error, int) {

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
	myStat, errStat := s.DBConnect.QueryContext(s.Ctx, queryStat, socGroupId, slotId)
	if errStat != nil {
		return errStat, 0
	}

	myBannerList, errStruct := rowsToStruct(myStat)
	if errStruct != nil {
		return errStruct, 0
	}

	arrNum := s.MyBandit.GetBannerNum(myBannerList)
	myBannerId := myBannerList[arrNum].BannerId
	query := `
					insert into banner_stat(slot_id,  banner_id, soc_group_id, stat_type)
					values($1, $2, $3, 'S')
	`
	result, err := s.DBConnect.ExecContext(s.Ctx, query, slotId, myBannerId, socGroupId)
	fmt.Println(result, err)

	return err, myBannerId
}

func (s *Storage) Close() error {
	fmt.Println("Start close postgr")
	err := s.DBConnect.Close()
	fmt.Println("Finiah close postgr")
	return err
}
