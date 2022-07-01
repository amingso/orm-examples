package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	. "github.com/x-ream/sqlxb"
	"testing"
)

var Db *sqlx.DB
func InitSqlxDB() *sqlx.DB {

	var err interface{}
	Db, err = sqlx.Connect("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True")
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		if Db != nil  {
			Db.Close()
		}
	}
	Db.SetMaxOpenConns(16)
	Db.SetMaxIdleConns(8)

	return Db
}

func TestBulderX(t *testing.T) {

	pet := Pet{}
	cat := Cat{}
	dog := Dog{}


	builder := NewBuilderX(&cat,"c")
	builder.ResultKeys( "c.color"," d.id")
	builder.Eq("p.id", 1)

	subP := Sub()
	subP.ResultKeys("id").Source(&pet)
	builder.SourceBuilder().Sub(subP).Alia("p").JoinOn(LEFT_JOIN,ON("id","c","pet_id"))

	arr := []interface{}{3000,4000,5000,6000,6666}
	sub := Sub()
	sub.ResultKeys("pet_id").Source(&dog).Eq("age",2).In("weight",arr...)
	builder.SourceBuilder().Sub(sub).Alia("d").JoinOn(LEFT_JOIN,ON("id","c","pet_id"))

	builder.SourceBuilder().Source(&cat).Alia("cat").JoinOn(LEFT_JOIN,ON("pet_id","p","id"))
	builder.GroupBy("c.color").
		Having(Gt,"id",1000).
		Sort("p.id",DESC).
		Paged().Rows(10).Last(101)

	vs, dataSql, countSql, kmp:= builder.WithoutOptimization().Build().Sql()
	fmt.Println(dataSql)
	fmt.Println(vs)
	fmt.Println(kmp)
	fmt.Println(countSql)

	InitSqlxDB()

	catList := []Cat{}
	err := Db.Select(&catList, dataSql,vs...)
	if err != nil {
		fmt.Println(err)
	}
	s := fmt.Sprintf("price : %v", *(catList[0].Price))
	fmt.Println(s)
}

