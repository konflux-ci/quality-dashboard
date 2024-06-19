import React, { useState } from "react";
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    ActionsColumn,
    IAction,
} from '@patternfly/react-table';
import { ITeam } from './TeamsSelect';
import { createUser, deleteInApi, updateTeam } from "@app/utils/APIService";
import { ReactReduxContext, useSelector } from 'react-redux';
import { Alert, AlertGroup, AlertVariant, Button, Form, FormGroup, Modal, ModalVariant, TextArea, TextInput } from "@patternfly/react-core";
import { AlertInfo, JiraProjects } from './TeamsOnboarding';
import { UserConfig } from "./User";
import { generateJiraConfig } from "./Configuration";
import { validate } from "@app/Jira/Jira";

export const TeamsTable: React.FunctionComponent = () => {
    // In real usage, this data would come from some external source like an API via props.
    const columns = {
        team_name: 'Team',
        description: 'Description',
        jira_keys: "Jira Projects"
    }

    const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false);
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [toUpdateTeam, setToUpdateTeam] = useState<ITeam>();
    const [toDeleteTeam, setToDeleteTeam] = useState<ITeam>();
    const [newTeamName, setNewTeamName] = useState<string>("");
    const [newTeamDesc, setNewTeamDesc] = useState<string>("");
    const [jiraProjects, setJiraProjects] = useState<Array<string>>([])
    const [query, setQuery] = useState<string>("")
    const [isJqlQueryValid, setIsJqlQueryValid] = useState<validate>('success')
    const [alerts, setAlerts] = React.useState<AlertInfo[]>([]);
    const [isAlertModalOpen, setIsAlertModalOpen] = React.useState(false);
    const { store } = React.useContext(ReactReduxContext);
    const state = store.getState();
    const redux_dispatch = store.dispatch;
    let currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);

    const editTeam = (team: ITeam) => {
        setIsUpdateModalOpen(true)
        setToUpdateTeam(team)
        setNewTeamName(team.team_name)
        setNewTeamDesc(team.description)
        setJiraProjects(team.jira_keys?.split(","))
    }

    const deleteTeam = (team: ITeam) => {
        setIsDeleteModalOpen(true)
        setToDeleteTeam(team)
    }

    const defaultActions = (team: ITeam): IAction[] => [
        {
            title: 'Delete Team',
            onClick: () => deleteTeam(team)
        },
        {
            title: 'Edit Team',
            onClick: () => editTeam(team)
        },
    ];

    const onUpdateSubmit = async () => {
        if (toUpdateTeam != undefined) {
            try {
                const data = {
                    team_name: newTeamName,
                    description: newTeamDesc,
                    target: toUpdateTeam.team_name,
                    jira_keys: jiraProjects.join(","),
                    jira_config: generateJiraConfig(query),
                }
                updateTeam(data)
                clear();
                setAlerts(prevAlertInfo => [...prevAlertInfo, {
                    title: 'Your changes will be updated in the background. You will need to wait a few minutes until the update is finished.',
                    variant: AlertVariant.info,
                    key: "all-created"
                }]);
                setIsAlertModalOpen(!isAlertModalOpen);
            }
            catch (error) {
                console.log(error)
            }
        }
    }

    const onDeleteSubmit = async () => {
        if (toDeleteTeam != undefined) {
            const data = {
                team_name: toDeleteTeam.team_name,
                team_description: toDeleteTeam.description,
            };
            try {
                await deleteInApi(data, '/api/quality/teams/delete');
                clear()

                if (state.auth.USER_CONFIG != "" && state.auth.USER_CONFIG != undefined) {
                    let userConfig = JSON.parse(state.auth.USER_CONFIG) as UserConfig;
                    const userConfigDefaultTeam = userConfig.teams_configuration.default_team

                    // check that the team to delete is the default one
                    if (toDeleteTeam.team_name == userConfigDefaultTeam) {
                        userConfig.teams_configuration.default_team = "n/a"
                        // update user config
                        const config = JSON.stringify(userConfig)
                        const userClaims = JSON.parse(window.atob(state.auth.IDT.split('.')[1]))

                        await createUser(userClaims.email, config)
                        redux_dispatch({ type: "SET_USER_CONFIG", data: config });
                    }
                }
                window.location.reload();
            } catch (error) {
                console.log(error);
            }
        }
    }

    const clear = () => {
        setIsUpdateModalOpen(false)
        setIsDeleteModalOpen(false)
        setToUpdateTeam(undefined)
        setToDeleteTeam(undefined)
        setNewTeamName("")
        setNewTeamDesc("")
    }

    const handleAlertModalToggle = () => {
        setIsAlertModalOpen(!setIsAlertModalOpen);
        window.location.reload();
    };


    const onJiraProjectsSelected = (options: Array<string>, query: string, isJqlQueryValid: validate) => {
        setJiraProjects(options)
        setQuery(query)
        setIsJqlQueryValid(isJqlQueryValid)
    }

    return (
        <React.Fragment>
            <TableComposable aria-label="Actions table">
                <Thead>
                    <Tr>
                        <Th>{columns.team_name}</Th>
                        <Th>{columns.description}</Th>
                        <Th>{columns.jira_keys}</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {currentTeamsAvailable.map(team => {
                        const rowActions: IAction[] | null = defaultActions(team);
                        return (
                            <Tr key={team.team_name}>
                                <Td>{team.team_name}</Td>
                                <Td>{team.description}</Td>
                                <Td>{team.jira_keys}</Td>
                                <Td isActionCell>
                                    {rowActions ? (
                                        <ActionsColumn
                                            items={rowActions}
                                        />
                                    ) : null}
                                </Td>
                            </Tr>
                        );
                    })}
                </Tbody>
            </TableComposable>
            {isUpdateModalOpen &&
                <Modal
                    variant={ModalVariant.medium}
                    title={"Update team " + toUpdateTeam?.team_name}
                    isOpen={isUpdateModalOpen}
                    onClose={clear}
                    actions={[
                        <Button key="update" variant="primary" form="modal-with-form-form" isDisabled={isJqlQueryValid == "error" || query == ""} onClick={onUpdateSubmit}>
                            Update
                        </Button>,
                        <Button key="cancel" variant="link" onClick={clear}>
                            Cancel
                        </Button>
                    ]}
                >
                    <Form>
                        <FormGroup label="Team Name" isRequired fieldId="team-name" helperText="Update your team name">
                            <TextInput value={newTeamName} type="text" onChange={(value) => { setNewTeamName(value) }} aria-label="text input example" placeholder="Update your team name" />
                        </FormGroup>
                        <FormGroup label="Description" isRequired fieldId='team-description' helperText="Update your team description">
                            <TextArea value={newTeamDesc} type="text" onChange={(value) => { setNewTeamDesc(value) }} aria-label="text area example" placeholder="Update your team description" />
                        </FormGroup>
                        <FormGroup label="Jira Projects" fieldId='jira-projects'>
                            <JiraProjects onChange={onJiraProjectsSelected} teamJiraKeys={toUpdateTeam?.jira_keys} teamName={toUpdateTeam?.team_name}></JiraProjects>
                        </FormGroup>
                    </Form>
                    <div style={{ marginTop: "2em" }}>
                        <AlertGroup isLiveRegion aria-live="polite" aria-relevant="additions text" aria-atomic="false">
                            {alerts.map(({ title, variant, key }) => (
                                <Alert variant={variant} isInline isPlain title={title} key={key} />
                            ))}
                        </AlertGroup>
                    </div>
                </Modal>
            }
            <Modal
                variant={ModalVariant.large}
                isOpen={isAlertModalOpen}
                aria-label="No header/footer modal"
                aria-describedby="modal-no-header-description"
                onClose={handleAlertModalToggle}
            >
                <div style={{ marginTop: "1em" }}>
                    <AlertGroup>
                        {
                            alerts.map(({ title, variant, key }) => (
                                <Alert variant={variant} isInline isPlain title={title} key={key} />
                            ))
                        }
                    </AlertGroup>
                </div>
            </Modal>
            {isDeleteModalOpen && <Modal
                variant={ModalVariant.small}
                title={"Delete team " + toDeleteTeam?.team_name}
                description={"All resources related to team " + toDeleteTeam?.team_name + " will be deleted. The resources cannot be recovered. Please, be sure you want to delete it."}
                isOpen={isDeleteModalOpen}
                onClose={clear}
                actions={[
                    <Button key="delete" variant="primary" form="modal-with-form-form" onClick={onDeleteSubmit}>
                        Confirm
                    </Button>,
                    <Button key="cancel" variant="link" onClick={clear}>
                        Cancel
                    </Button>
                ]}
            >
            </Modal>
            }
        </React.Fragment>
    );
};
