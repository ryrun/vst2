// +build plugin

package vst2

//#include "include/vst.h"
import "C"
import (
	"unsafe"
)

//export newGoPlugin
// instantiate go plugin
func newGoPlugin(cp *C.CPlugin, c C.HostCallback) {
	p := PluginAllocator(HostCallback{c})
	cp.magic = C.int(EffectMagic)
	cp.numInputs = C.int(p.InputChannels)
	cp.numOutputs = C.int(p.OutputChannels)
	cp.flags = cp.flags | C.int(PluginDoubleProcessing)
	// cp.flags = cp.flags | C.int(PluginFloatProcessing)
	plugins.Lock()
	p.inputDouble = DoubleBuffer{
		numChannels: p.InputChannels,
		data:        make([]*C.double, p.InputChannels),
	}
	p.outputDouble = DoubleBuffer{
		numChannels: p.OutputChannels,
		data:        make([]*C.double, p.OutputChannels),
	}
	plugins.mapping[unsafe.Pointer(cp)] = &p
	plugins.Unlock()
}

//export dispatchPluginBridge
// global dispatch, calls real plugin dispatch.
func dispatchPluginBridge(cp *C.CPlugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	p, ok := getPlugin(cp)
	if !ok {
		return 0
	}
	return p.DispatchFunc(PluginOpcode(opcode), index, value, ptr, opt)
}

//export processDoublePluginBridge
// global processDouble, calls real plugin processDouble.
func processDoublePluginBridge(cp *C.CPlugin, in, out **C.double, sampleFrames int32) {
	p, ok := getPlugin(cp)
	if !ok {
		return
	}
	for i := range p.inputDouble.data {
		p.inputDouble.data[i] = getDoubleChannel(in, i)
	}
	for i := range p.outputDouble.data {
		p.outputDouble.data[i] = getDoubleChannel(out, i)
	}
	p.inputDouble.size = int(sampleFrames)
	p.outputDouble.size = int(sampleFrames)
	p.ProcessDoubleFunc(p.inputDouble, p.outputDouble)
	return
}

//export processFloatPluginBridge
// global processFloat, calls real plugin processFloat.
func processFloatPluginBridge(cp *C.CPlugin, in, out **float32, sampleFrames int32) {
	return
}

//export getParameterPluginBridge
// global getParameter, calls real plugin getParameter.
func getParameterPluginBridge(cp *C.CPlugin, index int32) float32 {
	return 0
}

//export setParameterPluginBridge
// global setParameter, calls real plugin setParameter.
func setParameterPluginBridge(cp *C.CPlugin, index int32, value float32) {
	return
}

func getDoubleChannel(buf **C.double, i int) *C.double {
	ptrPtr := (**C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + uintptr(i)*unsafe.Sizeof(*buf)))
	return *ptrPtr
}
