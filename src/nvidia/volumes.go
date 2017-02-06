// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package nvidia

import (
	"bufio"
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/NVIDIA/nvidia-docker/src/ldcache"
)

const (
	binDir   = "bin"
	lib32Dir = "lib"
	lib64Dir = "lib64"
)

type components map[string][]string

type volumeDir struct {
	name  string
	files []string
}

type VolumeInfo struct {
	Name         string
	Mountpoint   string
	MountOptions string
	Components   components
}

type Volume struct {
	*VolumeInfo

	Path    string
	Version string
	dirs    []volumeDir
}

type VolumeMap map[string]*Volume

type FileCloneStrategy interface {
	Clone(src, dst string) error
}

type LinkStrategy struct{}

func (s LinkStrategy) Clone(src, dst string) error {
	return os.Link(src, dst)
}

type LinkOrCopyStrategy struct{}

func (s LinkOrCopyStrategy) Clone(src, dst string) error {
	// Prefer hard link, fallback to copy
	err := os.Link(src, dst)
	if err != nil {
		err = Copy(src, dst)
	}
	return err
}

func Copy(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	fi, err := s.Stat()
	if err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	if err := d.Chmod(fi.Mode()); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

var Volumes = []VolumeInfo{
	{
		"nvidia_driver",
		"/usr/local/nvidia",
		"ro",
		components{
			"binaries": {
				//"nvidia-modprobe",       // Kernel module loader
				//"nvidia-settings",       // X server settings
				//"nvidia-xconfig",        // X xorg.conf editor
				"nvidia-cuda-mps-control", // Multi process service CLI
				"nvidia-cuda-mps-server",  // Multi process service server
				"nvidia-debugdump",        // GPU coredump utility
				"nvidia-persistenced",     // Persistence mode utility
				"nvidia-smi",              // System management interface
			},
			"libraries": {
				// ------- X11 -------

				//"libnvidia-cfg.so",  // GPU configuration (used by nvidia-xconfig)
				//"libnvidia-gtk2.so", // GTK2 (used by nvidia-settings)
				//"libnvidia-gtk3.so", // GTK3 (used by nvidia-settings)
				//"libnvidia-wfb.so",  // Wrapped software rendering module for X server
				//"libglx.so",         // GLX extension module for X server

				// ----- Compute -----

				"libnvidia-ml.so",              // Management library
				"libcuda.so",                   // CUDA driver library
				"libnvidia-ptxjitcompiler.so",  // PTX-SASS JIT compiler (used by libcuda)
				"libnvidia-fatbinaryloader.so", // fatbin loader (used by libcuda)
				"libnvidia-opencl.so",          // NVIDIA OpenCL ICD
				"libnvidia-compiler.so",        // NVVM-PTX compiler for OpenCL (used by libnvidia-opencl)
				//"libOpenCL.so",               // OpenCL ICD loader

				// ------ Video ------

				"libvdpau_nvidia.so",  // NVIDIA VDPAU ICD
				"libnvidia-encode.so", // Video encoder
				"libnvcuvid.so",       // Video decoder
				"libnvidia-fbc.so",    // Framebuffer capture
				"libnvidia-ifr.so",    // OpenGL framebuffer capture

				// ----- Graphic -----

				// XXX In an ideal world we would only mount nvidia_* vendor specific libraries and
				// install ICD loaders inside the container. However, for backward compatibility reason
				// we need to mount everything. This will hopefully change once GLVND is well established.

				"libGL.so",         // OpenGL/GLX legacy _or_ compatibility wrapper (GLVND)
				"libGLX.so",        // GLX ICD loader (GLVND)
				"libOpenGL.so",     // OpenGL ICD loader (GLVND)
				"libGLESv1_CM.so",  // OpenGL ES v1 common profile legacy _or_ ICD loader (GLVND)
				"libGLESv2.so",     // OpenGL ES v2 legacy _or_ ICD loader (GLVND)
				"libEGL.so",        // EGL ICD loader
				"libGLdispatch.so", // OpenGL dispatch (GLVND) (used by libOpenGL, libEGL and libGLES*)

				"libGLX_nvidia.so",         // OpenGL/GLX ICD (GLVND)
				"libEGL_nvidia.so",         // EGL ICD (GLVND)
				"libGLESv2_nvidia.so",      // OpenGL ES v2 ICD (GLVND)
				"libGLESv1_CM_nvidia.so",   // OpenGL ES v1 common profile ICD (GLVND)
				"libnvidia-eglcore.so",     // EGL core (used by libGLES* or libGLES*_nvidia and libEGL_nvidia)
				"libnvidia-egl-wayland.so", // EGL wayland extensions (used by libEGL_nvidia)
				"libnvidia-glcore.so",      // OpenGL core (used by libGL or libGLX_nvidia)
				"libnvidia-tls.so",         // Thread local storage (used by libGL or libGLX_nvidia)
				"libnvidia-glsi.so",        // OpenGL system interaction (used by libEGL_nvidia)
			},
		},
	},
}

func blacklisted(file string, obj *elf.File) (bool, error) {
	lib := regexp.MustCompile(`^.*/lib([\w-]+)\.so[\d.]*$`)
	glcore := regexp.MustCompile(`libnvidia-e?glcore\.so`)
	gldispatch := regexp.MustCompile(`libGLdispatch\.so`)

	if m := lib.FindStringSubmatch(file); m != nil {
		switch m[1] {

		// Blacklist EGL/OpenGL libraries issued by other vendors
		case "EGL":
			fallthrough
		case "GLESv1_CM":
			fallthrough
		case "GLESv2":
			fallthrough
		case "GL":
			deps, err := obj.DynString(elf.DT_NEEDED)
			if err != nil {
				return false, err
			}
			for _, d := range deps {
				if glcore.MatchString(d) || gldispatch.MatchString(d) {
					return false, nil
				}
			}
			return true, nil

		// Blacklist TLS libraries using the old ABI (!= 2.3.99)
		case "nvidia-tls":
			const abi = 0x6300000003
			s, err := obj.Section(".note.ABI-tag").Data()
			if err != nil {
				return false, err
			}
			return binary.LittleEndian.Uint64(s[24:]) != abi, nil
		}
	}
	return false, nil
}

func (v *Volume) Create(s FileCloneStrategy) (err error) {
	root := path.Join(v.Path, v.Version)
	if err = os.MkdirAll(root, 0755); err != nil {
		return
	}
	defer func() {
		if err != nil {
			v.Remove()
		}
	}()

	for _, d := range v.dirs {
		vpath := path.Join(root, d.name)
		if err := os.MkdirAll(vpath, 0755); err != nil {
			return err
		}

		// For each file matching the volume components (blacklist excluded), create a hardlink/copy
		// of it inside the volume directory. We also need to create soname symlinks similar to what
		// ldconfig does since our volume will only show up at runtime.
		for _, f := range d.files {
			obj, err := elf.Open(f)
			if err != nil {
				return fmt.Errorf("%s: %v", f, err)
			}
			defer obj.Close()

			ok, err := blacklisted(f, obj)
			if err != nil {
				return fmt.Errorf("%s: %v", f, err)
			}
			if ok {
				continue
			}

			l := path.Join(vpath, path.Base(f))
			if err := s.Clone(f, l); err != nil {
				return err
			}
			soname, err := obj.DynString(elf.DT_SONAME)
			if err != nil {
				return fmt.Errorf("%s: %v", f, err)
			}
			if len(soname) > 0 {
				l = path.Join(vpath, soname[0])
				if err := os.Symlink(path.Base(f), l); err != nil && !os.IsExist(err) {
					return err
				}
				// XXX Many applications (wrongly) assume that libcuda.so exists (e.g. with dlopen)
				// Hardcode the libcuda symlink for the time being.
				if strings.HasPrefix(soname[0], "libcuda") {
					l = strings.TrimRight(l, ".0123456789")
					if err := os.Symlink(path.Base(f), l); err != nil && !os.IsExist(err) {
						return err
					}
				}
				// XXX GLVND requires this symlink for indirect GLX support
				// It won't be needed once we have an indirect GLX vendor neutral library.
				if strings.HasPrefix(soname[0], "libGLX_nvidia") {
					l = strings.Replace(l, "GLX_nvidia", "GLX_indirect", 1)
					if err := os.Symlink(path.Base(f), l); err != nil && !os.IsExist(err) {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (v *Volume) Remove(version ...string) error {
	vv := v.Version
	if len(version) == 1 {
		vv = version[0]
	}
	return os.RemoveAll(path.Join(v.Path, vv))
}

func (v *Volume) Exists(version ...string) (bool, error) {
	vv := v.Version
	if len(version) == 1 {
		vv = version[0]
	}
	_, err := os.Stat(path.Join(v.Path, vv))
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (v *Volume) ListVersions() ([]string, error) {
	dirs, err := ioutil.ReadDir(v.Path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	versions := make([]string, len(dirs))
	for i := range dirs {
		versions[i] = dirs[i].Name()
	}
	return versions, nil
}

func which(bins ...string) ([]string, error) {
	paths := make([]string, 0, len(bins))

	out, _ := exec.Command("which", bins...).Output()
	r := bufio.NewReader(bytes.NewBuffer(out))
	for {
		p, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if p = strings.TrimSpace(p); !path.IsAbs(p) {
			continue
		}
		path, err := filepath.EvalSymlinks(p)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func LookupVolumes(prefix string) (vols VolumeMap, err error) {
	drv, err := GetDriverVersion()
	if err != nil {
		return nil, err
	}
	cache, err := ldcache.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := cache.Close(); err == nil {
			err = e
		}
	}()

	vols = make(VolumeMap, len(Volumes))

	for i := range Volumes {
		vol := &Volume{
			VolumeInfo: &Volumes[i],
			Path:       path.Join(prefix, Volumes[i].Name),
			Version:    drv,
		}

		for t, c := range vol.Components {
			switch t {
			case "binaries":
				bins, err := which(c...)
				if err != nil {
					return nil, err
				}
				vol.dirs = append(vol.dirs, volumeDir{binDir, bins})
			case "libraries":
				libs32, libs64 := cache.Lookup(c...)
				vol.dirs = append(vol.dirs,
					volumeDir{lib32Dir, libs32},
					volumeDir{lib64Dir, libs64},
				)
			}
		}
		vols[vol.Name] = vol
	}
	return
}
