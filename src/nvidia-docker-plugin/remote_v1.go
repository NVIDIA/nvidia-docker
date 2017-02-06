// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/NVIDIA/nvidia-docker/src/nvidia"
)

type remoteV10 struct{}

func (r *remoteV10) version() string { return "v1.0" }

func (r *remoteV10) gpuInfo(resp http.ResponseWriter, req *http.Request) {
	const tpl = `
	Driver version:  	{{driverVersion}}
	Supported CUDA version:  	{{cudaVersion}}
	{{range $i, $e := .}}
	Device #{{$i}}
	  Model:  	{{or .Model "N/A"}}
	  UUID:  	{{.UUID}}
	  Path:  	{{.Path}}
	  Family: 	{{or .Family "N/A"}}
	  Arch:  	{{or .Arch "N/A"}}
	  Cores:  	{{or .Cores "N/A"}}
	  Power:  	{{if .Power}}{{.Power}} W{{else}}N/A{{end}}
	  CPU Affinity:  	{{if .CPUAffinity}}NUMA node{{.CPUAffinity}}{{else}}N/A{{end}}
	  PCI
	    Bus ID:  	{{.PCI.BusID}}
	    BAR1:  	{{if .PCI.BAR1}}{{.PCI.BAR1}} MiB{{else}}N/A{{end}}
	    Bandwidth:  	{{if .PCI.Bandwidth}}{{.PCI.Bandwidth}} MB/s{{else}}N/A{{end}}
	  Memory
	    ECC:  	{{or .Memory.ECC "N/A"}}
	    Global:  	{{if .Memory.Global}}{{.Memory.Global}} MiB{{else}}N/A{{end}}
	    Constant:  	{{if .Memory.Constant}}{{.Memory.Constant}} KiB{{else}}N/A{{end}}
	    Shared:  	{{if .Memory.Shared}}{{.Memory.Shared}} KiB{{else}}N/A{{end}}
	    L2 Cache:  	{{if .Memory.L2Cache}}{{.Memory.L2Cache}} KiB{{else}}N/A{{end}}
	    Bandwidth:  	{{if .Memory.Bandwidth}}{{.Memory.Bandwidth}} MB/s{{else}}N/A{{end}}
	  Clocks
	    Cores:  	{{if .Clocks.Cores}}{{.Clocks.Cores}} MHz{{else}}N/A{{end}}
	    Memory:  	{{if .Clocks.Memory}}{{.Clocks.Memory}} MHz{{else}}N/A{{end}}
	  P2P Available{{if not .Topology}}:  	None{{else}}{{range .Topology}}
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

	writeGPUInfoJSON(&body)
	resp.Header().Set("Content-Type", "application/json")
	_, err := body.WriteTo(resp)
	assert(err)
}

func writeGPUInfoJSON(wr io.Writer) {
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
	  Power:  	{{if and $s.Power .Power}}{{$s.Power}} / {{.Power}} W{{else}}N/A{{end}}
	  Temperature:  	{{if $s.Temperature}}{{$s.Temperature}} °C{{else}}N/A{{end}}
	  Utilization
	    GPU:  	{{if $s.Utilization.GPU}}{{$s.Utilization.GPU}} %{{else}}N/A{{end}}
	    Memory:  	{{if $s.Utilization.Memory}}{{$s.Utilization.Memory}} %{{else}}N/A{{end}}
	    Encoder:  	{{if $s.Utilization.Encoder}}{{$s.Utilization.Encoder}} %{{else}}N/A{{end}}
	    Decoder:  	{{if $s.Utilization.Decoder}}{{$s.Utilization.Decoder}} %{{else}}N/A{{end}}
	  Memory
	    Global:  	{{if and $s.Memory.GlobalUsed .Memory.Global}}{{$s.Memory.GlobalUsed}} / {{.Memory.Global}} MiB{{else}}N/A{{end}}
	    ECC Errors
	      L1 Cache:  	{{or $s.Memory.ECCErrors.L1Cache "N/A"}}
	      L2 Cache:  	{{or $s.Memory.ECCErrors.L2Cache "N/A"}}
	      Global:  	{{or $s.Memory.ECCErrors.Global "N/A"}}
	  PCI
	    BAR1:  	{{if and $s.PCI.BAR1Used .PCI.BAR1}}{{$s.PCI.BAR1Used}} / {{.PCI.BAR1}} MiB{{else}}N/A{{end}}
	    Throughput
	      RX:  	{{if $s.PCI.Throughput.RX}}{{$s.PCI.Throughput.RX}} MB/s{{else}}N/A{{end}}
	      TX:  	{{if $s.PCI.Throughput.TX}}{{$s.PCI.Throughput.TX}} MB/s{{else}}N/A{{end}}
	  Clocks
	    Cores:  	{{if $s.Clocks.Cores}}{{$s.Clocks.Cores}} MHz{{else}}N/A{{end}}
	    Memory:  	{{if $s.Clocks.Memory}}{{$s.Clocks.Memory}} MHz{{else}}N/A{{end}}
	  Processes{{if not $s.Processes}}:  	None{{else}}{{range $s.Processes}}
	    - PID:  	{{.PID}}
	      Name:  	{{.Name}}
	      Memory:  	{{.MemoryUsed}} MiB{{end}}{{end}}
	{{end}}
	`
	t := template.Must(template.New("").Parse(tpl))
	w := tabwriter.NewWriter(resp, 0, 4, 0, ' ', 0)

	assert(t.Execute(w, Devices))
	assert(w.Flush())
}

