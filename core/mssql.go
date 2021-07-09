package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type EnvInfo struct {
	RegionId    string `field:"RegionId"`
	Mrc         string `field:"Mrc"`
	DbServer    string `field:"DbServer"`
	DaxDbName   string `field:"DaxDbName"`
	StageDbName string `field:"StageDbName"`
	Caption     string `field:"Caption"`
}

type MsSqlConn struct {
	SqlConn *sql.DB
}

var connString string

// Конструктов подключения к SQL
func NewMsSqlConn(service *Service) *MsSqlConn {
	userName, nu := os.LookupEnv("USER_NAME")
	userPassword, np := os.LookupEnv("USER_PASSWORD")

	if runtime.GOOS == "windows" {
		service.Stdlog.Print("Среда выполнения Windows.")
		if nu && np {
			service.Stdlog.Print("Производиться подключение к БД с учетными данными взятыми из переменных окружения.")
			connString = fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;Connect Timeout=600", service.Cfg.MonitoringDb.Server, service.Cfg.MonitoringDb.DbName, userName, userPassword)
		} else {
			service.Stdlog.Print("Производиться подключение к БД с учетными данными пользователя Windows.")
			connString = fmt.Sprintf("Server=%s;Database=%s;Trusted_Connection=True;MultipleActiveResultSets=true;Connect Timeout=600", service.Cfg.MonitoringDb.Server, service.Cfg.MonitoringDb.DbName)
		}
	} else {
		if nu && np {
			service.Stdlog.Print("Среда выполнения Linux. \n Производиться подключение к БД с учетными данными взятыми из переменных окружения.")
			connString = fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;Connect Timeout=600", service.Cfg.MonitoringDb.Server, service.Cfg.MonitoringDb.DbName, userName, userPassword)
		} else {
			service.Errlog.Fatal("Среда выполнения Linux. \n Не удалось обнаружить обну или обед переменные среды, сожержащие логин и пароль пользователя для подключения к БД.")
		}
	}

	cc, err := sql.Open("mssql", connString)

	if err != nil {
		service.Errlog.Fatal(fmt.Sprintf("Не удалось построить подклчюение к серверу %v с базой данных %v. Ошибка: \n %v", service.Cfg.MonitoringDb.Server, service.Cfg.MonitoringDb.DbName, err))
	}

	return &MsSqlConn{
		SqlConn: cc,
	}
}

// Получить инфомацию о окружениях из таблицы EnvInfo
func (c *MsSqlConn) SelectEnvInfo(query string) ([]EnvInfo, error) {
	ref := []EnvInfo{}
	rows, err := c.SqlConn.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		envinfo := EnvInfo{}

		s := reflect.ValueOf(&envinfo).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err := rows.Scan(columns...)
		if err != nil {
			log.Fatal(err)
		}
		ref = append(ref, envinfo)
	}

	defer rows.Close()

	return ref, err
}

// Вызов процедуры мердж для сбора информации о терминалах
func (c *MsSqlConn) MergeQuery(query string, envinfo EnvInfo) error {

	query = strings.Replace(query, "%regionId%", envinfo.RegionId, -1)
	query = strings.Replace(query, "%server%", envinfo.DbServer, -1)
	query = strings.Replace(query, "%stageDbName%", envinfo.StageDbName, -1)
	fmt.Printf("Старт выборки информации по терминалам с региона %v.\n", envinfo.RegionId)
	res, err := c.SqlConn.Exec(query)
	if err != nil {
		return err
	}
	fmt.Printf("Выборка с региона %v окончена. Произведе мердж в базу данных мониторинга.\n", envinfo.RegionId)
	res.RowsAffected()
	return err
}
