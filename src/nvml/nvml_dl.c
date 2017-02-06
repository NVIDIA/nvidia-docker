// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

#include <stddef.h>
#include <dlfcn.h>

#include "nvml_dl.h"

#define DLSYM(x, sym)                           \
do {                                            \
    dlerror();				        \
    x = dlsym(handle, #sym);                    \
    if (dlerror() != NULL) {                    \
        return (NVML_ERROR_FUNCTION_NOT_FOUND); \
    }                                           \
} while (0)

typedef nvmlReturn_t (*nvmlSym_t)();

static void *handle;

nvmlReturn_t NVML_DL(nvmlInit)(void)
{
    handle = dlopen("libnvidia-ml.so.1", RTLD_LAZY | RTLD_GLOBAL);
    if (handle == NULL) {
	return (NVML_ERROR_LIBRARY_NOT_FOUND);
    }
    return (nvmlInit());
}

nvmlReturn_t NVML_DL(nvmlShutdown)(void)
{
    nvmlReturn_t r = nvmlShutdown();
    if (r != NVML_SUCCESS) {
	return (r);
    }
    return (dlclose(handle) ? NVML_ERROR_UNKNOWN : NVML_SUCCESS);
}

nvmlReturn_t NVML_DL(nvmlDeviceGetTopologyCommonAncestor)(
  nvmlDevice_t dev1, nvmlDevice_t dev2, nvmlGpuTopologyLevel_t *info)
{
    nvmlSym_t sym;

    DLSYM(sym, nvmlDeviceGetTopologyCommonAncestor);
    return ((*sym)(dev1, dev2, info));
}
