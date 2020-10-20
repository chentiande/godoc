package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/bmizerany/pq"

	"os"
	"strconv"
	"time"
)

func test(num int, db *sql.DB, csvChan chan int) {
	tablename := "pmtest" + strconv.Itoa(num)
	create_sql1 := "DROP TABLE " + tablename
	create_sql2 := `CREATE TABLE ` + tablename + `(
		id  serial PRIMARY KEY,
		name1 char(20) DEFAULT NULL,
		name2 char(20) DEFAULT NULL,
		name3 char(20) DEFAULT NULL,
		name4 char(20) DEFAULT NULL,
		date1 timestamp  DEFAULT NULL,
		date2 timestamp  DEFAULT NULL,
		date3 timestamp  DEFAULT NULL,
		date4 timestamp  DEFAULT NULL,
		value1 float DEFAULT NULL,
		value2 float DEFAULT NULL,
		value3 float DEFAULT NULL,
		value4 float DEFAULT NULL
		
	  ) `

	_, err := db.Exec(create_sql1)
	if err != nil {
		fmt.Printf(err.Error())
	}
	_, err = db.Exec(create_sql2)
	if err != nil {
		fmt.Printf(err.Error())
	}

	time1 := time.Now().UnixNano()

	for i := 1; i <= 10; i++ {

		tx, err := db.Begin()
		for j := 1; j <= 10000; j++ {
			a := "字符串字段" + strconv.Itoa((i-1)*10000+j)
			_, err = tx.Exec("insert into "+tablename+" (name1,name2,name3,name4,date1,date2,date3,date4,value1,value2,value3,value4) values ($1,$2,$3,$4,current_timestamp,current_timestamp,current_timestamp,current_timestamp,$5,$6,$7,$8)", a, a, a, a, i, i, i, i)
			if err != nil {
				fmt.Printf(err.Error())
			}
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf(err.Error())
		}

	}
	fmt.Printf("第%v个程序执行插入10万条数据花费：%v毫秒,IO写速度：%vK/s\n", num, (time.Now().UnixNano()-time1)/1000000, 15872*1000000000/(time.Now().UnixNano()-time1))
	var m int64
	m = 1
	for k := 1; k <= 2; k++ {
		time1 = time.Now().UnixNano()
		_, err = db.Exec("insert into " + tablename + " (name1,name2,name3,name4,date1,date2,date3,date4,value1,value2,value3,value4) select name1,name2,name3,name4,date1,date2,date3,date4,value1,value2,value3,value4 from " + tablename)
		if err != nil {
			fmt.Printf(err.Error())
		}
		fmt.Printf("第%v个程序执行第%v次复制数据花费：%v毫秒，IO写速度：%vK/s\n", num, k, (time.Now().UnixNano()-time1)/1000000, 15872*m*1000000000/(time.Now().UnixNano()-time1))
		m = m * 2
	}
	if err != nil {
		fmt.Printf(err.Error())
	}

	time1 = time.Now().UnixNano()
	_, err = db.Exec("create index idx1_" + tablename + " on " + tablename + "(name1)")
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("第%v个程序创建索引花费：%v毫秒\n", num, (time.Now().UnixNano()-time1)/1000000)
	time1 = time.Now().UnixNano()
	for i := 1; i <= 1; i++ {
		_, err = db.Query("select sum(value1),avg(value2),min(value3),max(value4) from  (select * from " + tablename + " ) a")

		if err != nil {
			fmt.Printf(err.Error())
		}
	}
	fmt.Printf("第%v个程序聚合查询花费：%v毫秒\n", num, (time.Now().UnixNano()-time1)/1000000)

	csvChan <- 1
}

func main() {

	var showVer = flag.Bool("v", false, "show build version")

	var ip = flag.String("ip", "192.168.2.101", "pg db ipadress")

	var port = flag.String("port", "5432", "pg db port")

	var user = flag.String("user", "postgres", "pg db user")

	var passwd = flag.String("passwd", "postgre", "pg db passwd")

	var dbname = flag.String("dbname", "postgre", "pg database name")
	var runnum = flag.Int("n", 2, "pg test run number")

	flag.Parse()

	if *showVer {

		fmt.Printf("build name:\t%s\n", "dbbeach")
		fmt.Printf("build ver:\t%s\n", "20201020")
		fmt.Printf("build author:\t%s\n", "chentiande")

		os.Exit(0)
	}
	db, err := sql.Open("postgres", "user="+*user+" password="+*passwd+" host="+*ip+" dbname="+*dbname+"s port="+*port+" sslmode=disable")
	defer db.Close()

	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(*runnum)
	db.SetMaxIdleConns(10)

	var intChan chan int

	intChan = make(chan int, *runnum)

	for i := 1; i <= *runnum; i++ {
		go test(i, db, intChan)

	}
	for i := 1; i <= *runnum; i++ {
		<-intChan

	}

}
