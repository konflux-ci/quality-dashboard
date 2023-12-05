import React from 'react';
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    InnerScrollContainer,
    ExpandableRowContent,
    ThProps
} from '@patternfly/react-table';
import { Grid, GridItem, } from '@patternfly/react-core';
import { Card, CardTitle, CardBody } from '@patternfly/react-core';
import { Flakey, TestCase } from './Types';

export const InnerNestedComposableTableNestedExpandable: React.FunctionComponent<{ test_cases: TestCase[], rowIndex: number }> = ({ test_cases, rowIndex }) => {
    const [simpleSortIndex, setSimpleSortIndex] = React.useState<number>(0);
    const [simpleSortDirection, setSimpleSortDirection] = React.useState<'asc' | 'desc'>('desc');
    const onSimpleSort = (a, b) => {
        const aValue = getSimpleSortableRowValues(a)[simpleSortIndex];
        const bValue = getSimpleSortableRowValues(b)[simpleSortIndex];

        if (simpleSortDirection === 'asc') {
            return (aValue as number) - (bValue as number);
        }
        return (bValue as number) - (aValue as number);
    }

    const getSimpleSortableRowValues = (flakey: any): (number)[] => {
        const { count, test_case_impact } = flakey;
        return [count, test_case_impact];
    };

    // Expandable test cases
    const [expandedTestCaseNames, setExpandedTestCaseNames] = React.useState<string[]>([]);

    const setTestCaseExpanded = (suite_name: string, isExpanding = true) => {
        setExpandedTestCaseNames(prevExpanded => {
            const otherExpandedTestCaseNames = prevExpanded.filter(t => t !== suite_name);
            return isExpanding ? [...otherExpandedTestCaseNames, suite_name] : otherExpandedTestCaseNames;
        });
    }
    const isTestCaseExpanded = (test_case_name: string) => expandedTestCaseNames.includes(test_case_name);

    const expandPre = (e) => {
        if (e.currentTarget.classList.contains('expandedPre')) {
            e.currentTarget.classList.remove("expandedPre")
            e.currentTarget.classList.add('collapsedPre')

        } else {
            e.currentTarget.classList.remove('collapsedPre')
            e.currentTarget.classList.add('expandedPre')
        }
    }

    const getSimpleSortParams = (columnIndex: number): ThProps['sort'] => ({
        sortBy: {
            index: simpleSortIndex,
            direction: simpleSortDirection,
            defaultDirection: 'desc'
        },
        onSort: (_event, index, direction) => {
            setSimpleSortIndex(index);
            setSimpleSortDirection(direction);
        },
        columnIndex
    });

    const columnNames = {
        name: 'Test Case',
        status: 'Status',
        error_message: 'Failing Jobs',
        count: 'Count',
        suite_name: 'Suite Name',
        job_id: "Job ID",
        job_url: "Job Url",
        failure_date: "Falure Date",
        average_impact: "Impact"
    };

    return (

        <TableComposable aria-label="Error messages" variant="compact">
            <Thead>
                <Tr>
                    <Th width={80} rowSpan={2} />
                    <Th>Test Name</Th>
                    <Th width={10} sort={getSimpleSortParams(1)}>
                        Impact
                    </Th>
                    <Th width={10} sort={getSimpleSortParams(0)}>
                        Count
                    </Th>
                </Tr>
            </Thead>
            {test_cases && test_cases.sort(onSimpleSort).map((test_case, tc_idx) => (
                <Tbody key={test_case.name + tc_idx}>
                    <Tr>
                        <Td expand={{ rowIndex, isExpanded: isTestCaseExpanded(test_case.name), onToggle: () => setTestCaseExpanded(test_case.name, !isTestCaseExpanded(test_case.name)) }} />
                        <Td>
                            {test_case.name}
                        </Td>
                        <Td>
                            {test_case.test_case_impact}
                        </Td>
                        <Td>
                            {test_case.count}
                        </Td>
                    </Tr>
                    <Tr isExpanded={isTestCaseExpanded(test_case.name)}>
                        <Td></Td>
                        <Td colSpan={3}>
                            <ExpandableRowContent>
                                <TableComposable aria-label="Error messages" variant="compact">
                                    <Thead>
                                        <Tr>
                                            <Th />
                                            <Th width={10}>Job URL</Th>
                                            <Th width={70}>Error</Th>
                                            <Th width={20}>Failure Dates</Th>
                                        </Tr>
                                    </Thead>
                                    {test_case.messages && test_case.messages.map((message, m_idx) => (
                                        <Tbody key={message.job_id + m_idx}>
                                            <Tr>
                                                <Td></Td>
                                                <Td dataLabel={columnNames.job_id}>
                                                    <a href={message.job_url} rel="noreferrer noopener" target='_blank'>{message.job_id}</a>
                                                </Td>
                                                <Td dataLabel="Error messages" onClick={expandPre} className='collapsedPre'>
                                                    <p style={{ textAlign: 'center', color: 'var(--pf-global--link--Color)', cursor: "pointer" }}><u>Show error</u></p>
                                                    <pre>
                                                        {message.error_message}
                                                    </pre>
                                                </Td>
                                                <Td dataLabel={columnNames.failure_date}>
                                                    {message.failure_date}
                                                </Td>
                                            </Tr>
                                        </Tbody>
                                    ))
                                    }
                                </TableComposable>
                            </ExpandableRowContent>
                        </Td>
                    </Tr>
                </Tbody>
            ))}
        </TableComposable>
    )
}

