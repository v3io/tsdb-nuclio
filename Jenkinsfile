label = "${UUID.randomUUID().toString()}"
git_project = "tsdb-nuclio"
git_project_user = "v3io"
git_project_upstream_user = "iguazio"
git_deploy_user = "iguazio-prod-git-user"
git_deploy_user_token = "iguazio-prod-git-user-token"
git_deploy_user_private_key = "iguazio-prod-git-user-private-key"

podTemplate(label: "${git_project}-${label}", inheritFrom: "jnlp-docker") {
    node("${git_project}-${label}") {
        pipelinex = library(identifier: 'pipelinex@development', retriever: modernSCM(
                [$class       : 'GitSCMSource',
                 credentialsId: git_deploy_user_private_key,
                 remote       : "git@github.com:iguazio/pipelinex.git"])).com.iguazio.pipelinex
        common.notify_slack {
            withCredentials([
                    string(credentialsId: git_deploy_user_token, variable: 'GIT_TOKEN')
            ]) {
                github.release(git_deploy_user, git_project, git_project_user, git_project_upstream_user, true, GIT_TOKEN) {
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
                                    dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                        common.shellc("TSDB_DOCKER_REPO= TSDB_TAG=${github.DOCKER_TAG_VERSION} make ingest")
                                    }
                                }
                            },
                            'build tsdb-query in dood': {
                                container('docker-cmd') {
                                    dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                        common.shellc("TSDB_DOCKER_REPO= TSDB_TAG=${github.DOCKER_TAG_VERSION} make query")
                                    }
                                }
                            }
                    )

                    stage('push') {
                        container('docker-cmd') {
                            dockerx.images_push_multi_registries(["tsdb-ingest:${github.DOCKER_TAG_VERSION}","tsdb-query:${github.DOCKER_TAG_VERSION}"], [pipelinex.DockerRepo.ARTIFACTORY_IGUAZIO, pipelinex.DockerRepo.DOCKER_HUB, pipelinex.DockerRepo.QUAY_IO, pipelinex.DockerRepo.GCR_IO])
                        }
                    }
                }
            }
        }
    }
}
