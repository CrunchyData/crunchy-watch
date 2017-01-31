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
	"log"
)

type BackupJobParms struct {
	JOB_HOST           string
	CCP_IMAGE_TAG      string
	CMD                string
	PG_USER            string
	PG_PASSWORD        string
	PG_PORT            string
	BACKUP_PV_CAPACITY string
	BACKUP_PV_PATH     string
	BACKUP_PV_HOST     string
	BACKUP_PVC_STORAGE string
}

type BackupJob struct {
	Logger        *log.Logger
	Host          string
	CCP_IMAGE_TAG string
	Cmd           string
}

func DoSomething() {
}