export const ComposableTableNestedExpandable: React.FunctionComponent<{ teams: Flakey[] }> = ({ teams }) => {

    const columnNames = {
        name: 'Test Case',
        status: 'Status',
        error_message: 'Failing Jobs',
        count: 'Count',
        suite_name: 'Suite Name',
        job_id: "Job ID",
        job_url: "Job Url",
        failure_date: "Falure Date",
        average_impact: "Impact"
    };

    // Exapndable suites
    const [expandedSuitesNames, setExpandedSuitesNames] = React.useState<string[]>([]);

    const setSuiteExpanded = (suite_name: string, isExpanding = true) => {
        setExpandedSuitesNames(prevExpanded => {
            const otherExpandedSuiteNames = prevExpanded.filter(t => t !== suite_name);
            return isExpanding ? [...otherExpandedSuiteNames, suite_name] : otherExpandedSuiteNames;
        });
    }
    const isSuiteExpanded = (suite_name: string) => expandedSuitesNames.includes(suite_name);


    const [activeSortIndex, setActiveSortIndex] = React.useState<number>(0);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc'>('desc');

    const getSortParams = (columnIndex: number): ThProps['sort'] => ({
        sortBy: {
            index: activeSortIndex,
            direction: activeSortDirection,
            defaultDirection: 'desc'
        },
        onSort: (_event, index, direction) => {
            setActiveSortIndex(index);
            setActiveSortDirection(direction);
        },
        columnIndex
    });

    const getSortableRowValues = (flakey: Flakey): (string | number)[] => {
        const { average_impact, suite_name } = flakey;
        return [average_impact, suite_name];
    };

    const onSortFn = (a: Flakey, b: Flakey): number => {

        const aValue = getSortableRowValues(a)[activeSortIndex];
        const bValue = getSortableRowValues(b)[activeSortIndex];
        if (typeof aValue === 'number') {
            // Numeric sort
            if (activeSortDirection === 'asc') {
                return (aValue as number) - (bValue as number);
            }
            return (bValue as number) - (aValue as number);
        } else {
            // String sort
            if (activeSortDirection === 'asc') {
                return (aValue as string).localeCompare(bValue as string);
            }
            return (bValue as string).localeCompare(aValue as string);
        }
    }

    const countFailingSuites = (suite: Flakey) => {
        const sum = suite.test_cases.reduce((accumulator, object) => {
            return accumulator + object.count;
        }, 0);
        return sum
    }

    return (
        <InnerScrollContainer>
            <TableComposable aria-label="Nested column headers with expandable rows table" gridBreakPoint="">
                <Thead hasNestedHeader>
                    <Tr>
                        <Th rowSpan={2} />
                        <Th width={50} >
                            {columnNames.suite_name}
                        </Th>
                        <Th width={10} sort={getSortParams(0)}>
                            {columnNames.average_impact}
                        </Th>
                    </Tr>
                </Thead>
                {teams.sort(onSortFn).map((suite, rowIndex) => (
                    <Tbody key={suite.suite_name + rowIndex} isExpanded={isSuiteExpanded(suite.suite_name)}>
                        <Tr>
                            <Td expand={{ rowIndex, isExpanded: isSuiteExpanded(suite.suite_name), onToggle: () => setSuiteExpanded(suite.suite_name, !isSuiteExpanded(suite.suite_name)) }} />
                            <Td dataLabel={columnNames.name}>{suite.suite_name}</Td>
                            <Td dataLabel={columnNames.count}>{suite.average_impact.toFixed(2)}%</Td>
                        </Tr>
                        <Tr isExpanded={isSuiteExpanded(suite.suite_name)} className='pf-px-xl'>
                            <Td colSpan={3}>
                                <ExpandableRowContent>
                                    <Grid hasGutter>
                                        <GridItem span={3}>
                                            <Card className='card-no-border'>
                                                <CardTitle component="h3" style={{ color: "red" }}>Overall Impact</CardTitle>
                                                <CardBody>{suite.average_impact.toFixed(2)}%</CardBody>
                                            </Card>
                                            <Card className='card-no-border'>
                                                <CardTitle component="h3">Total count of failed test cases</CardTitle>
                                                <CardBody>{isSuiteExpanded(suite.suite_name) ? countFailingSuites(suite) : ""}</CardBody>
                                            </Card>
                                        </GridItem>
                                        <GridItem span={9}>
                                            <Card className='card-no-border'>
                                                <CardTitle>Failing test cases</CardTitle>
                                                <InnerNestedComposableTableNestedExpandable test_cases={suite.test_cases} rowIndex={rowIndex}></InnerNestedComposableTableNestedExpandable>
                                            </Card>
                                        </GridItem>
                                    </Grid>
                                </ExpandableRowContent>
                            </Td>
                        </Tr>
                    </Tbody>
                ))}
            </TableComposable>
        </InnerScrollContainer>
    );
};