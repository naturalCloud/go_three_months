package model

import (
	"database/sql"
	"github.com/pkg/errors"
)

//var sqlParseError = errors.New("errorHandle/model: student parse sql error ")

type sqlParseError struct {
	msg string
}

func (e *sqlParseError) Error() string {
	return e.msg
}

type dbQueryQueryError struct {
	msg string
}

func (e *dbQueryQueryError) Error() string {
	return e.msg
}

// IsSqlParseError
// sql 解析错误
func IsSqlParseError(err error) bool {
	var e *sqlParseError
	return errors.As(err, &e)
}

// IsDbQueryError
// db query 错误
func IsDbQueryError(err error) bool {
	var e *dbQueryQueryError
	return errors.As(err, &e)
}

type Student struct {
	Name   string
	Age    int
	Gender int
}

// InsertStudents
// 插入错误
func InsertStudents(db *sql.DB, student Student, sql string) error {
	prepare, err := db.Prepare(sql)

	if err != nil {
		return errors.Wrap(&sqlParseError{
			msg: "errorHandle/model: student parse sql error"}, err.Error())
	}
	_, err = prepare.Exec(student.Name, student.Age, student.Gender)

	return errors.Wrap(err, " prepare.Exec error")

}

// FindStudentByName
// 根据姓名查找学生
func FindStudentByName(db *sql.DB, name string) ([]Student, error) {
	sqlStr := `select name ,age , gender from student where name = $1`
	var err error
	rows, err := db.Query(sqlStr, name)
	if err != nil {
		return nil, errors.Wrap(
			&dbQueryQueryError{msg: "errorHandle/model: dbQuery "},
			err.Error(),
		)
	}

	students := make([]Student, 0)
	defer rows.Close()
	for rows.Next() {
		stu := Student{}
		err = rows.Scan(&stu.Name, &stu.Age, &stu.Gender)
		// 扫描器扫描数据下一行为空,不属于错误
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		students = append(students, stu)
	}
	return students, err
}
