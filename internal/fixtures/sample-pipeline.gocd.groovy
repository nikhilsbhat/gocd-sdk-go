package artifacts_and_html_reports

import cd.go.contrib.plugins.configrepo.groovy.dsl.GoCD

GoCD.script {
    pipelines {
        pipeline('website') {
            group = "example-group"

            trackingTool {
                link = 'https://github.com/gocd/api.go.cd/issues/${ID}'
                regex = ~/##(\\d+)/
            }

            materials {
                git {
                    url = 'https://github.com/gocd/api.go.cd'
                    branch = 'release-18.1.0'
                }
            }
            stages {
                stage('build-website') {
                    jobs {
                        job('build') {
                            tasks {
                                bash {
                                    commandString = 'bundle install --path .bundle -j4'
                                }
                                bash {
                                    commandString = 'bundle exec rake build'
                                }
                            }

                            artifacts {
                                build {
                                    source = 'build/18.1.0'
                                    destination = 'website'
                                }
                            }

                            tabs {
                                tab('website') { path = 'website/18.1.0/index.html' }
                            }
                        }
                    }
                }
            }
        }
    }
}
