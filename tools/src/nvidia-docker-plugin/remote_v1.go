// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"

	"nvidia"
)

type remoteV10 struct{}

func (r *remoteV10) version() string { return "v1.0" }

func (r *remoteV10) gpuInfo(resp http.ResponseWriter, req *http.Request) {
	const tpl = `
	Driver version:  	{{driverVersion}}
	Supported CUDA version:  	{{cudaVersion}}
	{{range $i, $e := .}}
	Device #{{$i}}
	  Model:  	{{.Model}}
	  UUID:  	{{.UUID}}
	  Path:  	{{.Path}}
	  Family: 	{{.Family}}
	  Arch:  	{{.Arch}}
	  Cores:  	{{.Cores}}
	  Power:  	{{.Power}} W
	  CPU Affinity:  	NUMA node{{.CPUAffinity}}
	  PCI
	    Bus ID:  	{{.PCI.BusID}}
	    BAR1:  	{{.PCI.BAR1}} MiB
	    Bandwidth:  	{{.PCI.Bandwidth}} GB/s
	  Memory
	    ECC:  	{{.Memory.ECC}}
	    Global:  	{{.Memory.Global}} MiB
	    Constant:  	{{.Memory.Constant}} KiB
	    L1 / Shared:  	{{.Memory.Shared}} KiB
	    L2 Cache:  	{{.Memory.L2Cache}} KiB
	    Bandwidth:  	{{.Memory.Bandwidth}} GB/s
	  Clocks
	    Core:  	{{.Clocks.Core}} MHz
	    Memory:  	{{.Clocks.Memory}} MHz
	  P2P Available{{if len .Topology | eq 0}}:  	None{{else}}{{range .Topology}}
	    {{.BusID}} - {{(.Link.String)}}{{end}}{{end}}
	{{end}}
	`
	m := template.FuncMap{
		"driverVersion": nvidia.GetDriverVersion,
		"cudaVersion":   nvidia.GetCUDAVersion,
	}
	t := template.Must(template.New("").Funcs(m).Parse(tpl))
	w := tabwriter.NewWriter(resp, 0, 4, 0, ' ', 0)

	assert(t.Execute(w, Devices))
	assert(w.Flush())
}

func (r *remoteV10) gpuInfoJSON(resp http.ResponseWriter, req *http.Request) {
	var body bytes.Buffer

	writeInfoJSON(&body)
	resp.Header().Set("Content-Type", "application/json")
	_, err := body.WriteTo(resp)
	assert(err)
}

func writeInfoJSON(wr io.Writer) {
	var err error

	r := struct {
		Version struct{ Driver, CUDA string }
		Devices []nvidia.Device
	}{
		Devices: Devices,
	}
	r.Version.Driver, err = nvidia.GetDriverVersion()
	assert(err)
	r.Version.CUDA, err = nvidia.GetCUDAVersion()
	assert(err)

	assert(json.NewEncoder(wr).Encode(r))
}

func (r *remoteV10) gpuStatus(resp http.ResponseWriter, req *http.Request) {
	const tpl = `{{range $i, $e := .}}{{$s := (.Status)}}
	Device #{{$i}}
	  Power:  	{{$s.Power}} / {{.Power}} W
	  Temperature:  	{{$s.Temperature}} °C
	  Utilization
	    GPU:  	{{$s.Utilization.GPU}} %
	    Encoder:  	{{$s.Utilization.Encoder}} %
	    Decoder:  	{{$s.Utilization.Decoder}} %
	  Memory
	    Global:  	{{$s.Memory.GlobalUsed}} / {{.Memory.Global}} MiB
	    ECC Errors{{if not $s.Memory.ECCErrors}}:  	N/A{{else}}
	      L1 Cache:  	{{$s.Memory.ECCErrors.L1Cache}}
	      L2 Cache:  	{{$s.Memory.ECCErrors.L2Cache}}
	      Global:  	{{$s.Memory.ECCErrors.Global}}{{end}}
	  PCI
	    BAR1:  	{{$s.PCI.BAR1Used}} / {{.PCI.BAR1}} MiB
	    Throughput{{if not $s.PCI.Throughput}}:  	N/A{{else}}
	      RX:  	{{$s.PCI.Throughput.RX}} KB/s
	      TX:  	{{$s.PCI.Throughput.TX}} KB/s{{end}}
	  Clocks
	    Core:  	{{$s.Clocks.Core}} MHz
	    Memory:  	{{$s.Clocks.Memory}} MHz
	  Processes{{if len $s.Processes | eq 0}}:  	None{{else}}{{range $s.Processes}}
	    {{.PID}} - {{.Name}}{{end}}{{end}}
	{{end}}
	`
	t := template.Must(template.New("").Parse(tpl))
	w := tabwriter.NewWriter(resp, 0, 4, 0, ' ', 0)

	assert(t.Execute(w, Devices))
	assert(w.Flush())
}

