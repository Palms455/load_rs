package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"load_rs/internal/config"
	"load_rs/internal/xml_parse"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var pool *pgxpool.Pool
var err error

var Flc_chan = make(chan int, 1000)

func init() {
	log.Println("Connected to : ", config.MainConfig.DbAddr)
	pool, err = pgxpool.Connect(context.Background(), config.MainConfig.DbAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}
}

type RsFile struct {
	Filename, Ftype, Mo, Tfoms, Period, Nn, Rsfile, Status string
	Rs_data xml_parse.RsData
	ErrorMsg map[string]string
}

func (r *RsFile) GetFromFile(file string) {
	rs := Rgx.FindStringSubmatch(file)
	if len(rs) == 0 {
		log.Fatalf("Файл %s не подходит для обработки", file)
	}
	r.Filename, r.Ftype, r.Mo, r.Tfoms, r.Period, r.Nn = rs[0], rs[1], rs[2], rs[3], rs[4], rs[5]
}

func (r *RsFile) GetRsFilePath(year, month int) {
	r.Rsfile = filepath.Join("for_tfoms", strconv.Itoa(year), fmt.Sprintf("%02d", month), r.Filename)
}


func GetCurrentPeriod(period_match string) (int, int, error) {
	var year, month int
	err = pool.QueryRow(context.Background(),
		"select year, month from rsj.period where id = $1 and is_open=True limit 1", period_match).Scan(&year, &month)
	if err != nil {
		log.Printf("Ошибка при получении периода: %s", err)
		return year, month, err
	}
	return year, month, err
}

func (r *RsFile) LoadPers(rs_id int) error {
	for _, pers := range r.Rs_data.LFile.PERS {
		json_pers, _ := json.Marshal(pers)
		_, err = pool.Exec(context.Background(),
			`insert into rsj.inbufpers (pers, rsfile_id) values ($1, $2)`,
			json_pers, rs_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RsFile) LoadZap(rs_id int) error {
	for _, zap := range r.Rs_data.Rs.ZAP {
		json_zap, _ := json.Marshal(zap)
		_, err = pool.Exec(context.Background(),
			`insert into rsj.inbufzap (zap, rsfile_id) values ($1, $2)`,
			json_zap, rs_id)
		if err != nil {
			return err
		}
	}
	return nil
}


func (r *RsFile) DeleteRs() error {
	var rs_id int
	err := pool.QueryRow(context.Background(),
		"select rs.id from rsj.inbufrs rs where lower(rs.filename) = lower($1)",
		r.Filename).Scan(&rs_id)
	if err != nil {
		log.Printf("Ошибка при получении id в rsj.inbufrs: %s", err)
	}
	if rs_id == 0 {
		return nil
	}
	fmt.Println(rs_id)
	_, err = pool.Exec(context.Background(), "delete from rsj.inbufpers ps where ps.rsfile_id = $1", rs_id)
	if err != nil {
		log.Println("del inbufpers")
		log.Println(err)
	}
	_, err = pool.Exec(context.Background(), "delete from rsj.inbufzap zap where zap.rsfile_id = $1", rs_id)
	if err != nil {
		log.Println("del inbufzap")
		log.Println(err)
	}
	_, err = pool.Exec(context.Background(), "delete from rsj.log_time_flc flc where flc.rsfile_id = $1", rs_id)
	if err != nil {
		log.Println(err)
	}
	_, err = pool.Exec(context.Background(), "delete from rsj.inbufrs rs where rs.id = $1", rs_id)
	if err != nil {
		log.Println("del inbufrs")
		log.Println(err)
	}
	return nil
}

func (r *RsFile) LoadRs() error {
	var rs_id int
	if err := r.DeleteRs(); err != nil {
		log.Fatal(err)
		return err
	}

	if r.Status == "20" {
		json_zglv, _ := json.Marshal(r.Rs_data.Rs.ZGLV)
		json_schet, _ := json.Marshal(r.Rs_data.Rs.SCHET)
		json_zglv_p, _ := json.Marshal(r.Rs_data.LFile.ZGLV)
		err = pool.QueryRow(context.Background(),
			`insert into rsj.inbufrs
					(filename, ftype, mo, tfoms, period_id, nn, ts_insert, zglv, schet, zglv_p, status, rsfile)
				values ($1, $2, $3, $4, $5, $6, now(), $7, $8, $9, 0, $10) returning id`,
			r.Filename,
			r.Ftype,
			r.Mo,
			r.Tfoms,
			r.Period,
			r.Nn,
			json_zglv,
			json_schet,
			json_zglv_p,
			r.Rsfile).Scan(&rs_id)
		if err != nil {
			log.Fatalf("Ошибка при вставке в rsj.inbufrs %s", err)
			return err
		}
		if err := r.LoadPers(rs_id); err != nil {
			log.Fatalf("Ошибка при вставке в rsj.inbufpers %s", err)
			return err
		}
		if err := r.LoadZap(rs_id); err != nil {
			log.Fatalf("Ошибка при вставке в rsj.inbufzap %s", err)
			return err
		}
		_, err = pool.Exec(context.Background(), "update rsj.inbufrs set status = 20 where id = $1", rs_id)
		if err != nil {
			log.Fatalf("Ошибка при обновлении rsj.inbufrs %s", err)
			return err
		}
		Flc_chan <- rs_id
	} else {
		json_error, _ := json.Marshal(r.ErrorMsg)
		err = pool.QueryRow(context.Background(),
			`insert into rsj.inbufrs
					(filename, ftype, mo, tfoms, period_id, nn, ts_insert, status, rsfile, errors)
				values ($1, $2, $3, $4, $5, $6, now(), $7, $8, $9) returning id`,
			r.Filename,
			r.Ftype,
			r.Mo,
			r.Tfoms,
			r.Period,
			r.Nn,
			r.Status,
			r.Rsfile,
			json_error).Scan(&rs_id)
		if err != nil {
			log.Fatalf("Ошибка при вставке в rsj.inbufrs %s", err)
			return err
		}
	}
	return nil
}


func FLC(file_id int) {
	start := time.Now()
	_, err = pool.Exec(context.Background(), "select rsj.process_flc($1)", file_id)
	if err != nil {
		log.Fatalf("Ошибка при проведении ФЛК %s", err)
	}
	end := time.Since(start)
	log.Println(file_id)
	log.Println("Время ФЛК для файла: ", end)
}