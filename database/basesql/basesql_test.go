package basesql

//
//import (
//	"encoding/json"
//	"fmt"
//	"testing"
//
//	_ "github.com/mattn/go-sqlite3"
//	"gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-config"
//	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/chain/database"
//	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/chain/log"
//)
//
//func getJSON(o interface{}) string {
//	j, err := json.Marshal(o)
//	if err != nil {
//		panic(err)
//	}
//	return string(j)
//}
//
//func getConn(dir, dbname string) *Basesql {
//	cfg := config.NewMapConfig(nil)
//	if dir == "" {
//		dir = "/Users/caojingqi/temp/sqlitedb/"
//	}
//	if dbname == "" {
//		dbname = "test.db"
//	}
//
//	cfg.Set("db_type", "sqlite3")
//	cfg.Set("db_dir", dir)
//
//	bs := &Basesql{}
//	logger := log.Initialize("", "output.log", "err.log")
//	err := bs.Init(dbname, cfg, logger)
//	if err != nil {
//		panic(err)
//	}
//
//	opt, opi, qt, qi := bs.GetInitSQLs()
//	err = bs.PrepareTables(opt, opi)
//	if err != nil {
//		panic(err)
//	}
//	err = bs.PrepareTables(qt, qi)
//	if err != nil {
//		panic(err)
//	}
//
//	return bs
//}
//
//func TestInit(t *testing.T) {
//	_ = getConn("", "")
//}
//
//func TestInsert(t *testing.T) {
//	bs := getConn("", "")
//	defer bs.Close()
//
//	fields := []database.Feild{
//		database.Feild{Name: "accountid", Value: "sfalsfj3454iosafjslfj"},
//		database.Feild{Name: "signer", Value: []byte("xkljsdofsdasdfsdslkfjwoijflsnfsdfjlj")},
//		database.Feild{Name: "weight", Value: 100},
//		database.Feild{Name: "keytype", Value: ""},
//	}
//
//	res, err := bs.Insert(database.TableSigners, fields)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(res.LastInsertId())
//}
//
//func TestDelete(t *testing.T) {
//	bs := getConn("", "")
//	defer bs.Close()
//
//	where := []database.Where{
//		database.Where{Name: "accountid", Value: "sfalsfj3454iosafjslfj"},
//		database.Where{Name: "signer", Value: []byte("xkljsdofsdasdfsdslkfjwoijflsnfsdfjlj")},
//	}
//
//	res, err := bs.Delete(database.TableSigners, where)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(res.RowsAffected())
//}
//
//func TestUpdate(t *testing.T) {
//	bs := getConn("", "")
//	defer bs.Close()
//
//	toupdate := []database.Feild{
//		database.Feild{Name: "weight", Value: 20},
//	}
//
//	where := []database.Where{
//		database.Where{Name: "accountid", Value: "sfalsfj3454iosafjslfj"},
//	}
//
//	res, err := bs.Update(database.TableSigners, toupdate, where)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(res.RowsAffected())
//}
//
//func TestSelectRows(t *testing.T) {
//	bs := getConn("", "")
//	defer bs.Close()
//
//	where := []database.Where{
//		database.Where{Name: "accountid", Value: "sfalsfj3454iosafjslfj"},
//	}
//	order := &database.Order{
//		Type: "desc",
//		Feilds: []string{
//			"signer",
//		},
//	}
//
//	var result []database.Signer
//	err := bs.SelectRows(database.TableSigners, where, order, nil, &result)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for _, v := range result {
//		fmt.Println(v)
//	}
//}
//
//func TestSql(t *testing.T) {
//	bs := getConn("/Users/caojingqi/temp/sqlitedb/", "test.db")
//	defer bs.Close()
//
//	values := []interface{}{20}
//
//	var result []database.Offer
//	err := bs.conn.Select(&result, "select * from offers where 1 = 1  order by price asc, offerid asc limit ?", values...)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for i := range result {
//		fmt.Println(getJSON(result[i]))
//		// fmt.Println(result[i].SellingAssetCode, result[i].SellingAmount, result[i].BoughtAssetCode, result[i].BoughtAmount)
//	}
//
//}
