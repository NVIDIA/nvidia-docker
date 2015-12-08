// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"bytes"
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
	  Name:  	{{.Name}}
	  UUID:  	{{.UUID}}
	  Path:  	{{.Path}}
	  Gen: 	{{.Gen}}
	  Arch:  	{{.Arch}}
	  Cores:  	{{.Cores}}
	  Power:  	{{.Power}} W
	  CPU Affinity:  	NUMA node{{.CPUAffinity}}
	  PCI
	    BusID:  	{{.PCI.BusID}}
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
	    SM:  	{{.Clocks.SM}} MHz
	    Memory:  	{{.Clocks.Memory}} MHz
	    Graphics:  	{{.Clocks.Graphics}} MHz
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
	if err := writeInfoJSON(&body); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(body.Bytes())
}

func writeInfoJSON(wr io.Writer) error {
	driverVersion, err := nvidia.GetDriverVersion()
	if err != nil {
		return err
	}

	cudaVersion, err := nvidia.GetCUDAVersion()
	if err != nil {
		return err
	}

	r := struct {
		DriverVersion string
		CUDAVersion   string
		Devices       []nvidia.Device
	}{
		driverVersion,
		cudaVersion,
		Devices,
	}
	assert(json.NewEncoder(wr).Encode(r))
	return nil
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
	      L1 Cache:  	{{$s.Memory.ECCErrors.L1}}
	      L2 Cache:  	{{$s.Memory.ECCErrors.L2}}
	      Memory:  	{{$s.Memory.ECCErrors.Global}}{{end}}
	  PCI
	    BAR1:  	{{$s.PCI.BAR1Used}} / {{.PCI.BAR1}} MiB
	    Throughput{{if not $s.PCI.Throughput}}:  	N/A{{else}}
	      RX:  	{{$s.PCI.Throughput.RX}} KB/s
	      TX:  	{{$s.PCI.Throughput.TX}} KB/s{{end}}
	  Clocks
	    SM:  	{{$s.Clocks.SM}} MHz
	    Memory:  	{{$s.Clocks.Memory}} MHz
	    Graphics:  	{{$s.Clocks.Graphics}} MHz
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
	if err := writeStatusJSON(&body); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(body.Bytes())
}

func writeStatusJSON(wr io.Writer) error {
	status := make([]*nvidia.DeviceStatus, 0, len(Devices))
	for i := range Devices {
		s, err := Devices[i].Status()
		if err != nil {
			return err
		}
		status = append(status, s)
	}

	r := struct{ Devices []*nvidia.DeviceStatus }{status}
	assert(json.NewEncoder(wr).Encode(r))
	return nil
}

func (r *remoteV10) dockerCLI(resp http.ResponseWriter, req *http.Request) {
	var body bytes.Buffer

	ids := strings.Split(req.FormValue("dev"), " ")
	if err := dockerCLIDevice(&body, ids); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	names := strings.Split(req.FormValue("vol"), " ")
	if err := dockerCLIVolume(&body, names); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Write(body.Bytes())
}

func dockerCLIDevice(wr io.Writer, ids []string) error {
	const tpl = "--device=/dev/nvidiactl --device=/dev/nvidia-uvm {{range .}}--device={{.}} {{end}}"

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

func dockerCLIVolume(wr io.Writer, names []string) error {
	const tpl = "--volume-driver=nvidia {{range .}}--volume={{.}} {{end}}"

	vols := make([]string, 0, len(Volumes))

	if len(names) == 1 && (names[0] == "*" || names[0] == "") {
		for _, v := range Volumes {
			vols = append(vols, fmt.Sprintf("%s:%s", v.Name, v.Mountpoint))
		}
	} else {
		for _, n := range names {
			v, ok := Volumes[n]
			if !ok {
				return fmt.Errorf("invalid volume: %s", n)
			}
			vols = append(vols, fmt.Sprintf("%s:%s", v.Name, v.Mountpoint))
		}
	}

	t := template.Must(template.New("").Parse(tpl))
	assert(t.Execute(wr, vols))
	return nil
}
