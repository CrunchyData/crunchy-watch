/*
 Copyright 2017 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package watchapi

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const DEFAULT_HEALTHCHECK_QUERY = "SELECT 1;"
const DEFAULT_HEALTHCHECK_DELAY = 10

func HealthCheck() bool {

	var conn *sql.DB
	var rows *sql.Rows
	var err error

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable", EnvVars.PG_MASTER_SERVICE, EnvVars.PG_MASTER_PORT, EnvVars.PG_MASTER_USER, EnvVars.PG_MASTER_PASSWORD, EnvVars.PG_DATABASE)

	conn, err = sql.Open("postgres", connectionString)

	if err != nil {
		Logger.Printf("[hc] Error creating connection : %s\n", err.Error())
		return false

	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	rows, err = conn.Query(DEFAULT_HEALTHCHECK_QUERY)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		Logger.Println("[hc] failed: error: " + err.Error())
		return false
	}

	return true

}
