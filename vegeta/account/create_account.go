package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Optimization 1: sync.Pool for byte buffers (*bytes.Buffer).
// Avoids memory allocations for each targeter request, reducing pressure on the Garbage Collector (GC).
var bufferPool = sync.Pool{
	New: func() interface{} {
		// Allocates a pre-sized buffer to avoid reallocations of the internal slice
		return bytes.NewBuffer(make([]byte, 0, 256))
	},
}

// Optimization 2: sync.Pool for rand.Rand instances.
// The global rand uses internal mutexes that lock under high concurrency. The Pool isolates instances per goroutine/use.
var rngPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

func main() {
	// Flags for command line control
	rateFlag := flag.Int("rate", 5000, "Request rate per second")
	durationFlag := flag.Duration("duration", 60*time.Second, "Test duration (e.g., 30s, 1m)")
	urlFlag := flag.String("url", "http://localhost:8080/accounts", "Target URL for the POST")
	connectionsFlag := flag.Int("connections", 10000, "Maximum number of simultaneous connections")
	workersFlag := flag.Int("workers", 500, "Number of parallel workers")
	flag.Parse()

	rate := vegeta.Rate{Freq: *rateFlag, Per: time.Second}
	duration := *durationFlag
	targetURL := *urlFlag

	accountingTypes := []string{"ASSET", "LIABILITY", "EQUITY", "REVENUE", "EXPENSE"}
	currencies := []string{"BRL", "USD", "EUR"}

	targeter := func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNoTargets
		}

		// Retrieves the random generator from the pool
		rng := rngPool.Get().(*rand.Rand)
		ownerID := uuid.New().String()
		accountExtID := uuid.New().String()

		// Optimization 3: Use of strconv.FormatInt instead of fmt.Sprintf to convert timestamp
		accountNum := strconv.FormatInt(time.Now().UnixNano(), 10)
		taxId := GerarCNPJ(false, rng)

		randomAccountingType := accountingTypes[rng.Intn(len(accountingTypes))]
		randomCurrency := currencies[rng.Intn(len(currencies))]

		// Returns the rand to the pool for reuse
		rngPool.Put(rng)

		tgt.Method = http.MethodPost
		tgt.URL = targetURL
		tgt.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}

		// Optimization 4: Use of bufferPool to assemble the JSON without allocating new strings via heavy fmt.Sprintf
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()

		buf.WriteString(`{"owner_id": "`)
		buf.WriteString(ownerID)
		buf.WriteString(`", "account_external_id": "`)
		buf.WriteString(accountExtID)
		buf.WriteString(`", "tax_id": "`)
		buf.WriteString(taxId)
		buf.WriteString(`", "account_number": "`)
		buf.WriteString(accountNum)
		buf.WriteString(`", "account_type": "`)
		buf.WriteString(randomAccountingType)
		buf.WriteString(`", "currency": "`)
		buf.WriteString(randomCurrency)
		buf.WriteString(`"}`)

		// Assigns the bytes from the buffer directly to the target's body
		tgt.Body = append([]byte(nil), buf.Bytes()...)

		// Returns the buffer to the pool
		bufferPool.Put(buf)

		return nil
	}

	// Optimization 5: Tuned transport for high performance (aggressive Keep-Alive, no unnecessary compression)
	tr := &http.Transport{
		MaxIdleConns:        50000,
		MaxIdleConnsPerHost: 50000,
		MaxConnsPerHost:     0, // No limit per host beyond global connections
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true, // Disables compression to save CPU on the client
		DisableKeepAlives:   false,
		ForceAttemptHTTP2:   true,
	}
	client := &http.Client{Transport: tr}

	attacker := vegeta.NewAttacker(
		vegeta.Client(client),
		vegeta.Connections(*connectionsFlag),
		vegeta.Workers(uint64(*workersFlag)),
	)
	var metrics vegeta.Metrics

	fmt.Printf("Starting dynamic load test (Optimized) \n"+
		"| URL: %s \n"+
		"| Rate: %d req/s \n"+
		"| Connections: %d \n"+
		"| Workers: %d \n"+
		"| Duration: %v...\n", targetURL, *rateFlag, *connectionsFlag, *workersFlag, duration)

	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}

func GerarCNPJ(comPontos bool, r *rand.Rand) string {
	nums := make([]int, 12)
	for i := 0; i < 8; i++ {
		nums[i] = r.Intn(10)
	}
	nums[8], nums[9], nums[10], nums[11] = 0, 0, 0, 1

	pesos1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	dv1 := calculaDV(nums, pesos1)

	numsComDv1 := append(nums, dv1)

	pesos2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	dv2 := calculaDV(numsComDv1, pesos2)

	cnpjNums := append(numsComDv1, dv2)

	if comPontos {
		return fmt.Sprintf("%d%d%d.%d%d%d.%d%d%d/%d%d%d-%d%d",
			cnpjNums[0], cnpjNums[1], cnpjNums[2],
			cnpjNums[3], cnpjNums[4], cnpjNums[5],
			cnpjNums[6], cnpjNums[7], cnpjNums[8],
			cnpjNums[9], cnpjNums[10], cnpjNums[11],
			cnpjNums[12], cnpjNums[13])
	}

	// Optimization 6: Direct string assembly without fmt.Sprintf for the numeric CNPJ
	buf := make([]byte, 14)
	for i, n := range cnpjNums {
		buf[i] = byte('0' + n)
	}
	return string(buf)
}

func calculaDV(numeros []int, pesos []int) int {
	soma := 0
	for i, v := range numeros {
		soma += v * pesos[i]
	}
	resto := soma % 11
	if resto < 2 {
		return 0
	}
	return 11 - resto
}
