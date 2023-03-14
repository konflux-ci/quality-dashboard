import React, { useEffect } from "react";
import { TableComposable, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { Card, CardBody, CardTitle, Pagination } from "@patternfly/react-core";
import { formatDate } from "./utils";

const columnsTests = {
    date: 'Date',
    job_type: 'Job Type',
    job_name: 'Job Name',
    job_id: 'Job ID',
    e2e_failed_messages: 'Failed E2E Tests',
}

export interface Job {
    job_id: string;
    e2e_failed_test_messages: string;
    suites_xml_url: string;
    created_at: string;
    job_type: string;
    job_name: string;
}

export const getFailedProwJobsInE2ETests = (prowJobs: Job[]) => {
    const failedProwJobs = new Array<Job>

    prowJobs?.map(job => {
        if (job.e2e_failed_test_messages != undefined) {
            failedProwJobs.push(job)
        }
    })

    failedProwJobs.sort((a, b) => (a.created_at < b.created_at ? -1 : 1));

    return failedProwJobs
}


export const FailedE2ETests: React.FC<{ failedProwJobs: any }> = ({ failedProwJobs }) => {
    const [failedProwJobsPage, setFailedProwJobsPage] = React.useState<Array<Job>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(4);

    useEffect(() => {
        if (failedProwJobs.length == 0) {
            setPage(1)
            setFailedProwJobsPage([])
        }
        if (failedProwJobs.length > 0) {
            setFailedProwJobsPage(failedProwJobs.slice(0, perPage))
            setPage(1)
        }
    }, [failedProwJobs]);


    useEffect(() => {
        if (failedProwJobs.length > 0) {
            let from = (page - 1) * perPage
            let to = (page - 1) * perPage + perPage >= failedProwJobs.length ? failedProwJobs.length : (page - 1) * perPage + perPage;
            setFailedProwJobsPage(failedProwJobs.slice(from, to))
        }
    }, [page, perPage]);



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


    return (
        <Card style={{ width: "100%", height: "100%", fontSize: "1rem" }}>
            <CardTitle>Total Failed Jobs In E2E Tests ({failedProwJobs.length})</CardTitle>
            <CardBody>
                <Pagination
                    perPageComponent="button"
                    itemCount={failedProwJobs.length}
                    perPage={perPage}
                    page={page}
                    onSetPage={onSetPage}
                    widgetId="top-example"
                    onPerPageSelect={onPerPageSelect}
                />
                <TableComposable aria-label="Actions table">
                    <Thead>
                        <Tr>
                            <Th>{columnsTests.date}</Th>
                            <Th>{columnsTests.job_type}</Th>
                            <Th>{columnsTests.job_name}</Th>
                            <Th>{columnsTests.job_id}</Th>
                            <Th>{columnsTests.e2e_failed_messages}</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {failedProwJobsPage?.map(job => {
                            if (job.e2e_failed_test_messages != undefined) {
                                return (
                                    <Tr key={job.job_id}>
                                        <Td>{formatDate(new Date(job.created_at))}</Td>
                                        <Td>{job.job_type}</Td>
                                        <Td>{job.job_name}</Td>
                                        <Td>
                                            <a href={job.suites_xml_url} target="blank" rel="noopener noreferrer">
                                                {job.job_id}
                                            </a>
                                        </Td>
                                        <Td style={{ whiteSpace: "pre-line" }}>{job.e2e_failed_test_messages}</Td>
                                    </Tr>
                                );
                            }
                        })}
                    </Tbody>
                </TableComposable>
                <Pagination
                    perPageComponent="button"
                    itemCount={failedProwJobs.length}
                    perPage={perPage}
                    page={page}
                    onSetPage={onSetPage}
                    widgetId="top-example"
                    onPerPageSelect={onPerPageSelect}
                />
            </CardBody>
        </Card>
    );
}