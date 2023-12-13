import React, { useState, useEffect } from 'react';
import { Pagination } from '@patternfly/react-core';
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    ThProps,
} from '@patternfly/react-table';
import { Bug } from './Types';
import { fillPopOver } from '@app/Github/Table';

export const OverviewTable: React.FC<{ bugSLIs: Array<Bug>, selected: string }> = ({ bugSLIs, selected }) => {
    const [bugSLIsPage, setBugSLIsPage] = useState<Array<Bug>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);
    const [count, setCount] = useState(bugSLIs.length);
    const [filters, setFilters] = useState({});
    const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);

    useEffect(() => {
        if (bugSLIs.length == 0) {
            setPage(1)
            setBugSLIsPage([])
        }
        if (bugSLIs.length > 0) {
            setBugSLIsPage(bugSLIs.slice(0, perPage))
            setPage(1)
        }
    }, [bugSLIs]);

    const columnNames = {
        jira_key: 'Jira Key',
        summary: 'Summary',
        priority: 'Priority',
        component: 'Component',
        labels: 'Labels',
        status: 'Status',
        days_without_assignee: 'Days Without Assignment',
        days_without_priority: 'Days Without Prioritization',
        days_without_resolution: 'Days Without Resolution',
    };

    useEffect(
        () => {
            setCount(bugSLIs.length);
        },
        [bugSLIs],
    );

    useEffect(() => {
        setCount(bugSLIs.length)
        if (bugSLIs.length > 0) {
            const filteredRows = filterRows(bugSLIs, filters)
            const sortedRows = sortRows(filteredRows)
            setCount(sortedRows.length)

            const from = (page - 1) * perPage
            const to = (page - 1) * perPage + perPage >= sortedRows.length ? sortedRows.length : (page - 1) * perPage + perPage;

            setBugSLIsPage(sortedRows.slice(from, to))
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
        { column: 'jira_key', label: 'Project Key' },
        { column: 'summary', label: 'Summary' },
        { column: 'priority', label: 'Priority' },
        { column: 'component', label: 'Component' },
        { column: 'labels', label: 'Labels' },
        { column: 'status', label: 'Status' },
    ]

    if (selected == "response") {
        columns.push({ column: 'days_without_assignee', label: 'Days Without Assignment' },)
    } else if (selected == "resolution") {
        columns.push({ column: 'days_without_resolution', label: 'Days Without Resolution' },)
    } else {
        columns.push({ column: 'days_without_priority', label: 'Days Without Prioritization' },)
    }

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
    const getSortableRowValues = (bug: any): (string | number)[] => {
        if (selected == "resolution") {
            const { jira_key, summary, priority, component, labels, status, days_without_resolution } = bug;
            return [jira_key, summary, priority, component, labels, status, days_without_resolution]
        } else if (selected == "response") {
            const { jira_key, summary, priority, component, labels, status, days_without_assignee } = bug;
            return [jira_key, summary, priority, component, labels, status, days_without_assignee]
        }

        const { jira_key, summary, priority, component, labels, status, days_without_priority } = bug;
        return [jira_key, summary, priority, component, labels, status, days_without_priority];

    };

    const sortRows = (rows) => {
        if (activeSortIndex !== null) {
            return rows.sort((a, b) => {
                const aValue = getSortableRowValues(a)[activeSortIndex] || a.frequency == 0 ? getSortableRowValues(a)[activeSortIndex] : "-";
                const bValue = getSortableRowValues(b)[activeSortIndex] || b.frequency == 0 ? getSortableRowValues(b)[activeSortIndex] : "-";

                if (aValue == "-") {
                    return 1
                }

                if (bValue == "-") {
                    return -1
                }

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


    const getSLI = (bug: Bug, selected: string) => {
        switch (selected) {
            case "response":
                return bug.response_sli.signal
            case "resolution":
                return bug.resolution_sli.signal
            default:
                return bug.triage_sli.signal
        }
    }

    const getColor = (bug: Bug, selected: string) => {
        const sli = getSLI(bug, selected)

        if (sli == "red") {
            return "#FA8A8A"
        } else if (sli == "yellow") {
            return "#FAEE84"
        }
    }
    const getInfo = (label) => {
        if (label == 'Days Without Assignment' || label == 'Days Without Prioritization' || label == 'Days Without Resolution') {
            return fillPopOver(label, "Only considers working days.")
        }
        return
    }

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

            <TableComposable aria-label="Simple table" >
                <Thead>
                    <Tr>
                        {columns.map((column, idx) => {
                            return (
                                <Th
                                    modifier="wrap"
                                    width={10}
                                    sort={getSortParams(idx)}
                                    key={idx}
                                    info={getInfo(column.label)}
                                >
                                    <div>
                                        {column.label}
                                    </div>
                                </Th>
                            )
                        })}
                    </Tr>
                    <Tr>
                        {columns.map(c => {
                            return (
                                <Th key={c.column}>
                                    <input style={{ width: '100%' }}
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
                    {bugSLIsPage.map((bug, index) => {
                        return (
                            <Tr key={index} style={{ background: getColor(bug, selected) }} {...(index % 2 === 0 && { isStriped: true })} >
                                <Td dataLabel={columnNames.jira_key}>
                                    <div>
                                        <a href={"https://issues.redhat.com/browse/" + bug.jira_key} target="blank" rel="noopener noreferrer">{bug.jira_key}</a>
                                    </div>
                                </Td>
                                <Td dataLabel={columnNames.summary}>{bug.summary}</Td>
                                <Td dataLabel={columnNames.priority}>{bug.priority}</Td>
                                <Td dataLabel={columnNames.component}>{bug.component}</Td>
                                <Td dataLabel={columnNames.labels}>{bug.labels}</Td>
                                <Td dataLabel={columnNames.status}>{bug.status}</Td>
                                {selected == "response" && <Td dataLabel={columnNames.days_without_assignee}>{bug.days_without_assignee}</Td>}
                                {selected == "resolution" && <Td dataLabel={columnNames.days_without_resolution}>{bug.days_without_resolution}</Td>}
                                {selected == "triage" && <Td dataLabel={columnNames.days_without_priority}>{bug.days_without_priority}</Td>}
                            </Tr>
                        )
                    })}
                </Tbody>
            </TableComposable>

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
    );
};