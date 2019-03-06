label = "${UUID.randomUUID().toString()}"
git_project = "tsdb-nuclio"
git_project_user = "gkirok"
git_deploy_user_token = "iguazio-prod-git-user-token"
git_deploy_user_private_key = "iguazio-prod-git-user-private-key"

podTemplate(label: "${git_project}-${label}", inheritFrom: "jnlp-docker") {
    node("${git_project}-${label}") {
        pipelinex = library(identifier: 'pipelinex@reduction', retriever: modernSCM(
                [$class       : 'GitSCMSource',
                 credentialsId: git_deploy_user_private_key,
                 remote       : "git@github.com:iguazio/pipelinex.git"])).com.iguazio.pipelinex
        common.notify_slack {
            withCredentials([
                    string(credentialsId: git_deploy_user_token, variable: 'GIT_TOKEN')
            ]) {
                github.init_project(git_project, git_project_user, GIT_TOKEN) {
                    stage('prepare sources') {
                        container('jnlp') {
                            dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                git(changelog: false, credentialsId: git_deploy_user_private_key, poll: false, url: "git@github.com:${git_project_user}/${git_project}.git")
                                common.shellc("git checkout ${github.TAG_VERSION}")
                            }
                        }
                    }

                    parallel(
                            'build tsdb-ingest in dood': {
                                container('docker-cmd') {
                                    dir("${BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                        sh("TSDB_DOCKER_REPO= TSDB_TAG=${DOCKER_TAG_VERSION} make ingest")
                                    }
                                }
                            },

                            'build tsdb-query in dood': {
                                container('docker-cmd') {
                                    dir("${BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                        sh("TSDB_DOCKER_REPO= TSDB_TAG=${DOCKER_TAG_VERSION} make query")
                                    }
                                }
                            }
                    )

                    stage('push') {
                        container('docker-cmd') {
                            dockerx.images_push_multi_registries(["tsdb-ingest:${github.DOCKER_TAG_VERSION}","tsdb-query:${github.DOCKER_TAG_VERSION}"], [pipelinex.DockerRepoDev.ARTIFACTORY_IGUAZIO, pipelinex.DockerRepoDev.DOCKER_HUB, pipelinex.DockerRepoDev.QUAY_IO])
                        }
                    }
                }
            }
        }
    }
}
