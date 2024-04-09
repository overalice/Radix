package radix

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/overalice/radix/dialect"
)

type User struct {
	ID   int `orm:"PRIMARY KEY"`
	Name string
	Age  int
}

var (
	user1 = &User{1, "Tom", 18}
	user2 = &User{2, "Sam", 25}
	user3 = &User{3, "Jack", 25}
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

func TestSqlite(t *testing.T) {
	Info("Start test sqlite")
	orm, _ := NewOrm("sqlite3", "radix.db")
	defer orm.Close()
	s := orm.NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	Info("End test sqlite")
}

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession(TestDB, TestDial).Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession(TestDB, TestDial).Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
