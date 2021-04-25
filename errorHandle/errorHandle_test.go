package errorHandle

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go_three_months/errorHandle/model"
	"testing"
)

var (
	driver = "sqlite3"
	dbPath = "db/student.sqllist"
	db     *sql.DB
)

func TestDbModelNotFound(t *testing.T) {
	tmpDb, err := sql.Open(driver, dbPath)
	if err != nil {
		panic(err)
	}

	db = tmpDb
	// 错误sql
	err = insert(`insert into  (name , age , gender) values(?,?,?)`)
	if model.IsSqlParseError(err) {
		t.Errorf("%+v", err)
	}
	// 正确sql
	err = insert(`insert into  student (name , age , gender) values(?,?,?)`)
	if err == nil {
		t.Log("\n", "write data ok")
	}

	// 查询数据 存在
	students, err := query(`张三`)
	if model.IsDbQueryError(err) {
		t.Errorf("%+v", err)
	} else {
		t.Log(students)
	}

	// 查询数据不存在
	students, err = query("张三丰")
	if !model.IsDbQueryError(err) {
		t.Log(students)
	} else {
		t.Errorf("%+v", err)
	}

}

// 创建一些数据
func insert(sql string) error {
	return model.InsertStudents(db, model.Student{
		Name:   "张三",
		Age:    3,
		Gender: 0,
	}, sql)
}

func query(name string) ([]model.Student, error) {
	return model.FindStudentByName(db, name)
}
