import React from "react";
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
import { deleteTeam } from "@app/utils/APIService";
import { useSelector } from 'react-redux';


export const TableComponent: React.FunctionComponent = () => {
    // In real usage, this data would come from some external source like an API via props.
    const columns = {
        team_name: 'Team',
        description: 'Description',
    }
    
    let currentTeamsAvailable = useSelector((state: any) => state.teams.TeamsAvailable);


    const defaultActions = (team: ITeam): IAction[] => [
        {
            title: 'Delete Team',
            onClick: () => deleteTeam(team.team_name, team.description)
        },
        // {
        //     title: 'Edit Team',
        //     onClick: () => editTeam(team)
        // },
    ];

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
                            console.log(team)
                            const rowActions: IAction[] | null = defaultActions(team);
                            return (
                                <Tr key={team.team_name}>
                                    <Td>{team.team_name}</Td>
                                    <Td>{team.description}</Td>
                                    <Td>
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
            </React.Fragment>
    );
};
