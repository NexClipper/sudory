package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro"
	dctv2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	env "github.com/NexClipper/sudory/pkg/server/model/service/gen_test_data/env"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	new_service := func(name string, cluster string, template string, stepcount int) (*servicev3.Service_create, []servicev3.ServiceStep_create) {
		uuid := macro.NewUuidString()
		created := time.Now()

		service := servicev3.Service_create{}
		service.PartitionDate = TrimDate(created)
		service.ClusterUuid = cluster
		service.Uuid = uuid
		service.Timestamp = created
		service.Name = name
		service.Summary = *vanilla.NewNullString(fmt.Sprintf("service %v", name))
		service.TemplateUuid = template
		service.StepCount = stepcount
		service.Priority = 0
		service.SubscribedChannel = vanilla.NullString{}
		service.StepPosition = 0
		service.Status = servicev3.StepStatusRegist
		service.Created = created

		step := servicev3.ServiceStep_create{}
		step.PartitionDate = TrimDate(created)
		step.ClusterUuid = cluster
		step.Uuid = uuid
		step.Sequence = 0
		step.Timestamp = created
		step.Name = name
		step.Summary = *vanilla.NewNullString(fmt.Sprintf("step %v", name))
		step.Method = "fake.test.method"
		step.Args = dctv2.CryptoObject{"name": "fake-name", "ns": "fake-ns"}
		step.Status = servicev3.StepStatusRegist
		step.ResultFilter = vanilla.NullString{}
		step.Created = created

		return &service, []servicev3.ServiceStep_create{step}
	}

	new_service_status_success := func(service servicev3.Service_create, client_uuid string) ([]servicev3.Service_update, []servicev3.ServiceStep_update) {
		// uuid := macro.NewUuidString()
		updated := time.Now()
		addtime := func() time.Time {
			updated = updated.Add(time.Duration(rand.Intn(1000)+1) * time.Millisecond)
			return updated
		}

		status1 := servicev3.Service_update{}
		status1.Status = 1
		status1.StepPosition = 1
		status1.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status1.Timestamp = addtime()

		status2 := servicev3.Service_update{}
		status2.Status = 2
		status2.StepPosition = 1
		status2.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status2.Timestamp = addtime()

		status4 := servicev3.Service_update{}
		status4.Status = 4
		status4.StepPosition = 1
		status4.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status4.Timestamp = addtime()

		service_status := []servicev3.Service_update{
			status1,
			status2,
			status4,
		}

		step_status1 := servicev3.ServiceStep_update{}
		step_status1.Status = 1
		step_status1.Started = *vanilla.NewNullTime(addtime())
		step_status1.Ended = *vanilla.NewNullTime(addtime())
		step_status1.Timestamp = addtime()

		step_status2 := servicev3.ServiceStep_update{}
		step_status2.Status = 2
		step_status2.Started = *vanilla.NewNullTime(addtime())
		step_status2.Ended = *vanilla.NewNullTime(addtime())
		step_status2.Timestamp = addtime()

		step_status4 := servicev3.ServiceStep_update{}
		step_status4.Status = 4
		step_status4.Started = *vanilla.NewNullTime(addtime())
		step_status4.Ended = *vanilla.NewNullTime(addtime())
		step_status4.Timestamp = addtime()

		service_step_status := []servicev3.ServiceStep_update{
			step_status1,
			step_status2,
			step_status4,
		}

		return service_status, service_step_status
	}

	new_service_status_none := func(service servicev3.Service_create, client_uuid string) ([]servicev3.Service_update, []servicev3.ServiceStep_update) {

		return []servicev3.Service_update{}, []servicev3.ServiceStep_update{}
	}

	new_service_status_progress := func(service servicev3.Service_create, client_uuid string) ([]servicev3.Service_update, []servicev3.ServiceStep_update) {
		// uuid := macro.NewUuidString()
		updated := time.Now()
		addtime := func() time.Time {
			updated = updated.Add(time.Duration(rand.Intn(1000)+1) * time.Millisecond)
			return updated
		}

		status1 := servicev3.Service_update{}
		status1.Status = 1
		status1.StepPosition = 1
		status1.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status1.Timestamp = addtime()

		status2 := servicev3.Service_update{}
		status2.Status = 2
		status2.StepPosition = 1
		status2.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status2.Timestamp = addtime()

		service_status := []servicev3.Service_update{
			status1,
			status2,
		}

		step_status1 := servicev3.ServiceStep_update{}
		step_status1.Status = 1
		step_status1.Started = *vanilla.NewNullTime(addtime())
		step_status1.Ended = *vanilla.NewNullTime(addtime())
		step_status1.Timestamp = addtime()

		step_status2 := servicev3.ServiceStep_update{}
		step_status2.Status = 2
		step_status2.Started = *vanilla.NewNullTime(addtime())
		step_status2.Ended = *vanilla.NewNullTime(addtime())
		step_status2.Timestamp = addtime()

		service_step_status := []servicev3.ServiceStep_update{
			step_status1,
			step_status2,
		}

		return service_status, service_step_status
	}

	new_service_status_error := func(service servicev3.Service_create, client_uuid string) ([]servicev3.Service_update, []servicev3.ServiceStep_update) {
		// uuid := macro.NewUuidString()
		updated := time.Now()
		addtime := func() time.Time {
			updated = updated.Add(time.Duration(rand.Intn(1000)+1) * time.Millisecond)
			return updated
		}

		status1 := servicev3.Service_update{}
		status1.Status = 1
		status1.StepPosition = 1
		status1.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status1.Timestamp = addtime()

		status2 := servicev3.Service_update{}
		status2.Status = 2
		status2.StepPosition = 1
		status2.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status2.Timestamp = addtime()

		status8 := servicev3.Service_update{}
		status8.Status = 8
		status8.StepPosition = 1
		status8.AssignedClientUuid = *vanilla.NewNullString(client_uuid)
		status8.Timestamp = addtime()
		status8.Message = *vanilla.NewNullString(fmt.Sprintf("error: created=%v", updated))

		service_status := []servicev3.Service_update{
			status1,
			status2,
			status8,
		}

		step_status1 := servicev3.ServiceStep_update{}
		step_status1.Status = 1
		step_status1.Started = *vanilla.NewNullTime(addtime())
		step_status1.Ended = *vanilla.NewNullTime(addtime())
		step_status1.Timestamp = addtime()

		step_status2 := servicev3.ServiceStep_update{}
		step_status2.Status = 2
		step_status2.Started = *vanilla.NewNullTime(addtime())
		step_status2.Ended = *vanilla.NewNullTime(addtime())
		step_status2.Timestamp = addtime()

		step_status8 := servicev3.ServiceStep_update{}
		step_status8.Status = 8
		step_status8.Started = *vanilla.NewNullTime(addtime())
		step_status8.Ended = *vanilla.NewNullTime(addtime())
		step_status8.Timestamp = addtime()

		service_step_status := []servicev3.ServiceStep_update{
			step_status1,
			step_status2,
			step_status8,
		}

		return service_status, service_step_status
	}

	new_service_result := func(service servicev3.Service_create) *servicev3.ServiceResult_create {
		created := time.Now()

		r := rand.Int63()

		service_result := &servicev3.ServiceResult_create{}
		service_result.ClusterUuid = service.ClusterUuid
		service_result.Uuid = service.Uuid
		service_result.PartitionDate = service.PartitionDate
		service_result.Timestamp = created.Add(time.Duration(rand.Intn(1000)) * time.Millisecond)
		service_result.ResultSaveType = servicev3.ResultSaveTypeDatabase
		service_result.Result = dctv2.CryptoString(fmt.Sprintf("%v", r))

		return service_result
	}

	laps := time.Now()
	st := time.Now()
	var atomic_count int32
	task := func(db *sql.DB, client_uuid string) func(count int) {

		fn := func(i int) error {

			name := fmt.Sprintf("service-%v", i)
			cluster := fmt.Sprintf("cluster-%v", client_uuid[0:5])
			template := "00000000000000000000000000000000"

			service, step := new_service(name, cluster, template, 1)

			var service_status []servicev3.Service_update
			var step_status []servicev3.ServiceStep_update
			switch i % 10 {
			case 1:
				service_status, step_status = new_service_status_none(*service, client_uuid)
			case 2:
				service_status, step_status = new_service_status_progress(*service, client_uuid)
			case 8:
				service_status, step_status = new_service_status_error(*service, client_uuid)
			default:
				service_status, step_status = new_service_status_success(*service, client_uuid)
			}
			service_result := new_service_result(*service)

			stmt_insert_service := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
				service.TableName(),
				strings.Join(service.ColumnNames(), ", "),
				strings.Join(Repeat(len(service.ColumnNames()), "?"), ", "),
			)

			_, err := db.Exec(stmt_insert_service, service.Values()...)
			if err != nil {
				return err
			}

			if count := atomic.AddInt32(&atomic_count, 1); (count)%1000 == 999 {
				fmt.Printf("since=%v\tdelta=%v\tclient_uuid=%v\tcount=%v\n", time.Since(st).String(), time.Since(laps), client_uuid, count+1)
				laps = time.Now()
			}

			for i := range step {
				stmt_insert_step := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
					step[i].TableName(),
					strings.Join(step[i].ColumnNames(), ", "),
					strings.Join(Repeat(len(step[i].ColumnNames()), "?"), ", "),
				)

				_, err := db.Exec(stmt_insert_step, step[i].Values()...)
				if err != nil {
					return err
				}
			}

			if true {

				for _, service_status := range service_status {

					// get target record
					record := new(servicev3.Service)
					query := fmt.Sprintf(`
SELECT %v 
  FROM %v 
 WHERE pdate = ? AND cluster_uuid = ? AND uuid = ?`,
						strings.Join(record.ColumnNames(), ", "),
						record.TableName(),
					)

					// fmt.Printf("%v\n", query)
					// fmt.Printf("%v\n", service.Values())
					rows, err := db.Query(query, []interface{}{service.PartitionDate, service.ClusterUuid, service.Uuid}...)
					if err != nil {
						return err
					}

					for rows.Next() {
						tmp := new(servicev3.Service)
						if err := tmp.Scan(rows); err != nil {
							return err
						}

						// get last revision
						if record.Timestamp.Before(tmp.Timestamp) {
							record = tmp
						}
					}

					record.Timestamp = service_status.Timestamp
					// reset update columns
					record.AssignedClientUuid = service_status.AssignedClientUuid
					record.StepPosition = service_status.StepPosition
					record.Status = service_status.Status
					record.Message = service_status.Message

					insert := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
						record.TableName(),
						strings.Join(record.ColumnNames(), ", "),
						strings.Join(Repeat(len(record.ColumnNames()), "?"), ", "),
					)

					// fmt.Printf("%v\n", insert)
					// fmt.Printf("%v\n", record.Values())
					if _, err := db.Exec(insert, record.Values()...); err != nil {
						return err
					}
				}

				for _, step := range step {
					for _, step_status := range step_status {

						// get target record
						record := new(servicev3.ServiceStep)
						query := fmt.Sprintf(`
SELECT %v 
  FROM %v 
  WHERE pdate = ? AND cluster_uuid = ? AND uuid = ? AND seq = ?`,
							strings.Join(record.ColumnNames(), ", "),
							record.TableName(),
						)

						// fmt.Printf("%v\n", query)
						// fmt.Printf("%v %v %v %v\n", step.PartitionDate, step.ClusterUuid, step.Uuid, step.Sequence)
						rows, err := db.Query(query, []interface{}{step.PartitionDate, step.ClusterUuid, step.Uuid, step.Sequence}...)
						if err != nil {
							return err
						}
						for rows.Next() {
							tmp := new(servicev3.ServiceStep)
							if err := tmp.Scan(rows); err != nil {
								return err
							}

							// get last revision
							if record.Timestamp.Before(tmp.Timestamp) {
								record = tmp
							}
						}

						record.Timestamp = step_status.Timestamp
						// reset update columns
						record.Status = step_status.Status
						record.Started = step_status.Started
						record.Ended = step_status.Ended

						insert := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
							record.TableName(),
							strings.Join(record.ColumnNames(), ", "),
							strings.Join(Repeat(len(record.ColumnNames()), "?"), ", "),
						)

						// fmt.Printf("%v\n", insert)
						// fmt.Printf("%v\n", record.Values())
						if _, err := db.Exec(insert, record.Values()...); err != nil {
							return err
						}

					}
				}
			}
			stmt_insert_service_result := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
				service_result.TableName(),
				strings.Join(service_result.ColumnNames(), ", "),
				strings.Join(Repeat(len(service_result.ColumnNames()), "?"), ", "),
			)

			_, err = db.Exec(stmt_insert_service_result, service_result.Values()...)
			if err != nil {
				return err
			}

			return nil
		}

		return func(count int) {

			var n int
			for i := 0; i < count; i++ {
				n = i
				if err := fn(i); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			fmt.Fprintf(os.Stdout, "client_uuid=%v n=%v\n", client_uuid, n)
		}
	}

	stime := time.Now()
	defer func() {
		etime := time.Since(stime)
		fmt.Fprintln(os.Stdout, etime)
	}()

	workers := 1
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	gencount := (10 * 1000)

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := env.NewSqlDB("sudory_schema_test_v3")
			defer db.Close()

			client_uuid := macro.NewUuidString()

			task(db, client_uuid)(gencount / workers)
		}()
	}

	wg.Wait()

}

