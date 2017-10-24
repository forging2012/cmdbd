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

import `fmt`

// SaveDeviceCheckin saves a device checkin to the database 'checkins' table.
func SaveDeviceCheckin(dev map[string]interface{}) (err error) {

	var vals []interface{}

	for _, col := range qy.Cols[`usbCiInsertCheckin`] {
		vals = append(vals, dev[col])
	}

	_, err = qy.Stmt[`usbCiInsertCheckin`].Exec(vals...)

	return err
}

// GetNewSerialNumber generates a new device serial number using the value
// from the auto-incremented ID column of the 'snrequest' table with the
// format string provided by the caller.
func GetNewSerialNumber(dev map[string]interface{}) (sn string, err error) {

	var (
		id int64
		serialFmt string
		vals []interface{}
	)

	if objectType, ok := dev[`object_type`].(string); !ok {
		return sn, fmt.Errorf(`cannot determine device type`)
	} else if serialFmt, ok = conf.SerialFmt[objectType]; !ok {
		return sn, fmt.Errorf(`missing SN format for device %q`, objectType)
	}

	tx, err := db.Begin()

	if err != nil {
		return sn, err
	}

	for _, col := range qy.Cols[`usbCiInsertSnRequest`] {
		vals = append(vals, dev[col])
	}

	if res, err := qy.Stmt[`usbCiInsertSnRequest`].Exec(vals...); err != nil {
		return sn, err
	} else if id, err = res.LastInsertId(); err != nil {
		return sn, err
	}

	if res, err := qy.Stmt[`cmdbInsertSequence`].Exec(); err != nil {
		return sn, err
	} else if seq, err := res.LastInsertId(); err != nil {
		return sn, err
	} else {
		sn = fmt.Sprintf(serialFmt, seq)
	}

	if _, err := qy.Stmt[`usbCiUpdateSnRequest`].Exec(sn, id); err != nil {
		return sn, err
	}

	err = tx.Commit()
	return sn, err
}

// SaveDeviceChanges records changes reported in a device audit in the 'changes'
// table in JSON format.
func SaveDeviceChanges(host, vid, pid, sn string, chgs []byte) (err error) {

	_, err = qy.Stmt[`usbCiInsertChanges`].Exec(host, vid, pid, sn, chgs)
	return err
}

// GetDeviceJSONObject retreives device properties from the 'serialized' device
// table and returns them to the caller in JSON format.
func GetDeviceJSONObject(vid, pid, sn string) (j []byte, err error) {

	err = qy.Stmt[`usbCiSelectJSONObject`].QueryRow(vid, pid, sn).Scan(&j)
	return j, err
}

// SaveUsbMeta updates the USB meta tables in the database.
func SaveUsbMeta() (err error) {

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	vendorStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceVendor`])
	productStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceProduct`])
	classStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceClass`])
	subclassStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceSubClass`])
	protocolStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceProtocol`])

	for vid, v := range conf.MetaUsb.Vendors {

		if _, err := vendorStmt.Exec(vid, v.String()); err != nil {
			return err
		}

		for pid, p := range v.Product {

			if _, err := productStmt.Exec(vid, pid, p.String()); err != nil {
				return err
			}
		}
	}

	for cid, c := range conf.MetaUsb.Classes {

		if _, err := classStmt.Exec(cid, c.String()); err != nil {
			return err
		}

		for sid, s := range c.SubClass {

			if _, err := subclassStmt.Exec(cid, sid, s.String()); err != nil {
				return err
			}

			for pid, p := range s.Protocol {

				if _, err := protocolStmt.Exec(cid, sid, pid, p.String()); err != nil {
					return err
				}
			}
		}
	}

	err = tx.Commit()
	return err
}
