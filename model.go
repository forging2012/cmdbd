// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`fmt`
	`github.com/jscherff/gocmdb/cmapi`
	`github.com/jscherff/goutil`
)

func SaveDeviceCheckin(dev *cmapi.UsbCi) (err error) {

	vals, err := goutil.ObjectDbValsByCol(dev, `db`, db.Columns[`usbciCheckinInsert`])

	if err != nil {
		elog.Print(err)
		return err
	}

	if _, err := db.Statements[`usbciCheckinInsert`].Exec(vals...); err != nil {
		elog.Print(err)
	}

	return err
}

func GetNewSerialNumber(sfmt string, dev *cmapi.UsbCi) (sn string, err error) {


	vals, err := goutil.ObjectDbValsByCol(dev, `db`, db.Columns[`usbciSnRequestInsert`])

	if err != nil {
		elog.Print(err)
		return sn, err
	}

	res, err := db.Statements[`usbciSnRequestInsert`].Exec(vals...)

	if err != nil {
		elog.Print(err)
		return sn, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		elog.Print(err)
		return sn, err
	}

	sn = fmt.Sprintf(`24F%04X`, id)

	if _, err = db.Statements[`usbciSnRequestUpdate`].Exec(sn, id); err != nil {
		elog.Print(err)
	}

	return sn, err
}

func SaveDeviceChanges(host, vid, pid, sn string, chgs []byte) (err error) {

	if _, err := db.Statements[`usbciChangeInsert`].Exec(host, vid, pid, sn, chgs); err != nil {
		elog.Print(err)
	}

	return err
}

func GetDeviceJSONObject(vid, pid, sn string) (j []byte, err error) {

	if err = db.Statements[`usbciJSONObjectSelect`].QueryRow(vid, pid, sn).Scan(&j); err != nil {
		elog.Print(err)
	}
	return j, err
}

// RowToMap converts a database row into a map of string values indexed by column name.
func RowToMap(vid, pid, sn string) (mss map[string]string, err error) {

	rows, err := db.Statements[`usbciObjectSelect`].Query(vid, pid, sn)

	if err != nil {
		elog.Print(err)
		return nil, err
	}

	defer rows.Close()

	var cols []string

	if cols, err = rows.Columns(); err != nil {
		elog.Print(err)
		return nil, err
	}

	for rows.Next() {

		vals := make([]interface{}, len(cols))
		pvals := make([]interface{}, len(cols))

		for i, _ := range vals {
			pvals[i] = &vals[i]
		}

		if err = rows.Scan(pvals...); err != nil {
			elog.Print(err)
			return nil, err
		}

		mss = make(map[string]string)

		for i, cn := range cols {
			if b, ok := vals[i].([]byte); ok {
				mss[cn] = string(b)
			} else {
				mss[cn] = fmt.Sprintf(`%v`, vals[i])
			}
		}
	}

	if rows.Err() != nil {
		err = rows.Err()
		elog.Print(err)
	}

	return mss, err
}