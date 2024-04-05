package radix

import (
	"database/sql"

	"github.com/overalice/radix/dialect"
)

type Orm struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewOrm(driver, source string) (orm *Orm, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		Error(err.Error())
		return
	}
	if err = db.Ping(); err != nil {
		Error(err.Error())
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		Error("dialect %s Not Found", driver)
		return
	}
	orm = &Orm{db: db, dialect: dial}
	Info("Connect database success")
	return
}

func (orm *Orm) NewSession() *Session {
	return NewSession(orm.db, orm.dialect)
}

func (orm *Orm) Close() {
	if err := orm.db.Close(); err != nil {
		Error("Failed to close database")
	}
	Info("Close database success")
}