func (r *remoteV10) gpuStatusJSON(resp http.ResponseWriter, req *http.Request) {
	var body bytes.Buffer

	writeGPUStatusJSON(&body)
	resp.Header().Set("Content-Type", "application/json")
	_, err := body.WriteTo(resp)
	assert(err)
}

func writeGPUStatusJSON(wr io.Writer) {
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
	const tpl = "--volume-driver={{.VolumeDriver}}{{range .Volumes}} --volume={{.}}{{end}}" +
		"{{range .Devices}} --device={{.}}{{end}}"

	devs := strings.Split(req.FormValue("dev"), " ")
	vols := strings.Split(req.FormValue("vol"), " ")

	args, err := dockerCLIArgs(devs, vols)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	t := template.Must(template.New("").Parse(tpl))
	assert(t.Execute(resp, args))
}

func (r *remoteV10) dockerCLIJSON(resp http.ResponseWriter, req *http.Request) {
	devs := strings.Split(req.FormValue("dev"), " ")
	vols := strings.Split(req.FormValue("vol"), " ")

	args, err := dockerCLIArgs(devs, vols)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	assert(json.NewEncoder(resp).Encode(args))
}

type dockerArgs struct {
	VolumeDriver string
	Volumes      []string
	Devices      []string
}

func dockerCLIArgs(devs, vols []string) (*dockerArgs, error) {
	cdevs, err := nvidia.GetControlDevicePaths()
	if err != nil {
		return nil, err
	}
	devs, err = dockerCLIDevices(devs)
	if err != nil {
		return nil, err
	}
	vols, err = dockerCLIVolumes(vols)
	if err != nil {
		return nil, err
	}
	return &dockerArgs{
		VolumeDriver: nvidia.DockerPlugin,
		Volumes:      vols,
		Devices:      append(cdevs, devs...),
	}, nil
}

func dockerCLIDevices(ids []string) ([]string, error) {
	devs := make([]string, 0, len(Devices))

	if len(ids) == 1 && (ids[0] == "*" || ids[0] == "") {
		for i := range Devices {
			devs = append(devs, Devices[i].Path)
		}
	} else {
		d, err := nvidia.FilterDevices(Devices, ids)
		if err != nil {
			return nil, err
		}
		for i := range d {
			devs = append(devs, d[i].Path)
		}
	}
	return devs, nil
}

func dockerCLIVolumes(names []string) ([]string, error) {
	vols := make([]string, 0, len(Volumes))

	drv, err := nvidia.GetDriverVersion()
	if err != nil {
		return nil, err
	}
	if len(names) == 1 && (names[0] == "*" || names[0] == "") {
		for _, v := range Volumes {
			vols = append(vols, fmt.Sprintf("%s_%s:%s:%s", v.Name, drv, v.Mountpoint, v.MountOptions))
		}
	} else {
		for _, n := range names {
			v, ok := Volumes[n]
			if !ok {
				return nil, fmt.Errorf("invalid volume: %s", n)
			}
			vols = append(vols, fmt.Sprintf("%s_%s:%s:%s", v.Name, drv, v.Mountpoint, v.MountOptions))
		}
	}
	return vols, nil
}

func (r *remoteV10) mesosCLI(resp http.ResponseWriter, req *http.Request) {
	const format = "--attributes=gpus:%s --resources=gpus:{%s}"

	// Generate Mesos attributes
	var b bytes.Buffer
	writeGPUInfoJSON(&b)
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
