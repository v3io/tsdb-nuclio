def label = "${UUID.randomUUID().toString()}"
def BUILD_FOLDER = '/go'
def github_user = "gkirok"
def docker_user = "gallziguazio"


properties([pipelineTriggers([[$class: 'PeriodicFolderTrigger', interval: '2m']])])
podTemplate(label: "tsdb-nuclio-${label}", inheritFrom: 'kube-slave-dood') {
    node("tsdb-nuclio-${label}") {
        withCredentials([
                usernamePassword(credentialsId: '4318b7db-a1af-4775-b871-5a35d3e75c21', passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME')
        ]) {
            stage('release') {
                def TAG_VERSION = sh(
                        script: "echo ${TAG_NAME} | tr -d '\\n' | egrep '^v[\\.0-9]*.*-v[\\.0-9]*\$'",
                        returnStdout: true
                ).trim()
                if ( TAG_VERSION ) {
                    print TAG_VERSION
                    def V3IO_TSDB_VERSION = sh(
                            script: "echo ${TAG_VERSION} | awk -F '-v' '{print \"v\"\$2}'",
                            returnStdout: true
                    ).trim()

                    stage('prepare sources') {
                        sh """ 
                                cd ${BUILD_FOLDER}
                                git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/${github_user}/tsdb-nuclio.git src/github.com/v3io/tsdb-nuclio
                                cd ${BUILD_FOLDER}/src/github.com/v3io/tsdb-nuclio
                                rm -rf functions/ingest/vendor/github.com/v3io/v3io-tsdb functions/query/vendor/github.com/v3io/v3io-tsdb
                                git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/${github_user}/v3io-tsdb.git functions/ingest/vendor/github.com/v3io/v3io-tsdb
                                cd functions/ingest/vendor/github.com/v3io/v3io-tsdb
                                git checkout ${V3IO_TSDB_VERSION}
                                rm -rf .git vendor/github.com/v3io vendor/github.com/nuclio
                                cd ${BUILD_FOLDER}/src/github.com/v3io/tsdb-nuclio
                                cp -R functions/ingest/vendor/github.com/v3io/v3io-tsdb functions/query/vendor/github.com/v3io/v3io-tsdb
                        """
                    }

                    stage('build in dood') {
                        container('docker-cmd') {
                            sh """
                                    cd ${BUILD_FOLDER}/src/github.com/v3io/tsdb-nuclio/functions/ingest
                                    docker build . --tag tsdb-ingest:latest --tag ${docker_user}/tsdb-ingest:${TAG_VERSION}

                                    cd ${BUILD_FOLDER}/src/github.com/v3io/tsdb-nuclio/functions/query
                                    docker build . --tag tsdb-query:latest --tag ${docker_user}/tsdb-query:${TAG_VERSION}
                            """
                            withDockerRegistry([credentialsId: "472293cc-61bc-4e9f-aecb-1d8a73827fae", url: ""]) {
                                sh "docker push ${docker_user}/tsdb-ingest:${TAG_VERSION}"
                                sh "docker push ${docker_user}/tsdb-query:${TAG_VERSION}"
                            }
                        }
                    }

                    stage('git push') {
                        try {
                            sh """
                                git config --global user.email '${GIT_USERNAME}@iguazio.com'
                                git config --global user.name '${GIT_USERNAME}'
                                cd ${BUILD_FOLDER}/src/github.com/v3io/tsdb-nuclio
                                git add *
                                git commit -am 'Updated TSDB to ${V3IO_TSDB_VERSION}';
                                git push origin master
                            """
                        } catch (err) {
                            echo "Can not push code to git"
                        }
                    }
                } else {
                    echo "${TAG_VERSION} is not release tag."
                }
            }
        }
    }
}