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
import { deleteInApi, updateTeam } from "@app/utils/APIService";
import { useSelector } from 'react-redux';
import { Button, Form, FormGroup, Modal, ModalVariant, TextArea, TextInput } from "@patternfly/react-core";

export const TeamsTable: React.FunctionComponent = () => {
    // In real usage, this data would come from some external source like an API via props.
    const columns = {
        team_name: 'Team',
        description: 'Description',
    }

    const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false);
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [toUpdateTeam, setToUpdateTeam] = useState<ITeam>();
    const [toDeleteTeam, setToDeleteTeam] = useState<ITeam>();
    const [newTeamName, setNewTeamName] = useState<string>("");
    const [newTeamDesc, setNewTeamDesc] = useState<string>("");

    let currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);

    const editTeam = (team: ITeam) => {
        setIsUpdateModalOpen(true)
        setToUpdateTeam(team)
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
                }
                await updateTeam(data)
                clear()
                window.location.reload();
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

    return (
        <React.Fragment>
            <TableComposable aria-label="Actions table">
                <Thead>
                    <Tr>
                        <Th>{columns.team_name}</Th>
                        <Th>{columns.description}</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {currentTeamsAvailable.map(team => {
                        console.log(toUpdateTeam)
                        const rowActions: IAction[] | null = defaultActions(team);
                        return (
                            <Tr key={team.team_name}>
                                <Td>{team.team_name}</Td>
                                <Td>{team.description}</Td>
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
            {isUpdateModalOpen && <Modal
                variant={ModalVariant.medium}
                title={"Update team " + toUpdateTeam?.team_name}
                isOpen={isUpdateModalOpen}
                onClose={clear}
                actions={[
                    <Button key="update" variant="primary" form="modal-with-form-form" onClick={onUpdateSubmit}>
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
                </Form>
            </Modal>
            }
            {isDeleteModalOpen && <Modal
                variant={ModalVariant.medium}
                title={"Delete team " + toDeleteTeam?.team_name}
                description={"All resources related to team " + toDeleteTeam?.team_name + " will be deleted. The resources cannot be recovered. Please, be sure you want to delete it."}
                isOpen={isDeleteModalOpen}
                onClose={clear}
                actions={[
                    <Button key="delete" variant="primary" form="modal-with-form-form" onClick={onDeleteSubmit}>
                        Delete
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
