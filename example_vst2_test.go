package vst2_test

import (
	"fmt"
	"log"
	"runtime"

	"pipelined.dev/signal"
	"pipelined.dev/vst2"
)

// PrinterHostCallback returns closure that prints received opcode with provided
// prefix. This technique allows to provide callback with any context needed.
func PrinterHostCallback(prefix string) vst2.HostCallbackFunc {
	return func(code vst2.HostOpcode, _ vst2.Index, _ vst2.Value, _ vst2.Ptr, _ vst2.Opt) vst2.Return {
		fmt.Printf("%s: %v\n", prefix, code)
		return 0
	}
}

// sample data that we will process.
var data = [][]float64{
	{
		-0.0027160645, -0.0039978027, -0.0071411133, -0.0065307617, 0.0038757324, 0.021972656, 0.041229248, 0.055511475, 0.064971924, 0.07342529, 0.08300781, 0.092681885, 0.10070801, 0.110809326, 0.12677002, 0.15231323, 0.19058228, 0.24459839, 0.3140869, 0.38861084, 0.44683838, 0.47177124, 0.46643066, 0.45007324, 0.4449768, 0.45724487, 0.47451782, 0.48321533, 0.47824097, 0.46679688, 0.45999146, 0.46765137, 0.491333, 0.52505493, 0.555542, 0.57055664, 0.5701599, 0.56591797, 0.5706787, 0.58740234, 0.60510254, 0.6090698, 0.5979004, 0.5837097, 0.58288574, 0.6016846, 0.63098145, 0.6526184, 0.6595459, 0.6639404, 0.6861267, 0.73825073, 0.8117676, 0.8830261, 0.9341736, 0.9536133, 0.95010376, 0.9467163, 0.9489136, 0.94799805, 0.9470825, 0.9473877, 0.9465332, 0.9461365, 0.9458313, 0.94555664, 0.94525146, 0.94454956, 0.94351196, 0.9433899, 0.9428711, 0.94226074, 0.94226074, 0.9415283, 0.94104004, 0.9407959, 0.93963623, 0.9390869, 0.93881226, 0.9384155, 0.93777466, 0.9367676, 0.93652344, 0.9362488, 0.9352417, 0.9342041, 0.93414307, 0.93356323, 0.93304443, 0.9326782, 0.93182373, 0.93136597, 0.93084717, 0.9299927, 0.92926025, 0.928833, 0.9283447, 0.9275818, 0.9270935, 0.9267273, 0.92578125, 0.92526245, 0.92489624, 0.9238281, 0.9232178, 0.92266846, 0.9219055, 0.92123413, 0.92123413, 0.9200134, 0.9194336, 0.9192505, 0.91848755, 0.9177246, 0.91744995, 0.9169922, 0.9158325, 0.9154968, 0.9150696, 0.91430664, 0.9137268, 0.91326904, 0.91275024, 0.91189575, 0.91137695, 0.9107971, 0.9099121, 0.9095154, 0.9087219, 0.9081421, 0.9077759, 0.9072571, 0.9066162, 0.90618896, 0.9053345, 0.9050293, 0.9039612, 0.9034729, 0.9031677, 0.90237427, 0.9014282, 0.90130615, 0.90078735, 0.89974976, 0.89959717, 0.8988342, 0.8982239, 0.8980713, 0.8967285, 0.89675903, 0.895813, 0.8949585, 0.89016724, 0.8718872, 0.8529358, 0.85061646, 0.8433838, 0.81521606, 0.78323364, 0.7631836, 0.7581787, 0.7541809, 0.7345276, 0.7017517, 0.6737976, 0.66192627, 0.6621704, 0.65826416, 0.6414795, 0.6186218, 0.59976196, 0.58813477, 0.5793457, 0.5661621, 0.54629517, 0.5232239, 0.49957275, 0.47854614, 0.4595642, 0.44082642, 0.41912842, 0.39367676, 0.3659668, 0.34313965, 0.32885742, 0.32141113, 0.31307983, 0.29656982, 0.2741394, 0.2564392, 0.24880981, 0.2480774, 0.24499512, 0.23535156, 0.2230835, 0.21395874, 0.21069336, 0.20672607, 0.19671631, 0.18139648, 0.16766357, 0.16088867, 0.16131592, 0.16290283, 0.15875244, 0.14880371, 0.14047241, 0.1375122, 0.13754272, 0.1312561, 0.11392212, 0.087524414, 0.06283569, 0.044067383, 0.028747559, 0.013336182, -0.002105713, -0.014251709, -0.022949219, -0.033081055, -0.050811768 - 0.07687378, -0.107299805, -0.13864136, -0.16897583, -0.19692993, -0.21844482, -0.23138428, -0.23873901, -0.24871826, -0.26620483, -0.28753662, -0.30618286, -0.3239441, -0.3468628, -0.37780762, -0.40966797, -0.4310913, -0.4385376, -0.43984985, -0.44317627, -0.4505005, -0.45721436, -0.46203613, -0.47052002, -0.49023438, -0.51968384, -0.5514221, -0.57769775, -0.5969238, -0.61032104, -0.6180115, -0.6163635, -0.6074829, -0.59988403, -0.60150146, -0.6123352, -0.62579346, -0.6381836, -0.65078735, -0.66641235, -0.682251, -0.69226074, -0.6951904, -0.697876, -0.70532227, -0.7165222, -0.724823, -0.7284851, -0.73114014, -0.73703, -0.744751, -0.7517395, -0.7571411, -0.7656555, -0.77575684, -0.78308105, -0.783905, -0.7837219, -0.7885742, -0.79852295, -0.8067322, -0.80700684, -0.8046875, -0.80966187, -0.8235779, -0.8360596, -0.836853, -0.82455444, -0.8106384, -0.80633545, -0.8121643, -0.82003784, -0.82110596, -0.8159485, -0.8136902, -0.8216858, -0.8374634, -0.8496399, -0.8527527, -0.8534546, -0.8637085, -0.88290405, -0.8952637, -0.88739014, -0.87106323, -0.8749695, -0.908844, -0.9442749, -0.94314575, -0.9057007, -0.8771057, -0.89553833, -0.9460144, -0.9734802, -0.95080566, -0.9076538, -0.8945923, -0.9199524, -0.94784546, -0.95092773, -0.9447632, -0.9534912, -0.9665222, -0.9493103, -0.90164185, -0.87387085, -0.90740967, -0.9699402, -0.98379517, -0.92019653, -0.84490967, -0.8469238, -0.9286194, -0.99887085, -0.9760437, -0.88500977, -0.82318115, -0.8460083, -0.90844727, -0.9246521, -0.8715515, -0.8070984, -0.7927246, -0.8227844, -0.8427124, -0.8262329, -0.801239, -0.8036194, -0.8232422, -0.8187561, -0.7727051, -0.7202759, -0.707428, -0.73535156, -0.7595215, -0.74105835, -0.69088745, -0.651947, -0.6512146, -0.67611694, -0.6911621, -0.67507935, -0.6367798, -0.60113525, -0.58914185, -0.602478, -0.6260681, -0.6374817, -0.62338257, -0.59176636, -0.5643921, -0.55389404, -0.5558777, -0.5494995, -0.52194214, -0.4753418, -0.42858887, -0.40090942, -0.40319824, -0.4307251 - 0.46392822, -0.4769287, -0.4588318, -0.42556763, -0.40551758, -0.40753174, -0.40896606, -0.3852539, -0.3413086, -0.31314087, -0.32635498, -0.36328125, -0.3829956, -0.3668213, -0.337677, -0.32650757, -0.3361206, -0.34222412, -0.33035278, -0.3109131, -0.29956055, -0.2885437, -0.2569275, -0.20565796, -0.1643982, -0.15802002, -0.17419434, -0.17623901, -0.14611816, -0.10662842, -0.09020996, -0.09765625, -0.09631348, -0.065216064, -0.0184021, 0.008453369, 0.0012207031, -0.017486572, -0.015106201, 0.014923096, 0.051605225, 0.075408936, 0.087524414, 0.10205078, 0.12762451, 0.1609497, 0.19122314, 0.2098999, 0.21704102, 0.21661377, 0.21902466, 0.23583984, 0.26986694, 0.31033325, 0.34545898, 0.37139893, 0.39263916, 0.41351318, 0.43130493, 0.44195557, 0.4508667, 0.46948242, 0.4954834, 0.5108032, 0.5023804, 0.48156738, 0.47988892, 0.51226807, 0.5572815, 0.58102417, 0.5742798, 0.559845, 0.56741333, 0.60110474, 0.64263916, 0.67437744, 0.6965332, 0.7150574, 0.7319641, 0.74438477, 0.75839233, 0.78186035, 0.81207275, 0.8309021, 0.8224182, 0.79211426, 0.7640991, 0.75460815, 0.7586365, 0.7642822, 0.7694092, 0.78253174, 0.8069763, 0.8355713, 0.86026, 0.8809509, 0.9004822, 0.9154663, 0.91726685, 0.90792847, 0.9007263, 0.90896606, 0.9289856, 0.9402771, 0.93685913, 0.93447876, 0.93066406, 0.92163086, 0.9219971, 0.93292236, 0.93548584, 0.9326172, 0.932312, 0.9321594, 0.93078613, 0.9276428, 0.9147949, 0.8964844, 0.8874512, 0.8934021, 0.91033936, 0.92663574, 0.92926025, 0.9249573, 0.92544556, 0.9252014, 0.923645, 0.92471313, 0.9231262, 0.9221802, 0.9222717, 0.9221802, 0.9222717,
	},
	{
		-0.0027160645, -0.0039978027, -0.0071411133, -0.0065307617, 0.0038757324, 0.021972656, 0.041229248, 0.055511475, 0.064971924, 0.07342529, 0.08300781, 0.092681885, 0.10070801, 0.110809326, 0.12677002, 0.15231323, 0.19058228, 0.24459839, 0.3140869, 0.38861084, 0.44683838, 0.47177124, 0.46643066, 0.45007324, 0.4449768, 0.45724487, 0.47451782, 0.48321533, 0.47824097, 0.46679688, 0.45999146, 0.46765137, 0.491333, 0.52505493, 0.555542, 0.57055664, 0.5701599, 0.56591797, 0.5706787, 0.58740234, 0.60510254, 0.6090698, 0.5979004, 0.5837097, 0.58288574, 0.6016846, 0.63098145, 0.6526184, 0.6595459, 0.6639404, 0.6861267, 0.73825073, 0.8117676, 0.8830261, 0.9341736, 0.9536133, 0.95010376, 0.9467163, 0.9489136, 0.94799805, 0.9470825, 0.9473877, 0.9465332, 0.9461365, 0.9458313, 0.94555664, 0.94525146, 0.94454956, 0.94351196, 0.9433899, 0.9428711, 0.94226074, 0.94226074, 0.9415283, 0.94104004, 0.9407959, 0.93963623, 0.9390869, 0.93881226, 0.9384155, 0.93777466, 0.9367676, 0.93652344, 0.9362488, 0.9352417, 0.9342041, 0.93414307, 0.93356323, 0.93304443, 0.9326782, 0.93182373, 0.93136597, 0.93084717, 0.9299927, 0.92926025, 0.928833, 0.9283447, 0.9275818, 0.9270935, 0.9267273, 0.92578125, 0.92526245, 0.92489624, 0.9238281, 0.9232178, 0.92266846, 0.9219055, 0.92123413, 0.92123413, 0.9200134, 0.9194336, 0.9192505, 0.91848755, 0.9177246, 0.91744995, 0.9169922, 0.9158325, 0.9154968, 0.9150696, 0.91430664, 0.9137268, 0.91326904, 0.91275024, 0.91189575, 0.91137695, 0.9107971, 0.9099121, 0.9095154, 0.9087219, 0.9081421, 0.9077759, 0.9072571, 0.9066162, 0.90618896, 0.9053345, 0.9050293, 0.9039612, 0.9034729, 0.9031677, 0.90237427, 0.9014282, 0.90130615, 0.90078735, 0.89974976, 0.89959717, 0.8988342, 0.8982239, 0.8980713, 0.8967285, 0.89675903, 0.895813, 0.8949585, 0.89016724, 0.8718872, 0.8529358, 0.85061646, 0.8433838, 0.81521606, 0.78323364, 0.7631836, 0.7581787, 0.7541809, 0.7345276, 0.7017517, 0.6737976, 0.66192627, 0.6621704, 0.65826416, 0.6414795, 0.6186218, 0.59976196, 0.58813477, 0.5793457, 0.5661621, 0.54629517, 0.5232239, 0.49957275, 0.47854614, 0.4595642, 0.44082642, 0.41912842, 0.39367676, 0.3659668, 0.34313965, 0.32885742, 0.32141113, 0.31307983, 0.29656982, 0.2741394, 0.2564392, 0.24880981, 0.2480774, 0.24499512, 0.23535156, 0.2230835, 0.21395874, 0.21069336, 0.20672607, 0.19671631, 0.18139648, 0.16766357, 0.16088867, 0.16131592, 0.16290283, 0.15875244, 0.14880371, 0.14047241, 0.1375122, 0.13754272, 0.1312561, 0.11392212, 0.087524414, 0.06283569, 0.044067383, 0.028747559, 0.013336182, -0.002105713, -0.014251709, -0.022949219, -0.033081055, -0.050811768 - 0.07687378, -0.107299805, -0.13864136, -0.16897583, -0.19692993, -0.21844482, -0.23138428, -0.23873901, -0.24871826, -0.26620483, -0.28753662, -0.30618286, -0.3239441, -0.3468628, -0.37780762, -0.40966797, -0.4310913, -0.4385376, -0.43984985, -0.44317627, -0.4505005, -0.45721436, -0.46203613, -0.47052002, -0.49023438, -0.51968384, -0.5514221, -0.57769775, -0.5969238, -0.61032104, -0.6180115, -0.6163635, -0.6074829, -0.59988403, -0.60150146, -0.6123352, -0.62579346, -0.6381836, -0.65078735, -0.66641235, -0.682251, -0.69226074, -0.6951904, -0.697876, -0.70532227, -0.7165222, -0.724823, -0.7284851, -0.73114014, -0.73703, -0.744751, -0.7517395, -0.7571411, -0.7656555, -0.77575684, -0.78308105, -0.783905, -0.7837219, -0.7885742, -0.79852295, -0.8067322, -0.80700684, -0.8046875, -0.80966187, -0.8235779, -0.8360596, -0.836853, -0.82455444, -0.8106384, -0.80633545, -0.8121643, -0.82003784, -0.82110596, -0.8159485, -0.8136902, -0.8216858, -0.8374634, -0.8496399, -0.8527527, -0.8534546, -0.8637085, -0.88290405, -0.8952637, -0.88739014, -0.87106323, -0.8749695, -0.908844, -0.9442749, -0.94314575, -0.9057007, -0.8771057, -0.89553833, -0.9460144, -0.9734802, -0.95080566, -0.9076538, -0.8945923, -0.9199524, -0.94784546, -0.95092773, -0.9447632, -0.9534912, -0.9665222, -0.9493103, -0.90164185, -0.87387085, -0.90740967, -0.9699402, -0.98379517, -0.92019653, -0.84490967, -0.8469238, -0.9286194, -0.99887085, -0.9760437, -0.88500977, -0.82318115, -0.8460083, -0.90844727, -0.9246521, -0.8715515, -0.8070984, -0.7927246, -0.8227844, -0.8427124, -0.8262329, -0.801239, -0.8036194, -0.8232422, -0.8187561, -0.7727051, -0.7202759, -0.707428, -0.73535156, -0.7595215, -0.74105835, -0.69088745, -0.651947, -0.6512146, -0.67611694, -0.6911621, -0.67507935, -0.6367798, -0.60113525, -0.58914185, -0.602478, -0.6260681, -0.6374817, -0.62338257, -0.59176636, -0.5643921, -0.55389404, -0.5558777, -0.5494995, -0.52194214, -0.4753418, -0.42858887, -0.40090942, -0.40319824, -0.4307251 - 0.46392822, -0.4769287, -0.4588318, -0.42556763, -0.40551758, -0.40753174, -0.40896606, -0.3852539, -0.3413086, -0.31314087, -0.32635498, -0.36328125, -0.3829956, -0.3668213, -0.337677, -0.32650757, -0.3361206, -0.34222412, -0.33035278, -0.3109131, -0.29956055, -0.2885437, -0.2569275, -0.20565796, -0.1643982, -0.15802002, -0.17419434, -0.17623901, -0.14611816, -0.10662842, -0.09020996, -0.09765625, -0.09631348, -0.065216064, -0.0184021, 0.008453369, 0.0012207031, -0.017486572, -0.015106201, 0.014923096, 0.051605225, 0.075408936, 0.087524414, 0.10205078, 0.12762451, 0.1609497, 0.19122314, 0.2098999, 0.21704102, 0.21661377, 0.21902466, 0.23583984, 0.26986694, 0.31033325, 0.34545898, 0.37139893, 0.39263916, 0.41351318, 0.43130493, 0.44195557, 0.4508667, 0.46948242, 0.4954834, 0.5108032, 0.5023804, 0.48156738, 0.47988892, 0.51226807, 0.5572815, 0.58102417, 0.5742798, 0.559845, 0.56741333, 0.60110474, 0.64263916, 0.67437744, 0.6965332, 0.7150574, 0.7319641, 0.74438477, 0.75839233, 0.78186035, 0.81207275, 0.8309021, 0.8224182, 0.79211426, 0.7640991, 0.75460815, 0.7586365, 0.7642822, 0.7694092, 0.78253174, 0.8069763, 0.8355713, 0.86026, 0.8809509, 0.9004822, 0.9154663, 0.91726685, 0.90792847, 0.9007263, 0.90896606, 0.9289856, 0.9402771, 0.93685913, 0.93447876, 0.93066406, 0.92163086, 0.9219971, 0.93292236, 0.93548584, 0.9326172, 0.932312, 0.9321594, 0.93078613, 0.9276428, 0.9147949, 0.8964844, 0.8874512, 0.8934021, 0.91033936, 0.92663574, 0.92926025, 0.9249573, 0.92544556, 0.9252014, 0.923645, 0.92471313, 0.9231262, 0.9221802, 0.9222717, 0.9221802, 0.9222717,
	},
}

