/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import React, { useContext, useEffect, useState } from 'react';
import {
    Card,
    CardBody,
    CardHeader,
    CardTitle,
    DescriptionList,
    DescriptionListDescription,
    DescriptionListGroup,
    DescriptionListTerm,
    Grid,
    GridItem,
    Stack,
    StackItem,
    Title,
} from '@patternfly/react-core';
import { getVersion } from '@app/utils/APIService';
import { ReactReduxContext } from 'react-redux';
import { CheckCircleIcon, ExclamationCircleIcon } from '@patternfly/react-icons';

export const About = () => {
    const { store } = useContext(ReactReduxContext);
    const dispatch = store.dispatch;

    /*
      SERVER INFO AND REPOSITORY TABLE METHODS
    */
    const [dashboardVersion, setVersion] = useState('unknown')
    const [serverAvailable, setServerAvailable] = useState<boolean>(false)

    useEffect(() => {
        getVersion().then((res) => { // making the api call here
            if (res.code === 200) {
                const result = res.data;
                dispatch({ type: "SET_Version", data: result['serverAPIVersion'] });
                // not really required to store it in the global state , just added it to make it better understandable
                setVersion(result['serverAPIVersion'])
                setServerAvailable(true)
            } else {
                setServerAvailable(false)
                dispatch({ type: "SET_ERROR", data: res });
            }
        });
    }, [dashboardVersion, setVersion, dispatch])

    return (
        <Grid hasGutter>
            <GridItem span={8}>
                <Card isLarge>
                    <CardTitle>
                        <Title headingLevel="h1" size="xl">
                            About
                        </Title>
                    </CardTitle>
                    <CardBody style={{ paddingLeft: '16px' }}>
                        <Card isPlain isCompact >
                            <CardTitle>What Is Quality Studio?</CardTitle>
                            <CardBody>
                                Quality Studio tries to unify in a single Dashboard all quality information about your product; from your commits to the costumers feedback.
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact >
                            <CardTitle>Connectors for everything</CardTitle>
                            <CardBody>Onboard team for your product and start to explore Quality Studio <b>Connectors</b>.</CardBody>
                        </Card>
                        <Card isPlain isCompact>
                            <CardHeader><b>Jira Connector</b></CardHeader>
                            <CardBody>
                                Measure your bugs resolution times or retrieve information about what bugs are blocking your tests.
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact>
                            <CardHeader><b>Openshift CI Connector</b></CardHeader>
                            <CardBody>
                                Connect with openshift CI and track information about presubmit, periodic or postsubmit jobs. Quality Studio offer a complete set of metrics about the impact of the tests execution
                                in your repositories. Note that only the repositories with OpenShift CI prow jobs are present, as well as the job types available for each repository.
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact>
                            <CardHeader><b>GitHub Connector</b></CardHeader>
                            <CardBody>
                                Observe detailed information from GitHub like GitHub Workflows, Pull Requests, Contributors.
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact>
                            <CardHeader><b>Coverage Connector</b></CardHeader>
                            <CardBody>
                                Connect to <a href="https://about.codecov.io/" target="blank" rel="noopener noreferrer"><b>CodeCov</b></a> and get the total coverage of your repository.
                            </CardBody>
                        </Card>
                    </CardBody>
                </Card>
            </GridItem>
            <GridItem span={4}>
                <Stack hasGutter>
                    <StackItem>
                        <Card isLarge>
                            <CardTitle>
                                <Title headingLevel="h1" size="xl">
                                    Red Hat Quality Studio Details
                                </Title>
                            </CardTitle>
                            <CardBody>
                                <DescriptionList>
                                    <DescriptionListGroup>
                                        <DescriptionListTerm>Quality Studio version</DescriptionListTerm>
                                        <DescriptionListDescription>
                                            <span>{dashboardVersion}</span>
                                        </DescriptionListDescription>
                                    </DescriptionListGroup>
                                    <DescriptionListGroup>
                                        <DescriptionListTerm>Server API Status</DescriptionListTerm>
                                        <DescriptionListDescription>
                                            {serverAvailable && <span style={{ color: "darkgreen", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500 }}> <CheckCircleIcon size={'sm'} ></CheckCircleIcon> OK </span>}
                                            {!serverAvailable && <span style={{ color: "darkred", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500 }}> <ExclamationCircleIcon size={'sm'} ></ExclamationCircleIcon> DOWN </span>}
                                        </DescriptionListDescription>
                                    </DescriptionListGroup>
                                    <DescriptionListGroup>
                                        <DescriptionListTerm>Database Status</DescriptionListTerm>
                                        <DescriptionListDescription>
                                            {serverAvailable && <span style={{ color: "darkgreen", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500 }}> <CheckCircleIcon size={'sm'} ></CheckCircleIcon> OK </span>}
                                            {!serverAvailable && <span style={{ color: "darkred", verticalAlign: "middle", lineHeight: "2em", fontWeight: 500 }}> <ExclamationCircleIcon size={'sm'} ></ExclamationCircleIcon> DOWN </span>}
                                        </DescriptionListDescription>
                                    </DescriptionListGroup>
                                </DescriptionList>
                            </CardBody>
                        </Card>
                    </StackItem>
                    <StackItem>
                        <Card isLarge>
                            <CardTitle>
                                <Title headingLevel="h1" size="xl">
                                    Documentation
                                </Title>
                            </CardTitle>
                            <CardBody>
                                Find our <a href="https://github.com/redhat-appstudio/quality-dashboard" target="_blank" rel="noreferrer"><b>code at GitHub</b></a>
                            </CardBody>
                        </Card>
                    </StackItem>
                </Stack>
            </GridItem>
        </Grid>
    )
};
