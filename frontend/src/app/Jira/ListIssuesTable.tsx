import React, { useEffect, useState } from 'react';
import { Table, Thead, Tr, Th, Tbody, Td, ThProps } from '@patternfly/react-table';
import { Pagination } from '@patternfly/react-core';

export interface Issue {
    jira_key: string;
    summary: string;
    status: string;
    priority: string;
    labels: string;
    component: string;
    age: string;
    assignee: string;
}

const columnNames = {
    jira_key: 'Key',
    summary: 'Summary',
    status: 'Status',
    priority: 'Priority',
    labels: 'Labels',
    component: 'Component',
    age: 'Age',
    assignee: 'Assignee'
};


export const ListIssues: React.FC<{ issues: any }> = ({ issues }) => {
    const [issuesPage, setIssuesPage] = useState<Array<Issue>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);
    const [count, setCount] = useState(issues?.length);
    const [filters, setFilters] = useState({});
    const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);

    useEffect(() => {
        if (issues.length == 0) {
            setPage(1)
            setIssuesPage([])
        }
        if (issues.length > 0) {
            setIssuesPage(issues.slice(0, perPage))
            setPage(1)
        }
    }, [issues]);

    useEffect(
        () => {
            setCount(issues.length);
        },
        [issues],
    );

    useEffect(() => {
        setCount(issues.length)
        if (issues.length > 0) {
            const filteredRows = filterRows(issues, filters)
            const sortedRows = sortRows(filteredRows)
            setCount(sortedRows.length)

            const from = (page - 1) * perPage
            const to = (page - 1) * perPage + perPage >= sortedRows.length ? sortedRows.length : (page - 1) * perPage + perPage;

            setIssuesPage(sortedRows.slice(from, to))
        }
    }, [page, perPage, filters, activeSortIndex, activeSortDirection]);

    const onSetPage = (_event: React.MouseEvent | React.KeyboardEvent | MouseEvent, newPage: number) => {
        setPage(newPage);
    };

    const onPerPageSelect = (
        _event: React.MouseEvent | React.KeyboardEvent | MouseEvent,
        newPerPage: number,
        newPage: number
    ) => {
        setPerPage(newPerPage);
        setPage(newPage);
    };

    // Filters helpers
    const columns = [
        { column: 'key', label: 'Key' },
        { column: 'summary', label: 'Summary' },
        { column: 'status', label: 'Status' },
        { column: 'priority', label: 'Priority' },
        { column: 'labels', label: 'Labels' },
        { column: 'component', label: 'Components' },
        { column: 'age', label: 'Age' },
        { column: 'assignee', label: 'Assignee' }


    ]

    function filterRows(rows, filters) {
        if (Object.keys(filters).length === 0) return rows

        return rows.filter(row => {
            return Object.keys(filters).every(column => {
                const value = row[column]
                const searchValue = filters[column]

                // handle Resolution Time filter
                if (typeof value === "number") {
                    return value == searchValue
                }

                if (typeof value === 'string') {
                    return value.toLocaleLowerCase().includes(searchValue.toLocaleLowerCase())
                }
                return false
            })
        })
    }

    const handleSearch = (value, column) => {
        if (value) {
            setFilters(prevFilters => ({
                ...prevFilters,
                [column]: value,
            }))
        } else {
            setFilters(prevFilters => {
                const updatedFilters = { ...prevFilters }
                delete updatedFilters[column]

                return updatedFilters
            })
        }
    }
    // End of filter helpers

    // Sort helpers
    const getSortableRowValues = (issue: Issue): (string | number)[] => {
        const { jira_key, summary, status, priority, labels, component, age, assignee } = issue;
        return [jira_key, summary, status, priority, labels, component, age, assignee];
    };

    const sortRows = (rows) => {
        if (activeSortIndex !== null) {
            return rows.sort((a, b) => {
                const aValue = getSortableRowValues(a)[activeSortIndex] || a.frequency == 0 ? getSortableRowValues(a)[activeSortIndex] : "-";
                const bValue = getSortableRowValues(b)[activeSortIndex] || b.frequency == 0 ? getSortableRowValues(b)[activeSortIndex] : "-";

                if (typeof aValue === 'number') {
                    // Numeric sort
                    if (activeSortDirection === 'asc') {
                        return (aValue as number) - (bValue as number);
                    }
                    return (bValue as number) - (aValue as number);
                } else if (typeof aValue === 'string') {
                    // String sort
                    if (activeSortDirection === 'asc') {
                        return aValue.localeCompare(bValue.toString(), undefined, { numeric: true, sensitivity: 'base' });
                    }
                    return bValue.toString().localeCompare(aValue as string, undefined, { numeric: true, sensitivity: 'base' });
                }
            });
        }
        return rows
    }

    const getSortParams = (columnIndex: number): ThProps['sort'] => ({
        sortBy: {
            index: activeSortIndex as number,
            direction: activeSortDirection as any
        },
        onSort: (_event, index, direction) => {
            setActiveSortIndex(index);
            setActiveSortDirection(direction);
        },
        columnIndex
    });
    // End of sort helpers
    
    return (
        <div>
            <Pagination
                perPageComponent="button"
                itemCount={count}
                perPage={perPage}
                page={page}
                onSetPage={onSetPage}
                widgetId="top-example"
                onPerPageSelect={onPerPageSelect}
            />
            <Table aria-label="Controlling text">
                <Thead>
                    <Tr>
                        {columns.map((column, idx) => {
                            return (
                                <Th sort={getSortParams(idx)} key={idx} modifier='wrap'>
                                    {column.label}
                                </Th>
                            )
                        })}
                    </Tr>
                    <Tr>
                        {columns.map(c => {
                            return (
                                <Th key={c.column}>
                                    <input
                                        style={{ width: "70px" }}
                                        key={`${c.column}-search`}
                                        type="search"
                                        placeholder={`Search`}
                                        value={filters[c.column]}
                                        onChange={event => handleSearch(event.target.value, c.column)}
                                    />
                                </Th>
                            )
                        })}
                    </Tr>
                </Thead>
                <Tbody>
                    {issuesPage.map((issue, index) => {
                        return (
                            <Tr key={index}>
                                <Td dataLabel={columnNames.jira_key}><a href={"https://issues.redhat.com/browse/" + issue.jira_key} target={"https://issues.redhat.com/browse/" + issue.jira_key}>{issue.jira_key}</a></Td>
                                <Td dataLabel={columnNames.summary}>{issue.summary}</Td>
                                <Td dataLabel={columnNames.status}>{issue.status}</Td>
                                <Td dataLabel={columnNames.priority}>{issue.priority}</Td>
                                <Td dataLabel={columnNames.priority}>{issue.labels}</Td>
                                <Td dataLabel={columnNames.component}>{issue.component == "undefined" ? "-" : issue.component}</Td>
                                <Td dataLabel={columnNames.age} modifier="nowrap">{issue.age +" days"}</Td>
                                <Td dataLabel={columnNames.assignee} modifier="nowrap">{issue.assignee}</Td>
                            </Tr>
                        )
                    })}
                </Tbody>
            </Table>
            <Pagination
                perPageComponent="button"
                itemCount={count}
                perPage={perPage}
                page={page}
                onSetPage={onSetPage}
                widgetId="top-example"
                onPerPageSelect={onPerPageSelect}
            />
        </div>
    )
};
