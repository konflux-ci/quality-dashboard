import React from 'react';
import { Caption, TableComposable, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { ExternalLinkAltIcon } from '@patternfly/react-icons';


interface Workflows {
    workflow_name: string;
    badge_url: string;
    job: string;
    html_url: string
};

const columnNames = {
    name: 'Name',
    html_url: 'Job URL',
    job: 'Job State',
    badge: 'Last Execution'
};


export const GithubActions = (props) => {
    return (
        <TableComposable aria-label="Actions table" style={{ padding: "10px" }}>
            <Caption>All GitHub Actions available in the repository {props.repoName}</Caption>
            <Thead>
                <Tr>
                    <Th>{columnNames.name}</Th>
                    <Th>{columnNames.job}</Th>
                    <Th>{columnNames.badge}</Th>
                    <Th>{columnNames.html_url}</Th>
                </Tr>
            </Thead>
            <Tbody>
                {props.workflows.map(workf => {
                    const workflow: Workflows = workf
                    return (
                        <Tr key={workflow.workflow_name || ""}>
                            <Td dataLabel={columnNames.name}>{workflow.workflow_name || ""}</Td>
                            <Td dataLabel={columnNames.job}>{workflow.job || ""}</Td>
                            <Td dataLabel={columnNames.badge}><img src={workflow.badge_url}></img></Td>
                            <Td dataLabel={columnNames.html_url}><a href={workflow.html_url} target="blank" rel="noopener noreferrer"><ExternalLinkAltIcon>Link</ExternalLinkAltIcon>Go to job</a></Td>
                        </Tr>
                    );
                })}
            </Tbody>
        </TableComposable>
    );
}