import React, { useContext, useState } from 'react';
import { CodeEditor, Language } from '@patternfly/react-code-editor';
import { Button, AlertVariant, Modal, ModalVariant, Spinner, AlertGroup, Alert } from '@patternfly/react-core';
import YAML from 'yaml';
import { ReactReduxContext, useSelector } from 'react-redux';
import { createRepository, createTeam, listJiraProjects } from '@app/utils/APIService';
import { AlertInfo } from '@app/Teams/TeamsOnboarding';


export const Editor: React.FunctionComponent = () => {
    const [code, setCode] = useState("");
    const [alerts, setAlerts] = React.useState<AlertInfo[]>([]);
    const [isOpen, setIsOpen] = React.useState(false);
    const [creationLoading, setCreationLoading] = useState<boolean>(false);
    const { store } = useContext(ReactReduxContext);
    const dispatch = store.dispatch;

    const onEditorDidMount = (editor, monaco) => {
        // eslint-disable-next-line no-console
        editor.layout();
        editor.focus();
        monaco.editor.getModels()[0].updateOptions({ tabSize: 5 });
    };

    let currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);

    const teamExists = (team) => {
        if (currentTeamsAvailable.find(t => t.team_name == team.name)) {
            return true;
        }
        return false
    }

    const handleModalToggle = () => {
        setIsOpen(!isOpen);
        const parsed = YAML.parse(code)

        if (isValidYaml(parsed)) {
            window.location.reload();
        }
    };

    const onChange = value => {
        setCode(value)
    };

    const isValidYaml = (parsed) => {
        if (parsed.teams == undefined) {
            return false
        }

        for (const team of parsed.teams) {
            if (team.name == undefined || team.description == undefined
                || team.jira_projects == undefined || team.repositories == undefined) {
                return false
            }

            for (const repo of team.repositories) {
                if (repo.name == undefined || repo.organization == undefined) {
                    return false
                }
            }

        }

        return true
    }

    const filterJiraKeys = (keys: string[]) => {
        return new Promise((resolve, reject) => {
            let valid = new Array<string>

            listJiraProjects().then((res) => {
                if (res.code === 200) {
                    const result = res.data;

                    keys.forEach((key, idx) => {
                        if (result.find(el => el.project_key == key)) {
                            valid.push(key)
                        } else {
                            setAlerts(prevAlertInfo => [...prevAlertInfo, {
                                title: 'There is no JIRA key ' + key,
                                variant: AlertVariant.danger,
                                key: "jira-not-added",
                                index: key,
                            }]);
                        }
                    })
                    resolve(valid)
                }
            }).catch(error => {
                console.log(error);
                reject(valid);
            })
        });
    }

    const applyConfig = () => {
        setIsOpen(true)
        setCreationLoading(true)
        setAlerts([])

        try {
            const parsed = YAML.parse(code)

            if (isValidYaml(parsed)) {
                parsed.teams.forEach(async (team, idx) => {
                    // create team, if does not exist
                    if (!teamExists(team)) {
                       await filterJiraKeys(team.jira_projects).then((jiraKeys: any) => {
                            const data = {
                                "team_name": team.name,
                                "description": team.description,
                                "jira_keys": jiraKeys.join(","),
                            }

                            createTeam(data).then(response => {
                                if (response.code == 200) {
                                    setAlerts(prevAlertInfo => [...prevAlertInfo, {
                                        title: 'Team created ' + team.name,
                                        variant: AlertVariant.success,
                                        key: "team-created",
                                        index: idx,
                                    }]);
                                } else {
                                    setAlerts(prevAlertInfo => [...prevAlertInfo, {
                                        title: 'Could not create new team ' + team.name,
                                        variant: AlertVariant.danger,
                                        key: "team-not-created",
                                        index: idx,
                                    }]);
                                }
                            }).catch(error => {
                                console.log("NEW TEAM", data);
                                console.log(error);
                            })
                        })
                    }

                    // create repositories
                    team.repositories.forEach(repo => {
                        const data = {
                            git_organization: repo.organization,
                            repository_name: repo.name,
                            jobs: {
                                github_actions: {
                                    monitor: false
                                }
                            },
                            artifacts: [],
                            team_name: team.name
                        }

                        createRepository(data).then(response => {
                            if (response.code == 200) {
                                setAlerts(prevAlertInfo => [...prevAlertInfo, {
                                    title: 'Repository created ' + repo.organization + "/" + repo.name,
                                    variant: AlertVariant.success,
                                    key: "repo-created",
                                    index: idx,
                                }]);
                            } else {
                                setAlerts(prevAlertInfo => [...prevAlertInfo, {
                                    title: 'Could not add repository ' + repo.organization + "/" + repo.name + " in team " + team.name,
                                    variant: AlertVariant.danger,
                                    key: "repo-not-created",
                                    index: idx,
                                }]);
                            }
                        }).catch(error => {
                            console.log("NEW REPOSITORY", data);
                            console.log(error);
                        })
                    })
                });
            } else {
                setAlerts(prevAlertInfo => [...prevAlertInfo, {
                    title: 'Invalid Yaml. Check the sample, please.',
                    variant: AlertVariant.danger,
                    key: "invalid-yaml",
                }]);
            }
        } catch (err) {
            setAlerts(prevAlertInfo => [...prevAlertInfo, {
                title: 'Invalid Yaml. Check the console, please.',
                variant: AlertVariant.danger,
                key: "invalid-yaml",
            }]);
        }

        setCreationLoading(false)
    }


    return (
        <React.Fragment>
            <Modal
                variant={ModalVariant.small}
                title="Applying config..."
                isOpen={isOpen}
                onClose={handleModalToggle}
            >
                {creationLoading && <Spinner isSVG aria-label="Contents of the basic example" />}
                <AlertGroup aria-live="polite" aria-relevant="additions text" aria-atomic="false">
                    {alerts.map(({ title, variant, key }) => (
                        <Alert variant={variant} isInline isPlain title={title} key={key} />
                    ))}
                </AlertGroup>
            </Modal>
            <CodeEditor
                onCodeChange={onChange}
                language={Language.yaml}
                height="500px"
                isUploadEnabled
                isDownloadEnabled
                isCopyEnabled
                isLanguageLabelVisible
                isMinimapVisible
                onEditorDidMount={onEditorDidMount}
            />
            {code?.length > 0 && <div style={{ display: "flex", justifyContent: "right", marginTop: "10px" }}>
                <Button onClick={applyConfig} style={{ width: "100px" }} >
                    Save
                </Button>
            </div>}
            {code?.length == 0 && <div style={{ display: "flex", justifyContent: "right", marginTop: "10px" }}>
                <Button isDisabled style={{ width: "100px" }} >
                    Save
                </Button>
            </div>}
        </React.Fragment>
    );
};

