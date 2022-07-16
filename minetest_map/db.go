package main

import (
	"bytes"
	"database/sql"
	"github.com/anon55555/mt"
	_ "github.com/mattn/go-sqlite3" // MIT licensed.
	"log"

	"github.com/EliasFleckenstein03/mtmap"
)

var db *sql.DB

var writeBlk *sql.Stmt
var readBlk *sql.Stmt

func OpenDB(file string) (err error) {
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal("cant open map.sqlite: ", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `blocks` (	`pos` INT PRIMARY KEY, `data` BLOB );")
	if err != nil {
		log.Fatal("cant create table blocks: ", err)
	}

	// prepare stms:
	// writeBlk, err = db.Prepare("UPSERT INTO blocks (pos, param0, param1, param2) VALUES (?, ?, ?, ?)")
	readBlk, err = db.Prepare("SELECT data FROM blocks WHERE pos = ?")
	if err != nil {
		log.Fatal("cant prepare read statement: ", err)
	}

	// return the return value to the caller
	// this is important so the caller can check whether
	// an error has occured during the execution of the function
	return
}

func GetBlk(p [3]int16) *mtmap.MapBlk {
	r := readBlk.QueryRow(Blk2DBPos(p))

	var buf []byte
	r.Scan(&buf)
	if len(buf) == 0 {
		return nil
	}

	reader := bytes.NewReader(buf)

	blk, err := mtmap.Deserialize(reader, nimap)
	if err != nil {
		log.Println("error", err)
	}

	return blk
}

func SetNode(pos [3]int16, node mt.Content) {
	blk, i := mt.Pos2Blkpos(pos)
	oldBlk := GetBlk(blk)

	if oldBlk == nil {
		oldBlk = EmptyBlk()
	}

	oldBlk.Param0[i] = node

	SetBlk(blk, oldBlk)
}

func SetBlk(p [3]int16, blk *mtmap.MapBlk) {
	q, err := db.Prepare("INSERT OR REPLACE INTO blocks (pos, data) VALUES (?, ?)")
	if err != nil {
		log.Fatal("can't set block: ", err)
	}

	defer q.Close()

	w := &bytes.Buffer{}

	err = mtmap.Serialize(blk, w, nimap)
	if err != nil {
		panic(err)
	}

	pos := Blk2DBPos(p)

	_, err = q.Exec(pos, w.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func Commit() {
}

func EmptyBlk() (blk *mtmap.MapBlk) {
	for k := range blk.Param0 {
		blk.Param0[k] = mt.Air
	}

	return
}
