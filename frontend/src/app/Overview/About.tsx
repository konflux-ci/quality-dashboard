import React, { useContext, useEffect, useState } from 'react';
import {
    Card,
    CardBody,
    CardTitle,
    DescriptionList,
    DescriptionListDescription,
    DescriptionListGroup,
    DescriptionListTerm,
    Gallery,
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
        <Gallery hasGutter style={{ display: "flex" }}>
            <Card isRounded isCompact style={{ width: "65%", padding: "5px" }}>
                <CardTitle>
                    <Title headingLevel="h1" size="xl">
                        About
                    </Title>
                </CardTitle>
                <CardBody style={{ paddingLeft: '16px' }}>
                    <Card isPlain isCompact >
                        <CardTitle>Jira Issues</CardTitle>
                        <CardBody>
                            Check the Jira issues metrics based on projects STONE, DEVHAS, SRVKP, GITOPSRVCE, and HACBS with label ci-fail.
                        </CardBody>
                    </Card>
                    <Card isPlain isCompact>
                        <CardTitle>OpenShift CI prow jobs</CardTitle>
                        <CardBody>
                            Track OpenShift CI by job type (presubmit, periodic, or postsubmit).
                            Note that are only present the repositories with OpenShift CI prow jobs, as well as the job types available for each repository.
                        </CardBody>
                    </Card>
                    <Card isPlain isCompact>
                        <CardTitle>GitHub Actions</CardTitle>
                        <CardBody>
                            Observe the last execution of the each GitHub Action.
                        </CardBody>
                    </Card>
                </CardBody>
            </Card>
            <Card isRounded style={{ width: "35%" }}>
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
        </Gallery>
    )
};