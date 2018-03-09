package vst2

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	wav "github.com/youpy/go-wav"
)

const (
	wavPath = "_testdata/test.wav"
)

var (
	samples64 [][]float64 //to test processDoubleReplacing
	samples32 [][]float32 //to test processReplacing

	wavFormat *wav.WavFormat

	outputFile string
)

//Read wav for test
func init() {
	out := flag.String("out", "", "Output file for processed audio")
	flag.Parse()
	outputFile = *out

	var wavSamples []wav.Sample
	inFile, _ := os.Open(wavPath)
	defer inFile.Close()
	reader := wav.NewReader(inFile)
	wavFormat, _ = reader.Format()

	for {
		read, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		wavSamples = append(wavSamples, read...)
	}

	samples64 = convertWavSamplesToFloat64(wavSamples)
	samples32 = convertWavSamplesToFloat32(wavSamples)
}

//Test load plugin
func TestPlugin(t *testing.T) {
	library, err := Open(pluginPath)
	if err != nil {
		t.Fatalf("Failed to open library: %v\n", err)
	}
	defer library.Close()
	assert.NotNil(t, library.entryPoint)
	assert.NotNil(t, library.library)
	assert.NotNil(t, library.Name)
	assert.NotNil(t, library.Path)

	plugin, err := library.Open()
	if err != nil {
		t.Fatalf("Failed to open plugin: %v\n", err)
	}
	defer plugin.Close()
	assert.Equal(t, len(plugins), 1)
	assert.NotNil(t, plugin.effect)
	assert.NotNil(t, plugin.Name)
	assert.NotNil(t, plugin.Path)
	assert.Equal(t, true, plugin.CanProcessFloat32())
	assert.Equal(t, false, plugin.CanProcessFloat64())

	plugin.Dispatch(EffOpen, 0, 0, nil, 0.0)

	// Set default sample rate and block size
	sampleRate := 44100.0
	blocksize := int64(512)
	plugin.Dispatch(EffSetSampleRate, 0, 0, nil, sampleRate)
	plugin.Dispatch(EffSetBlockSize, 0, blocksize, nil, 0.0)
	plugin.Dispatch(EffMainsChanged, 0, 1, nil, 0.0)

	processedSamples := plugin.ProcessFloat32(samples32)
	// processedSamples := plugin.ProcessFloat64(samples64)
	// processedSamples := make([][]float64, len(samples64))
	// for i := range processedSamples {
	// 	processedSamples[i] = make([]float64, 0, len(samples64[0]))
	// }
	// c := float64AsChan(samples64, int(blocksize))
	// for samples := range c {
	// 	processedChunk := plugin.Process(samples)
	// 	for i := range processedChunk {
	// 		processedSamples[i] = append(processedSamples[i], processedChunk[i]...)
	// 	}
	// }

	//processedSamples := plaugin.Process(samples64)

	assert.NotNil(t, processedSamples)
	assert.NotEmpty(t, processedSamples)
	assert.Equal(t, 2, len(processedSamples))
	assert.NotEqual(t, 0.0, processedSamples[0][0])
	assert.NotEqual(t, 0.0, processedSamples[0][1])

	if processedSamples == nil {
		return
	}
	for i := 0; i < 20; i++ {
		t.Logf("Sample %d: [%.6f][%.6f]\n", i, processedSamples[0][i], processedSamples[1][i])
	}

	if len(outputFile) > 0 {
		saveSamples32(t, processedSamples)
	}
}

// save samples to temp file
func saveSamples64(t *testing.T, processedSamples [][]float64) {
	outFile, err := ioutil.TempFile("./", outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	numChannels := uint16(len(processedSamples))
	numSamples := uint32(len(processedSamples[0]))
	writer := wav.NewWriter(outFile, numSamples, numChannels, wavFormat.SampleRate, wavFormat.BitsPerSample)
	writer.WriteSamples(convertFloat64ToWavSamples(processedSamples))
}

// save samples to temp file
func saveSamples32(t *testing.T, processedSamples [][]float32) {
	outFile, err := ioutil.TempFile("./", outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	numChannels := uint16(len(processedSamples))
	numSamples := uint32(len(processedSamples[0]))
	writer := wav.NewWriter(outFile, numSamples, numChannels, wavFormat.SampleRate, wavFormat.BitsPerSample)
	writer.WriteSamples(convertFloat32ToWavSamples(processedSamples))
}

//convert WAV samples to float64 slice
func convertWavSamplesToFloat64(wavSamples []wav.Sample) (samples [][]float64) {

	samples = make([][]float64, 2)

	samples[0] = make([]float64, 0, len(wavSamples))
	samples[1] = make([]float64, 0, len(wavSamples))

	for _, sample := range wavSamples {
		samples[0] = append(samples[0], float64(sample.Values[0])/0x8000)
		samples[1] = append(samples[1], float64(sample.Values[1])/0x8000)
	}
	return samples
}

//convert WAV samples to float32 slice
func convertWavSamplesToFloat32(wavSamples []wav.Sample) (samples [][]float32) {

	samples = make([][]float32, 2)

	samples[0] = make([]float32, 0, len(wavSamples))
	samples[1] = make([]float32, 0, len(wavSamples))

	for _, sample := range wavSamples {
		//if i < 512 {
		samples[0] = append(samples[0], float32(sample.Values[0])/0x8000)
		samples[1] = append(samples[1], float32(sample.Values[1])/0x8000)
		//}
	}
	return samples
}

func convertFloat64ToWavSamples(samples [][]float64) (wavSamples []wav.Sample) {
	wavSamples = make([]wav.Sample, len(samples[0]))
	for i := 0; i < len(samples[0]); i++ {
		wavSamples[i].Values[0] = int(samples[0][i] * 0x7FFF)
		wavSamples[i].Values[1] = int(samples[1][i] * 0x7FFF)
	}
	return
}

func convertFloat32ToWavSamples(samples [][]float32) (wavSamples []wav.Sample) {
	wavSamples = make([]wav.Sample, len(samples[0]))
	for i := 0; i < len(samples[0]); i++ {
		wavSamples[i].Values[0] = int(samples[0][i] * 0x7FFF)
		wavSamples[i].Values[1] = int(samples[1][i] * 0x7FFF)
	}
	return
}

func float64AsChan(samples [][]float64, blocksize int) chan [][]float64 {
	c := make(chan [][]float64)
	numChannels := len(samples)
	start, end := 0, blocksize
	go func() {
		defer close(c)
		for end < len(samples[0]) {
			s := make([][]float64, numChannels)
			for i := range samples {
				s[i] = samples[i][start:end]
			}
			end += blocksize
			c <- s
		}
	}()
	return c
}
