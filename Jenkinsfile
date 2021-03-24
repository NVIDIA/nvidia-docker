/*
# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

podTemplate (cloud:'sw-gpu-cloudnative',
    containers: [
    containerTemplate(name: 'docker', image: 'docker:dind', ttyEnabled: true, privileged: true)
  ]) {
    node(POD_LABEL) {
        stage('checkout') {
            checkout scm
        }
        stage('dependencies') {
            container('docker') {
                sh 'apk add --no-cache make bash'
            }
        }
        stage('build-one') {
            parallel (
                getSingleBuildForArchitectures(["amd64", "ppc64le", "arm64"])
            )
        }
        stage('build-all') {
            parallel (
                getAllBuildForArchitectures(["amd64", "ppc64le", "arm64", "x86_64", "aarch64"])
            )
        }
    }
}

def getBuildClosure(def architecture, def makeCommand, def makeTarget) {
    return {
        container('docker') {
            stage(architecture) {
                sh "${makeCommand} ${makeTarget}"
            }
        }
    }
}

def getBuildStagesForArchitectures(def architectures, def makeCommand, def makeTargetPrefix) {
    stages = [:]

    for (a in architectures) {
        stages[a] = getBuildClosure(a, makeCommand, "${makeTargetPrefix}-${a}")
    }

    return stages
}

def getSingleBuildForArchitectures(def architectures) {
    return getBuildStagesForArchitectures(architectures, "make", "ubuntu18.04")
}

def getAllBuildForArchitectures(def architectures) {
    // TODO: For the time being we only echo the command for the "all" stages
    return getBuildStagesForArchitectures(architectures, "echo make", "docker")
}