func SplitTrimDate(t time.Time) (year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) {
	year, month, day = t.Date()
	loc = t.Location()
	return
}

func TrimDate(t time.Time) time.Time {
	return time.Date(SplitTrimDate(t))
}

// func init() {
// 	enigmaConfigFilename := "enigma.yml"
// 	if err := newEnigma(enigmaConfigFilename); err != nil {
// 		panic(err)
// 	}
// }

// func newEnigma(configFilename string) error {
// 	config := enigma.Config{}
// 	if err := configor.Load(&config, configFilename); err != nil {
// 		return errors.Wrapf(err, "read enigma config file %v",
// 			logs.KVL(
// 				"filename", configFilename,
// 			))
// 	}

// 	if err := enigma.LoadConfig(config); err != nil {
// 		b, _ := ioutil.ReadFile(configFilename)

// 		return errors.Wrapf(err, "load enigma config %v",
// 			logs.KVL(
// 				"filename", configFilename,
// 				"config", string(b),
// 			))
// 	}

// 	if len(config.CryptoAlgorithmSet) == 0 {
// 		return errors.New("'enigma cripto alg set' is empty")
// 	}

// 	for _, k := range dctv1.CiperKeyNames() {
// 		if _, ok := config.CryptoAlgorithmSet[k]; !ok {
// 			return errors.Errorf("not found enigma machine name%s",
// 				logs.KVL(
// 					"key", k,
// 				))
// 		}
// 	}

// 	enigma.PrintConfig(os.Stdout, config)

// 	for key := range config.CryptoAlgorithmSet {
// 		const quickbrownfox = `the quick brown fox jumps over the lazy dog`
// 		encripted, err := enigma.CipherSet(key).Encode([]byte(quickbrownfox))
// 		if err != nil {
// 			return errors.Wrapf(err, "enigma test: encode %v",
// 				logs.KVL("config-name", key))
// 		}
// 		plain, err := enigma.CipherSet(key).Decode(encripted)
// 		if err != nil {
// 			return errors.Wrapf(err, "enigma test: decode %v",
// 				logs.KVL("config-name", key))
// 		}

// 		if strings.Compare(quickbrownfox, string(plain)) != 0 {
// 			return errors.Errorf("enigma test: diff result %v",
// 				logs.KVL("config-name", key))
// 		}
// 	}

// 	return nil
// }

func Repeat(n int, s string) []string {
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = s
	}
	return ss
}
