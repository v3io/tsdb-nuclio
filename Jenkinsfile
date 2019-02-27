label = "${UUID.randomUUID().toString()}"
BUILD_FOLDER = "/go"
git_project = "tsdb-nuclio"
git_project_user = "gkirok"
git_deploy_user_token = "iguazio-dev-git-user-token"
git_deploy_user_private_key = "iguazio-dev-git-user-private-key"

podTemplate(label: "${git_project}-${label}", inheritFrom: "jnlp-docker") {
    node("${git_project}-${label}") {
        withCredentials([
                string(credentialsId: git_deploy_user_token, variable: 'GIT_TOKEN')
        ]) {
            def TAG_VERSION
            def DOCKER_TAG_VERSION
            pipelinex = library(identifier: 'pipelinex@shellc', retriever: modernSCM(
                    [$class       : 'GitSCMSource',
                     credentialsId: git_deploy_user_private_key,
                     remote       : "git@github.com:iguazio/pipelinex.git"])).com.iguazio.pipelinex
            multi_credentials = [pipelinex.DockerRepoDev.ARTIFACTORY_IGUAZIO, pipelinex.DockerRepoDev.DOCKER_HUB, pipelinex.DockerRepoDev.QUAY_IO]

            common.notify_slack {
                stage('get tag data') {
                    container('jnlp') {
                        TAG_VERSION = github.get_tag_version(TAG_NAME)
                        DOCKER_TAG_VERSION = github.get_docker_tag_version(TAG_NAME)

                        echo "$TAG_VERSION"
                        echo "$DOCKER_TAG_VERSION"
                    }
                }

                if (github.check_tag_expiration(git_project, git_project_user, TAG_VERSION, GIT_TOKEN)) {
                    stage('prepare sources') {
                        container('jnlp') {
                            dir("${BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                git(changelog: false, credentialsId: git_deploy_user_private_key, poll: false, url: "git@github.com:${git_project_user}/${git_project}.git")
                                sh("git checkout ${TAG_VERSION}")
                            }
                        }
                    }

                    parallel(
                            'build tsdb-ingest': {
                                container('docker-cmd') {
                                    dir("${BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                        sh("TSDB_DOCKER_REPO= TSDB_TAG=${DOCKER_TAG_VERSION} make ingest")
                                    }
                                }

                                container('docker-cmd') {
                                    dockerx.images_push_multi_registries(["tsdb-ingest:${DOCKER_TAG_VERSION}"], multi_credentials)
                                }
                            },

                            'build tsdb-query': {
                                container('docker-cmd') {
                                    dir("${BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                        sh("TSDB_DOCKER_REPO= TSDB_TAG=${DOCKER_TAG_VERSION} make query")
                                    }
                                }

                                container('docker-cmd') {
                                    dockerx.images_push_multi_registries(["tsdb-query:${DOCKER_TAG_VERSION}"], multi_credentials)
                                }
                            }
                    )

                    stage('update release status') {
                        container('jnlp') {
                            github.update_release_status(git_project, git_project_user, "${TAG_VERSION}", GIT_TOKEN)
                        }
                    }
                }
            }
        }
    }
}
