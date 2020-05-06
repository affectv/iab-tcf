package iab_tcf

import (
	"github.com/montanaflynn/stats"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type data struct {
	message   string
	consent string
	result  string
}

var _ = Describe("performance", func() {

	const (
		times = 1000
	)

	var (
		testResults = map[string][]int64{}
		testData    = []data{
			data{
				message: "v1 consent",
				consent: "BOlLbqtOlLbqtAVABADECg-AAAApp7v______9______9uz_Ov_v_f__33e8__9v_l_7_-___u_-3zd4u_1vf99yfm1-7etr3tp_87ues2_Xur__79__3z3_9phP78k89r7337Ew-v02", 
				result: "111101110111111111111111111111111111111111111111111110111111111111111111111111111111111111111110110111011001111111100111010111111111110111111111101111111111111111111011111011101111011110011111111111111110110111111111110010111111111101111111111111011111111111111111110111011111111111011011111001101110111100010111011111111010110111101111111110111110111001001111110011011010111111011101101111010110110101111011110110110100111111111110011101110111001111010110011011011111101011110111010101111111111111111101111110111111111111111011111001111011111111111110110100110000100111111101111110010010011110011110110101111101111011111011111101100010011000011111010111111010011011",
			},
			data{
				message: "v2 consent",
				consent: "COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA", 
				result: "010001011111001001001010100000011001100011011010110110000100100011011001100001010100010111000100000100010000100001100011000001111100100111000111001010001010000000110010000000001010000110000001000000100010110001001011000000110101000000100001111000010100101100000110100000011000101110001000000000000000011100000100011001100000100100000001000000000000011000000100000000010000010000000000001000000100000001000000100010010000011100011000000100110000001001000000000000000010000000000010000000000110001001000100001000110001000000010000011000110000001011001100110100100000000100100000000100000010000010000000010001101101100011010000010100000000001000001001010110011000011110010000011101001010011001100100001001100011101111010011101011100000111011111111111101",
			},
		}
	)

	AfterSuite(func() {
		for testName, results := range testResults {
			p99, _ := stats.Percentile(stats.LoadRawData(results), 99)
			Expect(p99).Should(BeNumerically("<", 2000), testName+" shouldn't take too long.")
		}
	})

	var run = func(consent string, result string) func(b Benchmarker) {
		return func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				output, err := NewConsent(consent)
				Expect(err).NotTo(HaveOccurred())
				Expect(output.GetConsentBitstring()).To(Equal(result))
			})
			testName := CurrentGinkgoTestDescription().TestText
			testResults[testName] = append(testResults[testName], runtime.Microseconds())
			b.RecordValue("spent in microseconds", float64(runtime.Microseconds()))
		}
	}

	Context("custom parser", func() {

		var runCustom = func(consent string, result string) func(b Benchmarker) {
			return run(consent, result)
		}

		for _, data := range testData {
			testResults[data.message] = []int64{}
			Measure(data.message, runCustom(data.consent, data.result), times)
		}
	})
})