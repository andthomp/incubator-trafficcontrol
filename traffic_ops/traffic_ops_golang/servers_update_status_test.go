package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetServerUpdateStatus(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	row := sqlmock.NewRows([]string{"value"})
	row.AddRow("true")
	mock.ExpectQuery("SELECT").WillReturnRows(row)

	serverStatusRow := sqlmock.NewRows([]string{"id", "host_name", "type", "reval_pending", "upd_pending", "status", "parent_upd_pending", "parent_reval_pending"})
	serverStatusRow.AddRow(1, "host_name_1", "EDGE", true, true, "ONLINE", true, false)

	mock.ExpectPrepare("WITH").ExpectQuery().WithArgs("host_name_1").WillReturnRows(serverStatusRow)

	result, err := getServerUpdateStatus("host_name_1", db)

	expected := []tc.ServerUpdateStatus{{"host_name_1", true, true, 1, "ONLINE", true, false}}

	reflect.DeepEqual(expected, result)
}