func (r *remoteV10) gpuStatusJSON(resp http.ResponseWriter, req *http.Request) {
	var body bytes.Buffer

	writeStatusJSON(&body)
	resp.Header().Set("Content-Type", "application/json")
	_, err := body.WriteTo(resp)
	assert(err)
}

func writeStatusJSON(wr io.Writer) {
	status := make([]*nvidia.DeviceStatus, 0, len(Devices))

	for i := range Devices {
		s, err := Devices[i].Status()
		assert(err)
		status = append(status, s)
	}
	r := struct{ Devices []*nvidia.DeviceStatus }{status}
	assert(json.NewEncoder(wr).Encode(r))
}

func (r *remoteV10) dockerCLI(resp http.ResponseWriter, req *http.Request) {
	var body bytes.Buffer

	devs := strings.Split(req.FormValue("dev"), " ")
	vols := strings.Split(req.FormValue("vol"), " ")

	if err := dockerCLIDevices(&body, devs); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	body.WriteRune(' ')
	if err := dockerCLIVolumes(&body, vols); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := body.WriteTo(resp)
	assert(err)
}

func dockerCLIDevices(wr io.Writer, ids []string) error {
	const tpl = "--device=/dev/nvidiactl --device=/dev/nvidia-uvm{{range .}} --device={{.}}{{end}}"

	devs := make([]string, 0, len(Devices))

	if len(ids) == 1 && (ids[0] == "*" || ids[0] == "") {
		for i := range Devices {
			devs = append(devs, Devices[i].Path)
		}
	} else {
		for _, id := range ids {
			i, err := strconv.Atoi(id)
			if err != nil || i < 0 || i >= len(Devices) {
				return fmt.Errorf("invalid device: %s", id)
			}
			devs = append(devs, Devices[i].Path)
		}
	}
	t := template.Must(template.New("").Parse(tpl))
	assert(t.Execute(wr, devs))
	return nil
}

func dockerCLIVolumes(wr io.Writer, names []string) error {
	const tpl = "--volume-driver=nvidia{{range .}} --volume={{.}}{{end}}"

	vols := make([]string, 0, len(Volumes))

	if len(names) == 1 && (names[0] == "*" || names[0] == "") {
		for _, v := range Volumes {
			vols = append(vols, fmt.Sprintf("%s:%s:ro", v.Name, v.Mountpoint))
		}
	} else {
		for _, n := range names {
			v, ok := Volumes[n]
			if !ok {
				return fmt.Errorf("invalid volume: %s", n)
			}
			vols = append(vols, fmt.Sprintf("%s:%s:ro", v.Name, v.Mountpoint))
		}
	}
	t := template.Must(template.New("").Parse(tpl))
	assert(t.Execute(wr, vols))
	return nil
}

func (r *remoteV10) mesosCLI(resp http.ResponseWriter, req *http.Request) {
	const format = "--attributes=gpus:%s --resources=gpus:{%s}"

	// Generate Mesos attributes
	var b bytes.Buffer
	writeInfoJSON(&b)
	attr := base64Encode(zlibCompress(b.Bytes()))

	// Generate Mesos custom resources
	uuids := make([]string, 0, len(Devices))
	for i := range Devices {
		uuids = append(uuids, Devices[i].UUID)
	}
	res := strings.Join(uuids, ",")

	_, err := fmt.Fprintf(resp, format, attr, res)
	assert(err)
}

func zlibCompress(buf []byte) []byte {
	b := bytes.NewBuffer(make([]byte, 0, len(buf)))
	w := zlib.NewWriter(b)
	_, err := w.Write(buf)
	assert(err)
	err = w.Close()
	assert(err)
	return b.Bytes()
}

func base64Encode(buf []byte) string {
	s := base64.URLEncoding.EncodeToString(buf)
	if n := len(buf) % 3; n > 0 {
		s = s[:len(s)-(3-n)] // remove padding (RFC 6920)
	}
	return s
}
