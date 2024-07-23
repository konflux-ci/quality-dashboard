/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import * as React from 'react';
import {
    Card,
    CardBody,
    Flex,
    FlexItem,
    CardTitle,
    Bullseye,
    TextVariants,
    Divider,
    Text,
    Icon,
} from '@patternfly/react-core';
import { BugIcon, GithubIcon, MonitoringIcon } from '@patternfly/react-icons';

export const InfoBanner = () => (
    <Card isLarge>
        <CardBody style={{ paddingLeft: '16px' }}>
            <Flex
                justifyContent={{ default: 'justifyContentSpaceEvenly' }}
                flexWrap={{ default: 'nowrap' }}
                direction={{ default: 'column', sm: 'row' }}
            >
                <FlexItem flex={{ default: 'flex_1' }}>
                    <Card isPlain isCompact>
                        <CardTitle style={{ textAlign: "center" }}>
                            <Icon size="xl" iconSize="lg">
                                <BugIcon color='#bd2c00'/>
                            </Icon>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>
                                <Text component={TextVariants.p}>
                                    Jira Bugs Observability
                                </Text>
                            </Bullseye>
                        </CardBody>
                    </Card>
                </FlexItem>
                <Divider
                    orientation={{
                        default: 'vertical',
                    }}
                />
                <FlexItem flex={{ default: 'flex_1' }}>
                    <Card isPlain isCompact>
                        <CardTitle style={{ textAlign: "center" }}>
                            <Icon size="xl" iconSize="lg">
                                <MonitoringIcon color='green'/>
                            </Icon>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>
                                CI metrics
                            </Bullseye>
                        </CardBody>
                    </Card>
                </FlexItem>
                <Divider
                    orientation={{
                        default: 'vertical',
                    }}
                />
                <FlexItem flex={{ default: 'flex_1' }}>
                    <Card isPlain isCompact>
                        <CardTitle style={{ textAlign: "center" }}>
                            <Icon size="xl" iconSize="lg">
                                <GithubIcon color='#4078c0'/>
                            </Icon>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>GitHub metrics</Bullseye>
                        </CardBody>
                    </Card>
                </FlexItem>
            </Flex>
        </CardBody>
    </Card>
);
