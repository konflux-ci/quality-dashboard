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
    ActionsColumn,
    IAction,
} from '@patternfly/react-table';
import { useHistory } from 'react-router-dom';
import { deleteInApi } from '@app/utils/APIService';
import { FailureInfo } from './CiFailures';


export const ComposableTable: React.FC<{ failures: any, modal: any }> = ({ failures, modal }) => {
    const [failuresPage, setFailuresPage] = useState<Array<FailureInfo>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);
    const [count, setCount] = useState(failures.length);
    const [filters, setFilters] = useState({});
    const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
    const history = useHistory();
    const params = new URLSearchParams(window.location.search);


    // const defaultModalContext = useMetricsModalContextState();
    // const modalContext = useMetricsModalContext();

    async function deleteFailure(failure) {
       try {
            await deleteInApi(failure, '/api/quality/failures/delete');
            window.location.reload();
        } catch (error) {
            console.log(error);
        }
    }

    async function updateFailure(failure) {
        try {
            modal.handleModalToggle(true, failure);
        } catch (error) {
            console.log(error);
        }
    }

    const defaultActions = (failure): IAction[] => [
        {
            title: 'Delete',
            onClick: () => deleteFailure(failure),
        },
        {
            title: 'Update',
            onClick: () => updateFailure(failure),
        },
    ];

    useEffect(() => {
        if (failures.length == 0) {
            setPage(1)
            setFailuresPage([])
        }
        if (failures.length > 0) {
            setFailuresPage(failures.slice(0, perPage))
            setPage(1)
        }
    }, [failures]);

    const columnNames = {
        jira_key: "Jira Key",
        title_from_jira: "Title from Jira",
        jira_status: "Jira Status",
        error_message: "Error Message",
        frequency: "Frequency",
        created_date: "Created Date",
        closed_date: "Closed Date",
        labels: "Labels",
    };

    useEffect(
        () => {
            setCount(failures.length);
        },
        [failures],
    );

    useEffect(() => {
        setCount(failures.length)
        if (failures.length > 0) {
            const filteredRows = filterRows(failures, filters)
            const sortedRows = sortRows(filteredRows)
            setCount(sortedRows.length)

            const from = (page - 1) * perPage
            const to = (page - 1) * perPage + perPage >= sortedRows.length ? sortedRows.length : (page - 1) * perPage + perPage;

            setFailuresPage(sortedRows.slice(from, to))
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
        { column: 'jira_key', label: 'Jira Key' },
        { column: 'title_from_jira', label: 'Title from Jira'},
        { column: 'jira_status', label: 'Jira Status' },
        { column: 'error_message', label: 'Error Message' },
        { column: 'frequency', label: 'Frequency' },
        { column: 'created_date', label: 'Created Date'},
        { column: 'closed_date', label: 'Closed Date'},
        { column: 'labels', label: 'Labels'},

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
    const getSortableRowValues = (failure: FailureInfo): (string | number)[] => {
        const { jira_key, jira_status, error_message, frequency } = failure;
        return [jira_key, jira_status, error_message, frequency];
    };

    const sortRows = (rows) => {
        if (activeSortIndex !== null) {
            return rows.sort((a, b) => {
                const aValue = getSortableRowValues(a)[activeSortIndex] || a.frequency == 0  ? getSortableRowValues(a)[activeSortIndex] : "-";
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

                {/* <GetMetrics></GetMetrics> */}

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
                        {failuresPage.map((failure, index) => {
                            const rowActions: IAction[] | null = defaultActions(failure);

                            return (
                                <Tr key={index} {...(index % 2 === 0 && { isStriped: true })}>
                                    <Td dataLabel={columnNames.jira_key}>
                                        <div>
                                            <a href={"https://issues.redhat.com/browse/"+failure.jira_key} target="blank" rel="noopener noreferrer">{failure.jira_key}</a>
                                        </div>
                                    </Td>
                                    <Td dataLabel={columnNames.title_from_jira}>{failure.title_from_jira}</Td>
                                    <Td dataLabel={columnNames.jira_status}>{failure.jira_status}</Td>
                                    <Td dataLabel={columnNames.error_message}>{failure.error_message}</Td>
                                    <Td dataLabel={columnNames.frequency}>{failure.frequency+"%"}</Td>
                                    <Td dataLabel={columnNames.created_date}>{failure.created_date}</Td>
                                    <Td dataLabel={columnNames.closed_date}>{failure.closed_date}</Td>
                                    <Td dataLabel={columnNames.closed_date}>{failure.labels}</Td>
                                    <Td isActionCell>
                                        {rowActions ? (
                                            <ActionsColumn
                                                items={rowActions}
                                            />
                                        ) : null}
                                    </Td>
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