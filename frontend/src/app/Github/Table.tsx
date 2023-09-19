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
    ActionsColumn,
    IAction,
} from '@patternfly/react-table';
import { GetCodeCovInfo } from './CodeCov';
import { deleteInApi } from '@app/utils/APIService';
import { useHistory } from 'react-router-dom';
import { RepositoryInfo } from './Github';
import { GetMetrics, MetricsModalContext, useMetricsModalContext, useMetricsModalContextState } from './Metrics';


export const ComposableTable: React.FC<{ repos: any, modal: any }> = ({ repos, modal }) => {
    const [reposPage, setReposPage] = useState<Array<RepositoryInfo>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);
    const [count, setCount] = useState(repos.length);
    const [filters, setFilters] = useState({});
    const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
    const history = useHistory();
    const params = new URLSearchParams(window.location.search);
    const defaultModalContext = useMetricsModalContextState();
    const modalContext = useMetricsModalContext();

    async function deleteRepository(gitOrg: string, repoName: string) {
        const data = {
            git_organization: gitOrg,
            repository_name: repoName,
        };
        try {
            await deleteInApi(data, '/api/quality/repositories/delete');
            history.push(window.location.pathname + '?' + 'team=' + params.get('team'));
            window.location.reload();
        } catch (error) {
            console.log(error);
        }
    }

    async function editRepository(repo) {
        try {
            modal.handleModalToggle(true, repo);
        } catch (error) {
            console.log(error);
        }
    }

    const defaultActions = (repo): IAction[] => [
        {
            title: 'Delete Repository',
            onClick: () => deleteRepository(repo.git_organization, repo.repository_name),
        },
        {
            title: 'Edit Repository',
            onClick: () => editRepository(repo),
        },
    ];

    useEffect(() => {
        if (repos.length == 0) {
            setPage(1)
            setReposPage([])
        }
        if (repos.length > 0) {
            setReposPage(repos.slice(0, perPage))
            setPage(1)
        }
    }, [repos]);

    const columnNames = {
        repository_name: "Repository",
        git_organization: "Organization",
        description: "Description",
        code_cov: "CodeCov",
        retest_before_merge_avg: "Retest Before Merge Avg",
        open_prs: "Open PRs",
        merged_prs: "Merged PRs",
        time_to_merge_pr_avg_days: "Time To Merge PR Avg Days",
    };

    useEffect(
        () => {
            setCount(repos.length);
        },
        [repos],
    );

    useEffect(() => {
        setCount(repos.length)
        if (repos.length > 0) {
            const filteredRows = filterRows(repos, filters)
            const sortedRows = sortRows(filteredRows)
            setCount(sortedRows.length)

            const from = (page - 1) * perPage
            const to = (page - 1) * perPage + perPage >= sortedRows.length ? sortedRows.length : (page - 1) * perPage + perPage;

            setReposPage(sortedRows.slice(from, to))
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
        { column: 'repository_name', label: 'Repository' },
        { column: 'git_organization', label: 'Organization' },
        { column: 'description', label: 'Description' },
        { column: 'code_cov', label: 'Code Cov' },
        { column: 'retest_before_merge_avg', label: 'Retest Before Merge Avg' },
        { column: 'open_prs', label: 'Total Open PRs' },
        { column: 'merged_prs', label: 'Total Merged PRs' },
        { column: 'time_to_merge_pr_avg_days', label: 'Time To Merge PR Avg Days' },
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

                // handle ID, Summary, Created at, Updated at, and Resolved at filters
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
    const getSortableRowValues = (repo: RepositoryInfo): (string | number)[] => {
        const { repository_name, git_organization, description, code_cov, retest_before_merge_avg, open_prs, merged_prs, time_to_merge_pr_avg_days } = repo;
        return [repository_name, git_organization, description, code_cov, retest_before_merge_avg, open_prs, merged_prs, time_to_merge_pr_avg_days];
    };

    const sortRows = (rows) => {
        if (activeSortIndex !== null) {
            return rows.sort((a, b) => {
                const aValue = getSortableRowValues(a)[activeSortIndex] ? getSortableRowValues(a)[activeSortIndex] : "-";
                const bValue = getSortableRowValues(b)[activeSortIndex] ? getSortableRowValues(b)[activeSortIndex] : "-";

                if (aValue == 'N/A') {
                    return 1
                }

                if (bValue == 'N/A') {
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


    const fillPopOver = (title, description) => {
        return {
            popover: (description),
            ariaLabel: 'More information on' + title,
            popoverProps: {
                headerContent: title,
            }

        }
    }

    const getInfo = (label) => {
        if (label == "Retest Before Merge Avg") {
            return fillPopOver("Retest Before Merge Avg", "Calculates an average how many /test and /retest comments were issued after the last code push (in the selected time range)")
        }
        if (label == "Time To Merge PR Avg Days") {
            return fillPopOver("Time To Merge PR Avg Days", "Calculates an average of how many days were needed to merge a PR (difference between creation and merged date in the selected time range)")
        }
        if (label == "Code Cov") {
            return fillPopOver("Code Cov", "The coverage trend is calculated through the two last commits. No trend arrow means that the coverage trend is stable")
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
                        {reposPage.map((repo, index) => {
                            const rowActions: IAction[] | null = defaultActions(repo);

                            return (
                                <Tr key={index} {...(index % 2 === 0 && { isStriped: true })}>
                                    <Td dataLabel={columnNames.repository_name}>
                                        <div>
                                            <a href={repo.git_url} target={repo.git_url}>{repo.repository_name}</a>
                                        </div>
                                        <div style={{ marginTop: 2 }}>
                                            <Button style={{ fontSize: 14 }} variant="link" onClick={() => modalContext.handleModalToggle(repo)}>{"> Show detailed metrics"}</Button>
                                        </div>
                                    </Td>
                                    <Td dataLabel={columnNames.git_organization}>{repo.git_organization}</Td>
                                    <Td dataLabel={columnNames.description}>{repo.description}</Td>
                                    <Td dataLabel={columnNames.code_cov}>{repo.code_cov == 'N/A' ? "N/A" : GetCodeCovInfo(repo, 'left')}</Td>
                                    <Td dataLabel={columnNames.retest_before_merge_avg}>{repo.retest_before_merge_avg}</Td>
                                    <Td dataLabel={columnNames.open_prs}>{repo.open_prs}</Td>
                                    <Td dataLabel={columnNames.merged_prs}>{repo.merged_prs}</Td>
                                    <Td dataLabel={columnNames.time_to_merge_pr_avg_days}>{repo.time_to_merge_pr_avg_days}</Td>
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
        </MetricsModalContext.Provider >
    );
};