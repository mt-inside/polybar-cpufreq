package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TODO:
// * scaling_cur_freq
// * scaling_max_freq / /sys/devices/system/cpu/cpu0/cpufreq/bios_limit
// * bold most common p-state, grey any at 0 instances

/*
* See:
* - turbostat
* - cpupower
* - is turbo on: /sys/devices/system/cpu/cpufreq/boost
* - all p-states: /sys/devices/system/cpu/cpufreq/policy0/scaling_available_frequencies
 */

func getFreqs() []string {
	// Note: only includes what I think are p-states; there's more frequencies in the middle that can be adopted, plus these don't include boost frequencies
	buf, err := os.ReadFile("/sys/devices/system/cpu/cpufreq/policy0/scaling_available_frequencies")
	if err != nil {
		panic(err)
	}

	return strings.Split(string(buf), " ")
}

func normalise(x float64) string {
	return strconv.FormatFloat(x/1000000, 'f', 1, 64)
}

var blocks = []string{"▔" /* space is no good, cause it's a variable-width font */, "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

func main() {
	cpus, err := filepath.Glob("/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq")
	if err != nil {
		panic(err)
	}

	freqs := map[string]int{}

	//  TODO: highlight the op freqs if they're "available" ie a p-state
	// availFreqs := getFreqs()
	// for _, availFreqStr := range availFreqs {
	// 	availFreq, err := strconv.ParseFloat(strings.TrimSpace(availFreqStr), 64)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	bucket := normalise(availFreq)
	// 	freqs[bucket] = 0
	// }

	for {
		for k := range freqs {
			freqs[k] = 0
		}

		for _, cpu := range cpus {
			buf, err := os.ReadFile(cpu)
			if err != nil {
				panic(err)
			}
			freq, err := strconv.ParseFloat(strings.TrimSpace(string(buf)), 64)
			if err != nil {
				continue
			}
			bucket := normalise(freq)
			freqs[bucket] += 1
		}

		keys := make([]string, 0, len(freqs))
		for key := range freqs {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		render := ""
		fastest := ""
		for i := len(keys) - 1; i >= 0; i-- {
			if freqs[keys[i]] != 0 && fastest == "" {
				fastest = keys[i]
			}
			// Pick a non-zero block, unless the number is actually 0
			blk := float64(freqs[keys[i]])/float64(len(cpus))*float64((len(blocks)-2)) + 1
			if freqs[keys[i]] == 0 {
				blk = 0
			}
			render += blocks[int32(blk)]
		}
		fmt.Println(fastest + " " + render)

		time.Sleep(1 * time.Second)
	}
}
