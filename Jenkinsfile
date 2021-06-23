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
        def scmInfo

        stage('checkout') {
            scmInfo = checkout(scm)
        }

        stage('dependencies') {
            container('docker') {
                sh 'apk add --no-cache make bash git'
            }
        }

        def versionInfo
        stage('version') {
            container('docker') {
                versionInfo = getVersionInfo(scmInfo)
                println "versionInfo=${versionInfo}"
            }
        }

        def dist = 'ubuntu20.04'
        def arch = 'amd64'
        def stageLabel = "${dist}-${arch}"

        stage('build-one') {
            container('docker') {
                stage (stageLabel) {
                    sh "make ADD_DOCKER_PLATFORM_ARGS=true ${dist}-${arch}"
                }
            }
        }

        stage('release') {
            container('docker') {
                stage (stageLabel) {

                    def component = 'main'
                    def repository = 'sw-gpu-cloudnative-debian-local/pool/main/'

                    def uploadSpec = """{
                                        "files":
                                        [  {
                                                "pattern": "./dist/${dist}/${arch}/*.deb",
                                                "target": "${repository}",
                                                "props": "deb.distribution=${dist};deb.component=${component};deb.architecture=${arch}"
                                            }
                                        ]
                                    }"""

                    sh "echo starting release with versionInfo=${versionInfo}"
                    if (versionInfo.isTag) {
                        // upload to artifactory repository
                        def server = Artifactory.server 'sw-gpu-artifactory'
                        server.upload spec: uploadSpec
                    } else {
                        sh "echo skipping release for non-tagged build"
                    }
                }
            }
        }
    }
}

// getVersionInfo returns a hash of version info
def getVersionInfo(def scmInfo) {
    def versionInfo = [
        isTag: isTag(scmInfo.GIT_BRANCH)
    ]

    scmInfo.each { k, v -> versionInfo[k] = v }
    return versionInfo
}

def isTag(def branch) {
    if (!branch.startsWith('v')) {
        return false
    }

    def version = shOutput('git describe --all --exact-match --always')
    return version == "tags/${branch}"
}

def shOuptut(def script) {
    return sh(script: script, returnStdout: true).trim()
}
