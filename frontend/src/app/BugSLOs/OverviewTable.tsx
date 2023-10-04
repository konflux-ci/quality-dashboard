import React, { useState, useEffect } from 'react';
import { Button, Pagination } from '@patternfly/react-core';
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    ThProps,
} from '@patternfly/react-table';
import { fillPopOver } from '@app/Github/Table';
import { GetMetrics, MetricsModalContext, useMetricsModalContext, useMetricsModalContextState } from './Metrics';

interface BugSLOInfo {
    total: number,
    average: number,
}

interface Alert {
    alert_message: string,
    signal: string,
}

interface BugSLO {
    jira_key: string,
    jira_url: string,
    triage_sli: Alert,
    response_sli: Alert,
    resolution_sli: Alert,
    days_without_assignee: number,
    days_without_priority: number,
    days_without_resolution: number,
}

interface ProjectInfo {
    project_key: string;
    bug_slos: BugSLO[];
    red_triage_time_bug_slo_info: BugSLOInfo;
    yellow_triage_time_bug_slo_info: BugSLOInfo;
    red_response_time_bug_slo_info: BugSLOInfo;
    red_resolution_time_bug_slo_info: BugSLOInfo;
    yellow_resolution_time_bug_slo_info: BugSLOInfo;

}

export const OverviewTable: React.FC<{ bugSLOs: any }> = ({ bugSLOs }) => {
    const [bugSLOsPage, setBugSLOsPage] = useState<Array<ProjectInfo>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);
    const [count, setCount] = useState(bugSLOs.length);
    const [filters, setFilters] = useState({});
    const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
    const defaultModalContext = useMetricsModalContextState();
    const modalContext = useMetricsModalContext();


    useEffect(() => {
        if (bugSLOs.length == 0) {
            setPage(1)
            setBugSLOsPage([])
        }
        if (bugSLOs.length > 0) {
            setBugSLOsPage(bugSLOs.slice(0, perPage))
            setPage(1)
        }
    }, [bugSLOs]);

    const columnNames = {
        project_key: "Project Key",
        red_triage_time_bug_slo_info: "Total Bugs Meeting Triage Time Bug SLO",
        yellow_triage_time_bug_slo_info: "Total Bugs At Risk of Meeting Triage Time Bug SLO",
        red_response_time_bug_slo_info: "Total Bugs Meeting Response Time Bug SLO",
        red_resolution_time_bug_slo_info: "Total Bugs Meeting Resolution Time Bug SLO",
        yellow_resolution_time_bug_slo_info: "Total Bugs At Risk of Meeting Resolution Time Bug SLO",

    };

    useEffect(
        () => {
            setCount(bugSLOs.length);
        },
        [bugSLOs],
    );

    useEffect(() => {
        setCount(bugSLOs.length)
        if (bugSLOs.length > 0) {
            const filteredRows = filterRows(bugSLOs, filters)
            const sortedRows = sortRows(filteredRows)
            setCount(sortedRows.length)

            const from = (page - 1) * perPage
            const to = (page - 1) * perPage + perPage >= sortedRows.length ? sortedRows.length : (page - 1) * perPage + perPage;

            setBugSLOsPage(sortedRows.slice(from, to))
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
        { column: 'project_key', label: 'Project Key' },
        { column: 'red_triage_time_bug_slo_info', label: 'Total Bugs Meeting Triage Time Bug SLO' },
        { column: 'yellow_triage_time_bug_slo_info', label: 'Total Bugs At Risk of Meeting Triage Time Bug SLO' },
        { column: 'red_response_time_bug_slo_info', label: 'Total Bugs Meeting Response Time Bug SLO' },
        { column: 'red_resolution_time_bug_slo_info', label: 'Total Bugs Meeting Resolution Time Bug SLO' },
        { column: 'yellow_resolution_time_bug_slo_info', label: 'Total Bugs At Risk of Meeting Resolution Time Bug SLO' },

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
    const getSortableRowValues = (project: any): (string | number | BugSLOInfo)[] => {
        const { project_key, red_triage_time_bug_slo_info, yellow_triage_time_bug_slo_info, red_response_time_bug_slo_info, red_resolution_time_bug_slo_info, yellow_resolution_time_bug_slo_info } = project;
        return [project_key, red_triage_time_bug_slo_info, yellow_triage_time_bug_slo_info, red_response_time_bug_slo_info, red_resolution_time_bug_slo_info, yellow_resolution_time_bug_slo_info];
    };

    const sortRows = (rows) => {
        if (activeSortIndex !== null) {
            return rows.sort((a, b) => {
                let aValue = getSortableRowValues(a)[activeSortIndex] || a.frequency == 0 ? getSortableRowValues(a)[activeSortIndex] : "-";
                let bValue = getSortableRowValues(b)[activeSortIndex] || b.frequency == 0 ? getSortableRowValues(b)[activeSortIndex] : "-";


                // workaround for sorting the following fields:
                // red_triage_time_bug_slo_info, red_response_time_bug_slo_info, and red_resolution_time_bug_slo_info
                if (typeof getSortableRowValues(a)[activeSortIndex] == 'object') {
                    const targetA = getSortableRowValues(a)[activeSortIndex] as BugSLOInfo
                    const targetB = getSortableRowValues(b)[activeSortIndex] as BugSLOInfo

                    aValue = targetA.total || a.frequency == 0 ? targetA.total : "-";
                    bValue = targetB.total || b.frequency == 0 ? targetB.total : "-";

                }

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

    const getInfo = (label) => {
        if (label == "Total Bugs Meeting Triage Time Bug SLO") {
            return fillPopOver(label, "Number of bugs that meet Triage Time Bug SLO (priority should not be undefined for more than 2 days on untriaged bugs).")
        }
        if (label == "Total Bugs Meeting Response Time Bug SLO") {
            return fillPopOver(label, "Number of bugs that meet Response Time Bug SLO (assignee should not be undefined for more than 2 days in Blocker or Critical bugs).")
        }
        if (label == "Total Bugs Meeting Resolution Time Bug SLO") {
            return fillPopOver(label, <div><div>Number of bugs that meet Red Resolution Time Bug SLO.</div>Blocker bugs should not take more than 10 days to be resolved.<div>Critical bugs should not take more than 20 days to be resolved.</div><div>Major bugs should not take more than 40 days to be resolved.</div></div>)
        }
        return
    }

    return (
        <MetricsModalContext.Provider value={defaultModalContext}>
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

                <GetMetrics></GetMetrics>

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
                        {bugSLOsPage.map((bugSLO, index) => {
                            return (
                                <Tr key={index} {...(index % 2 === 0 && { isStriped: true })}>
                                    <Td dataLabel={columnNames.project_key}>
                                        <div>
                                            {bugSLO.project_key}
                                        </div>
                                        <div style={{ marginTop: 2 }}>
                                            <Button style={{ fontSize: 14 }} variant="link" onClick={() => modalContext.handleModalToggle(bugSLO)}>{"> Show detailed metrics"}</Button>
                                        </div>
                                    </Td>
                                    <Td dataLabel={columnNames.red_triage_time_bug_slo_info}>{bugSLO.red_triage_time_bug_slo_info.total}</Td>
                                    <Td dataLabel={columnNames.yellow_triage_time_bug_slo_info}>{bugSLO.yellow_triage_time_bug_slo_info.total}</Td>
                                    <Td dataLabel={columnNames.red_response_time_bug_slo_info}>{bugSLO.red_response_time_bug_slo_info.total}</Td>
                                    <Td dataLabel={columnNames.red_resolution_time_bug_slo_info}>{bugSLO.red_resolution_time_bug_slo_info.total}</Td>
                                    <Td dataLabel={columnNames.yellow_resolution_time_bug_slo_info}>{bugSLO.yellow_resolution_time_bug_slo_info.total}</Td>

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
        </MetricsModalContext.Provider >
    );
};