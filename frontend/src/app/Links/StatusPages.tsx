/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import React from 'react';
import {
    Card,
    CardBody,
    CardHeader,
    CardTitle,
    Grid,
    GridItem,
    Title,
} from '@patternfly/react-core';
export const StatusPages = () => {
    return (
        <Grid hasGutter>
            <GridItem span={8}>
                <Card isLarge>
                    <CardTitle>
                        <Title headingLevel="h1" size="xl">
                            Status pages links
                        </Title>
                    </CardTitle>
                    <CardBody style={{ paddingLeft: '16px' }}>
                        <Card isPlain isCompact >
                            <CardTitle>Red Hat Status Page</CardTitle>
                            <CardBody>
                                <a href="https://status.redhat.com/" target="blank" rel="noreferrer"><b>https://status.redhat.com/</b></a>
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact >
                            <CardTitle>Quay.io Status Page</CardTitle>
                            <CardBody>
                                <a href="https://status.quay.io/" target="blank" rel="noreferrer"><b>https://status.quay.io/</b></a>
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact>
                            <CardHeader><b>GitHub Status Page</b></CardHeader>
                            <CardBody>
                                <a href="https://www.githubstatus.com/" target="blank" rel="noreferrer"><b>https://www.githubstatus.com/</b></a>
                            </CardBody>
                        </Card>
                        <Card isPlain isCompact>
                            <CardHeader><b>AWS Status Page</b></CardHeader>
                            <CardBody>
                                <a href="https://health.aws.amazon.com/health/status" target="blank" rel="noreferrer"><b>https://health.aws.amazon.com/health/status</b></a>
                            </CardBody>
                        </Card>
                    </CardBody>
                </Card>
            </GridItem>
        </Grid>
    )
};
