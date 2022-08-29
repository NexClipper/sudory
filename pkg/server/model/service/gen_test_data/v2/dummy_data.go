package main

import (
	"database/sql"
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
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	new_service := func(name string, cluster string, template string, stepcount int) (*servicev2.Service, []servicev2.ServiceStep) {
		uuid := macro.NewUuidString()
		created := time.Now()

		service := servicev2.Service{}
		service.Uuid = uuid
		service.Created = created
		service.Name = name
		service.Summary = *vanilla.NewNullString(fmt.Sprintf("service %v", name))
		service.ClusterUuid = cluster
		service.TemplateUuid = template
		service.StepCount = stepcount

		step := servicev2.ServiceStep{}
		step.Uuid = uuid
		step.Created = created
		step.Sequence = 0
		step.Name = name
		step.Summary = *vanilla.NewNullString(fmt.Sprintf("step %v", name))
		step.Method = "fake.test.method"
		step.Args = dctv2.CryptoObject{"name": "fake-name", "ns": "fake-ns"}
		step.Sequence = 0

		return &service, []servicev2.ServiceStep{step}
	}

	new_service_status_success := func(service servicev2.Service, client_uuid string) ([]servicev2.ServiceStatus, []servicev2.ServiceStepStatus) {
		// uuid := macro.NewUuidString()
		created := time.Now()
		addtime := func() time.Time {
			created = created.Add(time.Duration(rand.Intn(1000)+1) * time.Millisecond)
			return created
		}

		status1 := servicev2.ServiceStatus{}
		status1.Uuid = service.Uuid
		status1.Created = addtime()
		status1.Status = 1
		status1.StepPosition = 1
		status1.AssignedClientUuid = client_uuid

		status2 := servicev2.ServiceStatus{}
		status2.Uuid = service.Uuid
		status2.Created = addtime()
		status2.Status = 2
		status2.StepPosition = 1
		status2.AssignedClientUuid = client_uuid

		status4 := servicev2.ServiceStatus{}
		status4.Uuid = service.Uuid
		status4.Created = addtime()
		status4.Status = 4
		status4.StepPosition = 1
		status4.AssignedClientUuid = client_uuid

		service_status := []servicev2.ServiceStatus{
			status1,
			status2,
			status4,
		}

		step_status1 := servicev2.ServiceStepStatus{}
		step_status1.Uuid = service.Uuid
		step_status1.Created = addtime()
		step_status1.Sequence = 0
		step_status1.Status = 1
		step_status1.Started = *vanilla.NewNullTime(addtime())
		step_status1.Ended = *vanilla.NewNullTime(addtime())

		step_status2 := servicev2.ServiceStepStatus{}
		step_status2.Uuid = service.Uuid
		step_status2.Created = addtime()
		step_status2.Sequence = 0
		step_status2.Status = 2
		step_status2.Started = *vanilla.NewNullTime(addtime())
		step_status2.Ended = *vanilla.NewNullTime(addtime())

		step_status4 := servicev2.ServiceStepStatus{}
		step_status4.Uuid = service.Uuid
		step_status4.Created = addtime()
		step_status4.Sequence = 0
		step_status4.Status = 4
		step_status4.Started = *vanilla.NewNullTime(addtime())
		step_status4.Ended = *vanilla.NewNullTime(addtime())

		service_step_status := []servicev2.ServiceStepStatus{
			step_status1,
			step_status2,
			step_status4,
		}

		return service_status, service_step_status
	}

	new_service_status_none := func(service servicev2.Service, client_uuid string) ([]servicev2.ServiceStatus, []servicev2.ServiceStepStatus) {

		return []servicev2.ServiceStatus{}, []servicev2.ServiceStepStatus{}
	}

	new_service_status_progress := func(service servicev2.Service, client_uuid string) ([]servicev2.ServiceStatus, []servicev2.ServiceStepStatus) {
		// uuid := macro.NewUuidString()
		created := time.Now()
		addtime := func() time.Time {
			created = created.Add(time.Duration(rand.Intn(1000)+1) * time.Millisecond)
			return created
		}

		status1 := servicev2.ServiceStatus{}
		status1.Uuid = service.Uuid
		status1.Created = addtime()
		status1.Status = 1
		status1.StepPosition = 1
		status1.AssignedClientUuid = client_uuid

		status2 := servicev2.ServiceStatus{}
		status2.Uuid = service.Uuid
		status2.Created = addtime()
		status2.Status = 2
		status2.StepPosition = 1
		status2.AssignedClientUuid = client_uuid

		service_status := []servicev2.ServiceStatus{
			status1,
			status2,
		}

		step_status1 := servicev2.ServiceStepStatus{}
		step_status1.Uuid = service.Uuid
		step_status1.Created = addtime()
		step_status1.Sequence = 0
		step_status1.Status = 1
		step_status1.Started = *vanilla.NewNullTime(addtime())
		step_status1.Ended = *vanilla.NewNullTime(addtime())

		step_status2 := servicev2.ServiceStepStatus{}
		step_status2.Uuid = service.Uuid
		step_status2.Created = addtime()
		step_status2.Sequence = 0
		step_status2.Status = 2
		step_status2.Started = *vanilla.NewNullTime(addtime())
		step_status2.Ended = *vanilla.NewNullTime(addtime())

		service_step_status := []servicev2.ServiceStepStatus{
			step_status1,
			step_status2,
		}

		return service_status, service_step_status
	}

	new_service_status_error := func(service servicev2.Service, client_uuid string) ([]servicev2.ServiceStatus, []servicev2.ServiceStepStatus) {
		// uuid := macro.NewUuidString()
		created := time.Now()
		addtime := func() time.Time {
			created = created.Add(time.Duration(rand.Intn(1000)+1) * time.Millisecond)
			return created
		}

		status1 := servicev2.ServiceStatus{}
		status1.Uuid = service.Uuid
		status1.Created = addtime()
		status1.Status = 1
		status1.StepPosition = 1
		status1.AssignedClientUuid = client_uuid

		status2 := servicev2.ServiceStatus{}
		status2.Uuid = service.Uuid
		status2.Created = addtime()
		status2.Status = 2
		status2.StepPosition = 1
		status2.AssignedClientUuid = client_uuid

		status8 := servicev2.ServiceStatus{}
		status8.Uuid = service.Uuid
		status8.Created = addtime()
		status8.Status = 8
		status8.StepPosition = 1
		status8.AssignedClientUuid = client_uuid
		status8.Message = *vanilla.NewNullString(fmt.Sprintf("error: created=%v", created))

		service_status := []servicev2.ServiceStatus{
			status1,
			status2,
			status8,
		}

		step_status1 := servicev2.ServiceStepStatus{}
		step_status1.Uuid = service.Uuid
		step_status1.Created = addtime()
		step_status1.Sequence = 0
		step_status1.Status = 1
		step_status1.Started = *vanilla.NewNullTime(addtime())
		step_status1.Ended = *vanilla.NewNullTime(addtime())

		step_status2 := servicev2.ServiceStepStatus{}
		step_status2.Uuid = service.Uuid
		step_status2.Created = addtime()
		step_status2.Sequence = 0
		step_status2.Status = 2
		step_status2.Started = *vanilla.NewNullTime(addtime())
		step_status2.Ended = *vanilla.NewNullTime(addtime())

		step_status8 := servicev2.ServiceStepStatus{}
		step_status8.Uuid = service.Uuid
		step_status8.Created = addtime()
		step_status8.Sequence = 0
		step_status8.Status = 8
		step_status8.Started = *vanilla.NewNullTime(addtime())
		step_status8.Ended = *vanilla.NewNullTime(addtime())

		service_step_status := []servicev2.ServiceStepStatus{
			step_status1,
			step_status2,
			step_status8,
		}

		return service_status, service_step_status
	}

	new_service_result := func(service servicev2.Service) *servicev2.ServiceResult {
		created := time.Now()

		r := rand.Int63()

		service_result := &servicev2.ServiceResult{}
		service_result.Uuid = service.Uuid
		service_result.Created = created.Add(time.Duration(rand.Intn(1000)) * time.Millisecond)
		service_result.ResultSaveType = servicev2.ResultSaveTypeDatabase
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

			var service_status []servicev2.ServiceStatus
			var step_status []servicev2.ServiceStepStatus
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

			for _, step := range step {
				stmt_insert_step := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
					step.TableName(),
					strings.Join(step.ColumnNames(), ", "),
					strings.Join(Repeat(len(step.ColumnNames()), "?"), ", "),
				)

				_, err := db.Exec(stmt_insert_step, step.Values()...)
				if err != nil {
					return err
				}
			}

			if true {
				for _, service_status := range service_status {
					stmt_insert_service_status := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
						service_status.TableName(),
						strings.Join(service_status.ColumnNames(), ", "),
						strings.Join(Repeat(len(service_status.ColumnNames()), "?"), ", "),
					)

					_, err := db.Exec(stmt_insert_service_status, service_status.Values()...)
					if err != nil {

						return err
					}
				}

				for _, step_status := range step_status {
					stmt_insert_step_status := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
						step_status.TableName(),
						strings.Join(step_status.ColumnNames(), ", "),
						strings.Join(Repeat(len(step_status.ColumnNames()), "?"), ", "),
					)

					_, err := db.Exec(stmt_insert_step_status, step_status.Values()...)
					if err != nil {

						return err
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
			db := env.NewSqlDB("sudory_schema_test_v2")
			defer db.Close()

			client_uuid := macro.NewUuidString()

			task(db, client_uuid)(gencount / workers)
		}()
	}

	wg.Wait()

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
