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
                                <BugIcon />
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
                                <MonitoringIcon />
                            </Icon>
                        </CardTitle>
                        <CardBody>
                            <Bullseye>
                                Openshift CI metrics
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
                                <GithubIcon />
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