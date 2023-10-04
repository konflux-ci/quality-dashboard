import React from 'react';
import { TableComposable, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { fillPopOver } from '@app/Github/Table';

export interface Workflows {
    workflow_name: string;
    badge_url: string;
    job: string;
    html_url: string;
}

const columnNames = {
    jira_key: 'Name',
    days_without_priority: 'Days Without Priority',
    days_without_assignee: 'Days Without Assignee',
    days_without_resolution: 'Days Without Resolution',
};


const getValue = (value, color) => {
    if (color == "yellow"){
        color = "#FDDA0D"
    }
    return <div style={{ color: color }}>{value}</div>
}

export const BugSLOTable = (props) => {
    const bugSLOs = props.bugSLOs

    if (bugSLOs != undefined) {
        bugSLOs.sort((a, b) => (a.days_without_resolution > b.days_without_resolution ? -1 : 1));
    }

    const getInfo = (label) => {
        if (label == "Days Without Priority") {
            return fillPopOver("Days Without Priority", "priority should not be undefined for more than 2 days on untriaged bugs.")
        }
        if (label == "Days Without Assignee") {
            return fillPopOver("Days Without Assignee", "Assignee should not be undefined for more than 2 days in Blocker or Critical bugs.")
        }
        if (label == "Days Without Resolution") {
            return fillPopOver("Days Without Resolution", <div>Blocker bugs should not take more than 10 days to be resolved.<div>Critical bugs should not take more than 20 days to be resolved.</div><div>Major bugs should not take more than 40 days to be resolved.</div></div>)
        }
        return
    }

    return (
        <TableComposable aria-label="Actions table" style={{ padding: '10px' }}>
            <Thead>
                <Tr>
                    <Th>{columnNames.jira_key}</Th>
                    <Th info={getInfo(columnNames.days_without_priority)}>{columnNames.days_without_priority}</Th>
                    <Th info={getInfo(columnNames.days_without_assignee)}>{columnNames.days_without_assignee}</Th>
                    <Th info={getInfo(columnNames.days_without_resolution)}>{columnNames.days_without_resolution}</Th>
                </Tr>
            </Thead>
            <Tbody>
                {bugSLOs?.map((bugSLO) => {
                    return (
                        <Tr key={bugSLO.jira_key || ''}>
                            <Td dataLabel={columnNames.jira_key}>
                                <div>
                                    <a href={"https://issues.redhat.com/browse/" + bugSLO.jira_key} target="blank" rel="noopener noreferrer">{bugSLO.jira_key}</a>
                                </div>
                            </Td>
                            <Td dataLabel={columnNames.days_without_priority}>{bugSLO.days_without_priority != -1 ? getValue(bugSLO.days_without_priority, bugSLO.triage_sli.signal) : "-"}</Td>
                            <Td dataLabel={columnNames.days_without_assignee}>{bugSLO.days_without_assignee != -1 ? getValue(bugSLO.days_without_assignee, bugSLO.response_sli.signal) : "-"}</Td>
                            <Td dataLabel={columnNames.days_without_resolution}>{getValue(bugSLO.days_without_resolution, bugSLO.resolution_sli.signal)}</Td>
                        </Tr>
                    );
                })}
            </Tbody>
        </TableComposable>
    );
};