func Example_plugin() {
	buffer := signal.Allocator{
		Channels: len(data),
		Length:   len(data[0]),
		Capacity: len(data[0]),
	}.Float64()

	signal.WriteStripedFloat64(data, buffer)

	// Open VST library. Library contains a reference to
	// OS-specific handle, that needs to be freed with Close.
	vst, err := vst2.Open(pluginPath())
	if err != nil {
		log.Panicf("failed to open VST library: %v", err)
	}
	defer vst.Close()

	// Load VST plugin with example callback.
	plugin := vst.Load(PrinterHostCallback("Received opcode"))
	defer plugin.Close()

	// Set sample rate in Hertz.
	plugin.SetSampleRate(44100)
	// Set channels information.
	plugin.SetSpeakerArrangement(
		&vst2.SpeakerArrangement{
			Type:        vst2.SpeakerArrMono,
			NumChannels: int32(buffer.Channels()),
		},
		&vst2.SpeakerArrangement{
			Type:        vst2.SpeakerArrMono,
			NumChannels: int32(buffer.Channels()),
		},
	)
	// Set buffer size.
	plugin.SetBufferSize(buffer.Length())
	// Start the plugin.
	plugin.Start()

	// To process data with plugin, we need to use VST2 buffers.
	// It's needed because VST SDK was written in C and expected
	// memory layout differs from Golang slices.
	// We need two buffers for input and output.
	in := vst2.NewDoubleBuffer(buffer.Channels(), buffer.Length())
	out := vst2.NewDoubleBuffer(buffer.Channels(), buffer.Length())

	// Fill input with data values.
	in.CopyFrom(buffer)

	// Process data.
	plugin.ProcessDouble(in, out)
	// Copy processed data.
	out.CopyTo(buffer)

	// Output:
	// Received opcode: HostGetCurrentProcessLevel
	// Received opcode: HostGetCurrentProcessLevel
}

// pluginPath returns a path to OS-specific plugin. It will panic if OS is
// not supported.
func pluginPath() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return "_testdata/TAL-Reverb.dll"
	case "darwin":
		return "_testdata\\TAL-Reverb.vst"
	default:
		panic(fmt.Sprintf("unsupported OS: %v", os))
	}
}
